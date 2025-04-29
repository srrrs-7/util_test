package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

// Proxy configuration
type ProxyConfig struct {
	ListenAddr string // Address the proxy listens on (e.g., localhost:3307)
	TargetAddr string // Address of the actual MySQL server (e.g., localhost:3306)
	LogTraffic bool   // Whether to log traffic
}

// MySQL Proxy structure
type MySQLProxy struct {
	config ProxyConfig
	logger *log.Logger
}

func main() {
	// Proxy configuration
	config := ProxyConfig{
		ListenAddr: "0.0.0.0:8080", // Specify the port for the proxy to listen on
		TargetAddr: "mysql:3306",   // Address of the actual MySQL server
		LogTraffic: true,           // Log traffic
	}

	// Create and start the proxy
	proxy := NewMySQLProxy(config)
	if err := proxy.Start(); err != nil {
		log.Fatalf("Failed to start proxy: %v", err)
	}
}

// Create a new MySQL proxy instance
func NewMySQLProxy(config ProxyConfig) *MySQLProxy {
	return &MySQLProxy{
		config: config,
		logger: log.New(os.Stdout, "[MySQL Proxy] ", log.LstdFlags),
	}
}

// Start the proxy
func (p *MySQLProxy) Start() error {
	// Create a TCP listener
	listener, err := net.Listen("tcp", p.config.ListenAddr)
	if err != nil {
		return fmt.Errorf("failed to create listener: %v", err)
	}
	defer listener.Close()

	p.logger.Printf("MySQL Proxy started on %s (target: %s)",
		p.config.ListenAddr, p.config.TargetAddr)

	// Wait for and handle connections
	for {
		clientConn, err := listener.Accept()
		if err != nil {
			p.logger.Printf("failed to accept connection: %v", err)
			continue
		}

		p.logger.Printf("new connection: %s", clientConn.RemoteAddr())

		// Handle each connection in a separate goroutine
		go p.handleConnection(clientConn)
	}
}

// Handle client connections
func (p *MySQLProxy) handleConnection(clientConn net.Conn) {
	defer clientConn.Close()

	// Connect to the actual MySQL server
	serverConn, err := net.Dial("tcp", p.config.TargetAddr)
	if err != nil {
		p.logger.Printf("failed to connect to target MySQL server: %v", err)
		return
	}
	defer serverConn.Close()

	// Handle bidirectional data transfer in separate goroutines
	// Client -> Server
	go func() {
		if _, err := p.transfer(clientConn, serverConn, "C->S"); err != nil && err != io.EOF {
			p.logger.Printf("client->server transfer error: %v", err)
		}
		// Close server connection if client disconnects
		serverConn.Close()
	}()

	// Server -> Client
	if _, err := p.transfer(serverConn, clientConn, "S->C"); err != nil && err != io.EOF {
		p.logger.Printf("server->client transfer error: %v", err)
	}

	p.logger.Printf("connection closed: %s", clientConn.RemoteAddr())
}

// Data transfer with optional logging
func (p *MySQLProxy) transfer(src, dst net.Conn, direction string) (int64, error) {
	if p.config.LogTraffic {
		// Transfer while logging
		buffer := make([]byte, 4096)
		var total int64

		for {
			// Set a timeout
			src.SetReadDeadline(time.Now().Add(5 * time.Minute))

			// Read data
			n, err := src.Read(buffer)
			if err != nil {
				return total, err
			}
			total += int64(n)

			// Log the transfer (actual MySQL protocol parsing is not implemented)
			p.logger.Printf("%s: %d bytes transferred", direction, n)

			// Set a timeout for writing
			dst.SetWriteDeadline(time.Now().Add(30 * time.Second))

			// Write data
			_, err = dst.Write(buffer[:n])
			if err != nil {
				return total, err
			}
		}
	} else {
		// Simple transfer using io.Copy
		return io.Copy(dst, src)
	}
}
