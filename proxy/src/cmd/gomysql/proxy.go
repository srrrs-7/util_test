package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/go-mysql-org/go-mysql/client"
	"github.com/go-mysql-org/go-mysql/server"
)

// Default timeout for graceful shutdown
const (
	MinConnections         = 16
	MaxConnections         = 64
	MaxIdleConnections     = 32
	DefaultShutdownTimeout = 30 * time.Second
)

var errorConnectionClosed = errors.New("connection closed")

var atomPool atomic.Pointer[client.Pool]

// DBConfig is a struct that holds database connection configuration information
type DBConfig struct {
	DBName string
	Addr   string
	User   string
	Pass   string
}

// ProxyServer represents the MySQL proxy server
type ProxyServer struct {
	listener     net.Listener
	listenDbConf DBConfig
	targetDbConf DBConfig
	ctx          context.Context
	cancel       context.CancelFunc
	wg           sync.WaitGroup
	shutdownCh   chan struct{}
}

// NewProxyServer creates a new ProxyServer instance
func NewProxyServer(pAddr, pUser, pPass, tAddr, tUser, tPass, tDbName string) *ProxyServer {
	ctx, cancel := context.WithCancel(context.Background())
	return &ProxyServer{
		listenDbConf: DBConfig{Addr: pAddr, User: pUser, Pass: pPass},
		targetDbConf: DBConfig{Addr: tAddr, User: tUser, Pass: tPass, DBName: tDbName},
		wg:           sync.WaitGroup{},
		ctx:          ctx,
		cancel:       cancel,
		shutdownCh:   make(chan struct{}),
	}
}

// Start starts the MySQL proxy server
func (s *ProxyServer) Start() error {
	proxyAddr := s.listenDbConf.Addr
	l, err := net.Listen("tcp", proxyAddr)
	if err != nil {
		return fmt.Errorf("failed to listen on address: %w", err)
	}
	s.listener = l
	log.Printf("MySQL proxy server successfully listening on %s", proxyAddr)

	go s.handleSignals()

	pool, err := client.NewPoolWithOptions(
		s.targetDbConf.Addr,
		s.targetDbConf.User,
		s.targetDbConf.Pass,
		s.targetDbConf.DBName,
		client.WithPoolLimits(MinConnections, MaxConnections, MaxIdleConnections),
		client.WithLogger(slog.Default()),
		client.WithNewPoolPingTimeout(5*time.Second),
	)
	if err != nil {
		return fmt.Errorf("failed to create connection pool: %w", err)
	}
	atomPool.Store(pool)

	// Loop to accept connections
	for {
		conn, err := l.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				log.Println("Listener closed, exiting accept loop")
				break
			}
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		s.wg.Add(1)
		go s.handleConnection(conn)
	}

	return nil
}

// handleSignals handles OS signals for graceful shutdown
func (s *ProxyServer) handleSignals() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	log.Println("Received shutdown signal, shutting down...")
	s.Shutdown()
}

// handleConnection processes a new client connection
func (s *ProxyServer) handleConnection(conn net.Conn) {
	defer s.wg.Done()

	go func() {
		<-s.ctx.Done()
		log.Printf("Shutting down connection from %s", conn.RemoteAddr().String())
		conn.SetDeadline(time.Now())
	}()

	client, err := atomPool.Load().GetConn(context.Background())
	if err != nil {
		log.Printf("Failed to get connection from pool: %v", err)
		conn.Close()
		return
	}
	defer atomPool.Load().PutConn(client)

	// Create MySQL server connection
	serverConn, err := server.NewDefaultServer().NewConn(
		conn,
		s.listenDbConf.User,
		s.listenDbConf.Pass,
		NewQueryHandler(client),
	)
	if err != nil {
		log.Printf("Failed to create server connection: %v", err)
		return
	}

	// Process commands until client disconnects
	for {
		if err := serverConn.HandleCommand(); err != nil {
			if errors.Is(err, errorConnectionClosed) {
				log.Printf("Client disconnected: %v", err)
			} else {
				log.Printf("Error handling command: %v", err)
			}
			break
		}
	}
}

// Shutdown gracefully shuts down the proxy server
func (s *ProxyServer) Shutdown() {
	s.cancel()
	s.listener.Close()
	close(s.shutdownCh)

	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Println("Server shutdown complete")
	case <-time.After(DefaultShutdownTimeout):
		log.Println("Server shutdown timed out")
	}
}
