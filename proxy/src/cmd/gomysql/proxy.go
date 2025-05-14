package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net"
	"sync"
	"sync/atomic"
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
	wg           *sync.WaitGroup
	atomPool     atomic.Pointer[client.Pool]
}

// NewProxyServer creates a new ProxyServer instance
func NewProxyServer(pConf, tConf DBConfig, l net.Listener, wg *sync.WaitGroup) *ProxyServer {
	return &ProxyServer{
		listener:     l,
		listenDbConf: pConf,
		targetDbConf: tConf,
		wg:           wg,
		atomPool:     atomic.Pointer[client.Pool]{},
	}
}

// Start starts the MySQL proxy server
func (s *ProxyServer) Start(ctx context.Context) error {
	// Create a connection pool to the target database
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
	// Store the connection pool in the atomic pointer
	s.atomPool.Store(pool)

	// Loop to accept connections
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				log.Println("Listener closed, exiting accept loop")
				break
			}
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		s.wg.Add(1)
		go s.handleConnection(ctx, conn)
	}

	return nil
}

// handleConnection processes a new client connection
func (s *ProxyServer) handleConnection(ctx context.Context, conn net.Conn) {
	defer func() {
		s.wg.Done()
		conn.Close()
	}()

	client, err := s.atomPool.Load().GetConn(ctx)
	if err != nil {
		log.Printf("Failed to get connection from pool: %v", err)
		conn.Close()
		return
	}
	defer s.atomPool.Load().PutConn(client)

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
			if err.Error() != errorConnectionClosed.Error() {
				log.Printf("Error handling command: %v", err)
			}

			break
		}
	}
}

// Shutdown gracefully shuts down the proxy server
func (s *ProxyServer) Shutdown(ctx context.Context) {
	select {
	case <-ctx.Done():
		log.Println("Shutdown context done")
		if err := s.listener.Close(); err != nil {
			log.Printf("Failed to close listener: %v", err)
		}
	}

	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	s.atomPool.Load().Close()

	select {
	case <-done:
		log.Println("Server shutdown complete")
	case <-time.After(DefaultShutdownTimeout):
		log.Println("Server shutdown timed out")
	}
}
