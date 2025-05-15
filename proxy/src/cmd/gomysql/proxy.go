package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-mysql-org/go-mysql/client"
	"github.com/go-mysql-org/go-mysql/server"
)

// Connection pool configuration constants
const (
	MinConnections     = 16
	MaxConnections     = 64
	MaxIdleConnections = 32
	PoolPingTimeout    = 5 * time.Second
)

// Timeout constants
const (
	DefaultShutdownTimeout = 30 * time.Second
)

// Common error definitions
var (
	ErrConnectionClosed = errors.New("connection closed")
)

// DBConfig holds database connection configuration
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
	pool, err := s.createConnectionPool()
	if err != nil {
		return err
	}

	// Store the connection pool in the atomic pointer
	s.atomPool.Store(pool)

	// Loop to accept connections
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				slog.Info("Listener closed, exiting accept loop")
				break
			}
			slog.Error("Failed to accept connection", "error", err)
			continue
		}

		s.wg.Add(1)
		go s.handleConnection(ctx, conn)
	}

	return nil
}

// createConnectionPool initializes the connection pool to the target database
func (s *ProxyServer) createConnectionPool() (*client.Pool, error) {
	pool, err := client.NewPoolWithOptions(
		s.targetDbConf.Addr,
		s.targetDbConf.User,
		s.targetDbConf.Pass,
		s.targetDbConf.DBName,
		client.WithPoolLimits(MinConnections, MaxConnections, MaxIdleConnections),
		client.WithLogger(slog.Default()),
		client.WithNewPoolPingTimeout(PoolPingTimeout),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}
	return pool, nil
}

// handleConnection processes a new client connection
func (s *ProxyServer) handleConnection(ctx context.Context, conn net.Conn) {
	defer func() {
		s.wg.Done()
		conn.Close()
	}()

	dbConn, err := s.getDBConnection(ctx)
	if err != nil {
		slog.Error("Failed to get connection from pool", "error", err)
		return
	}
	defer s.atomPool.Load().PutConn(dbConn)

	serverConn, err := s.createServerConnection(conn, dbConn)
	if err != nil {
		slog.Error("Failed to create server connection", "error", err)
		return
	}

	s.processCommands(serverConn)
}

// getDBConnection retrieves a connection from the pool
func (s *ProxyServer) getDBConnection(ctx context.Context) (*client.Conn, error) {
	pool := s.atomPool.Load()
	if pool == nil {
		return nil, errors.New("connection pool not initialized")
	}
	return pool.GetConn(ctx)
}

// createServerConnection sets up the MySQL server connection
func (s *ProxyServer) createServerConnection(conn net.Conn, dbConn *client.Conn) (*server.Conn, error) {
	return server.NewDefaultServer().NewConn(
		conn,
		s.listenDbConf.User,
		s.listenDbConf.Pass,
		NewQueryHandler(dbConn),
	)
}

// processCommands handles MySQL commands until the client disconnects
func (s *ProxyServer) processCommands(serverConn *server.Conn) {
	for {
		if err := serverConn.HandleCommand(); err != nil {
			if err.Error() != ErrConnectionClosed.Error() {
				slog.Error("Error handling command", "error", err)
			}
			break
		}
	}
}

// Shutdown gracefully shuts down the proxy server
func (s *ProxyServer) Shutdown() {
	if err := s.listener.Close(); err != nil {
		slog.Error("Failed to close listener", "error", err)
	}

	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	if pool := s.atomPool.Load(); pool != nil {
		pool.Close()
	}

	select {
	case <-done:
		slog.Info("Server shutdown complete")
	case <-time.After(DefaultShutdownTimeout):
		slog.Info("Server shutdown timed out")
	}
}
