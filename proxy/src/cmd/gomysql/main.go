package main

import (
	"log"
	"net"

	"github.com/go-mysql-org/go-mysql/server"
)

const (
	LISTEN_ADDR = "0.0.0.0:8080"
	LISTEN_USER = "root"
	LISTEN_PASS = "root"
	TARGET_ADDR = "mysql:3306"
	TARGET_USER = "root"
	TARGET_PASS = "root"

	// Connection pool configuration
	MAX_IDLE_CONNS    = 50  // Maximum number of idle connections
	MIN_OPEN_CONNS    = 10  // Minimum number of open connections
	MAX_OPEN_CONNS    = 100 // Maximum number of open connections (reduced to prevent "Too many connections")
	MAX_CONN_LIFETIME = 300 // Maximum connection lifetime (seconds)
	CONN_TIMEOUT      = 60  // Connection timeout (seconds)
)

func main() {
	// Listen for connections on localhost port 3307
	l, err := net.Listen("tcp", LISTEN_ADDR)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	log.Printf("MySQL proxy server listening on %s", LISTEN_ADDR)

	// Accept connections in a loop
	for {
		// Accept a new connection
		c, err := l.Accept()
		if err != nil {
			log.Printf("Accept error: %v", err)
			continue
		}

		// Handle each connection in a goroutine
		handleListen(c)
	}
}

func handleListen(c net.Conn) {
	defer c.Close()

	handle := &QueryHandler{}

	if err := handle.UseDB(""); err != nil {
		log.Printf("Failed to initialize connection: %v", err)
		return
	}

	// Create a connection with user root and password root.
	conn, err := server.NewDefaultServer().NewConn(
		c,
		LISTEN_USER,
		LISTEN_PASS,
		handle,
	)
	if err != nil {
		log.Printf("Failed to create connection: %v", err)
		return
	}

	// Handle commands until the client disconnects
	for {
		if err := conn.HandleCommand(); err != nil {
			log.Printf("Command handling error: %v", err)
			break
		}
	}

	log.Printf("Connection closed from %s", c.RemoteAddr().String())
}
