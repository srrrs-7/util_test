package main

import (
	"log"
	"log/slog"
	"net"
	"proxy/config"
	"time"

	"github.com/go-mysql-org/go-mysql/client"
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
	conf := config.NewConfig()

	// Listen for connections on localhost port 3307
	l, err := net.Listen("tcp", conf.Server.Addr)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	log.Printf("MySQL proxy server listening on %s", conf.Server.Addr)

	// Accept connections in a loop
	for {
		// Accept a new connection
		c, err := l.Accept()
		if err != nil {
			log.Printf("Accept error: %v", err)
			continue
		}

		go func(c net.Conn, conf config.Config) {
			// Handle each connection in a goroutine
			go handleListen(c, conf)
		}(c, conf)
	}
}

func handleListen(c net.Conn, conf config.Config) {
	defer c.Close()

	// Establish a new connection
	pool, err := client.NewPoolWithOptions(
		conf.Database.Test.Addr,
		conf.Database.Test.User,
		conf.Database.Test.Pass,
		conf.Database.Test.Name,
		client.WithLogger(slog.Default()),
		client.WithNewPoolPingTimeout(
			time.Duration(conf.Setting.ConnLifetime)*time.Second,
		),
		client.WithPoolLimits(
			conf.Setting.MinOpenConns,
			conf.Setting.MaxIdleConns,
			conf.Setting.MaxOpenConns,
		),
	)
	if err != nil {
		panic(err.Error())
	}

	handle := &QueryHandler{
		config: conf,
		pool:   pool,
	}

	// Create a connection with user root and password root.
	conn, err := server.NewDefaultServer().NewConn(
		c,
		conf.Server.User,
		conf.Server.Pass,
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
