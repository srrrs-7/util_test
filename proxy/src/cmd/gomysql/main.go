package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Constants for environment variable names
const (
	EnvProxyAddr = "PROXY_ADDR"
	EnvProxyUser = "PROXY_USER"
	EnvProxyPass = "PROXY_PASS"

	EnvTestDBAddr = "TEST_DB_ADDR"
	EnvTestDBUser = "TEST_DB_USER"
	EnvTestDBPass = "TEST_DB_PASS"
	EnvTestDBName = "TEST_DB_NAME"

	// Timeout for graceful shutdown
	ShutdownTimeout = 30 * time.Second
)

// Config is a struct that holds database connection information loaded from environment variables
type Config struct {
	proxy  DBConfig
	testDB DBConfig
}

// loadConfig loads settings from environment variables and returns a Config struct
func loadConfig() (*Config, error) {
	// Load proxy config
	proxyAddr := os.Getenv(EnvProxyAddr)
	proxyUser := os.Getenv(EnvProxyUser)
	proxyPass := os.Getenv(EnvProxyPass)

	// Load test DB config
	testDBAddr := os.Getenv(EnvTestDBAddr)
	testDBUser := os.Getenv(EnvTestDBUser)
	testDBPass := os.Getenv(EnvTestDBPass)
	testDBName := os.Getenv(EnvTestDBName)

	// Validate required fields
	if proxyAddr == "" {
		return nil, fmt.Errorf("missing required environment variable: %s", EnvProxyAddr)
	}
	if proxyUser == "" {
		return nil, fmt.Errorf("missing required environment variable: %s", EnvProxyUser)
	}
	if proxyPass == "" {
		return nil, fmt.Errorf("missing required environment variable: %s", EnvProxyPass)
	}
	if testDBAddr == "" {
		return nil, fmt.Errorf("missing required environment variable: %s", EnvTestDBAddr)
	}
	if testDBUser == "" {
		return nil, fmt.Errorf("missing required environment variable: %s", EnvTestDBUser)
	}
	if testDBPass == "" {
		return nil, fmt.Errorf("missing required environment variable: %s", EnvTestDBPass)
	}
	if testDBName == "" {
		return nil, fmt.Errorf("missing required environment variable: %s", EnvTestDBName)
	}

	return &Config{
		proxy: DBConfig{
			Addr: proxyAddr,
			User: proxyUser,
			Pass: proxyPass,
		},
		testDB: DBConfig{
			Addr:   testDBAddr,
			User:   testDBUser,
			Pass:   testDBPass,
			DBName: testDBName,
		},
	}, nil
}

// setupNetListener creates a TCP listener for the proxy server
func setupNetListener(addr string) (net.Listener, error) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on proxy address: %w", err)
	}
	return l, nil
}

// setupSignalHandler creates a context that is canceled when the process receives a termination signal
func setupSignalHandler(cancel context.CancelFunc, server *ProxyServer) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		slog.Info("Received shutdown signal", "signal", sig)
		cancel()
		server.Shutdown()
	}()
}

// main is the entry point of the application
func main() {
	// Configure structured logging with JSON format
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	// Load configuration from environment variables
	config, err := loadConfig()
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		panic(fmt.Sprintf("Failed to load configuration: %v", err))
	}

	// Set up network listener
	listener, err := setupNetListener(config.proxy.Addr)
	if err != nil {
		slog.Error("Failed to set up network listener", "error", err)
		panic(fmt.Sprintf("Failed to set up network listener: %v", err))
	}
	defer listener.Close()

	slog.Info("MySQL proxy server successfully listening", "addr", config.proxy.Addr)

	// Create context with cancellation for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create and start the proxy server
	wg := &sync.WaitGroup{}
	proxyServer := NewProxyServer(
		config.proxy,
		config.testDB,
		listener,
		wg,
	)

	// Set up signal handler for graceful shutdown
	setupSignalHandler(cancel, proxyServer)

	// Start the proxy server
	if err := proxyServer.Start(ctx); err != nil {
		slog.Error("Failed to start proxy server", "error", err)
		panic(fmt.Sprintf("Failed to start proxy server: %v", err))
	}

	// Wait for context cancellation (triggered by signal handler)
	<-ctx.Done()

	// Wait for a grace period for server to shut down
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer shutdownCancel()

	select {
	case <-shutdownCtx.Done():
		slog.Info("Shutdown completed")
	}
}
