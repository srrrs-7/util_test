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
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/go-mysql-org/go-mysql/client"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/server"
	"github.com/xwb1989/sqlparser"
)

// Constants for database types
const (
	DBTypeProxy = "proxy"
	DBTypeTest  = "test"
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
)

// Default timeout for graceful shutdown
const (
	defaultShutdownTimeout = 30 * time.Second
)

var atomPool atomic.Pointer[client.Pool]

// Env is a struct that holds database connection information loaded from environment variables
type Env struct {
	proxyAddr  string
	proxyUser  string
	proxyPass  string
	testDBAddr string
	testDBUser string
	testDBPass string
	testDBName string
}

// newEnv loads settings from environment variables and generates an Env struct
func newEnv() Env {
	return Env{
		proxyAddr:  os.Getenv(EnvProxyAddr),
		proxyUser:  os.Getenv(EnvProxyUser),
		proxyPass:  os.Getenv(EnvProxyPass),
		testDBAddr: os.Getenv(EnvTestDBAddr),
		testDBUser: os.Getenv(EnvTestDBUser),
		testDBPass: os.Getenv(EnvTestDBPass),
		testDBName: os.Getenv(EnvTestDBName),
	}
}

// DBConfig is a struct that holds database connection configuration information
type DBConfig struct {
	DBName string
	Addr   string
	User   string
	Pass   string
}

// Creates a mapping of all database configurations loaded from environment variables
func createDBConfigMap(env Env) map[string]DBConfig {
	return map[string]DBConfig{
		DBTypeProxy: {
			Addr: env.proxyAddr,
			User: env.proxyUser,
			Pass: env.proxyPass,
		},
		DBTypeTest: {
			Addr:   env.testDBAddr,
			User:   env.testDBUser,
			Pass:   env.testDBPass,
			DBName: env.testDBName,
		},
	}
}

// ProxyServer represents the MySQL proxy server
type ProxyServer struct {
	listener   net.Listener
	dbConfMap  map[string]DBConfig
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
	shutdownCh chan struct{}
}

// NewProxyServer creates a new ProxyServer instance
func NewProxyServer(dbConfMap map[string]DBConfig) *ProxyServer {
	ctx, cancel := context.WithCancel(context.Background())
	return &ProxyServer{
		dbConfMap:  dbConfMap,
		ctx:        ctx,
		cancel:     cancel,
		shutdownCh: make(chan struct{}),
	}
}

// Start starts the MySQL proxy server
func (s *ProxyServer) Start() error {
	proxyAddr := s.dbConfMap[DBTypeProxy].Addr
	l, err := net.Listen("tcp", proxyAddr)
	if err != nil {
		return fmt.Errorf("failed to listen on address: %w", err)
	}
	s.listener = l
	log.Printf("MySQL proxy server successfully listening on %s", proxyAddr)

	go s.handleSignals()

	pool, err := client.NewPoolWithOptions(
		s.dbConfMap[DBTypeTest].Addr,
		s.dbConfMap[DBTypeTest].User,
		s.dbConfMap[DBTypeTest].Pass,
		s.dbConfMap[DBTypeTest].DBName,
		client.WithPoolLimits(16, 64, 32),
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

	// Create MySQL server connection
	serverConn, err := server.NewDefaultServer().NewConn(
		conn,
		s.dbConfMap[DBTypeProxy].User,
		s.dbConfMap[DBTypeProxy].Pass,
		NewQueryHandler(client),
	)
	if err != nil {
		log.Printf("Failed to create server connection: %v", err)
		return
	}

	// Process commands until client disconnects
	for {
		if err := serverConn.HandleCommand(); err != nil {
			if err.Error() == "connection closed" {
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
	case <-time.After(defaultShutdownTimeout):
		log.Println("Server shutdown timed out")
	}
}

// main is the entry point of the application
func main() {
	// Load settings from environment variables
	env := newEnv()
	log.Printf("Starting MySQL proxy server with environment configuration")

	// Create database configuration map
	dbConfMap := createDBConfigMap(env)

	// Create and start the proxy server
	proxyServer := NewProxyServer(dbConfMap)
	if err := proxyServer.Start(); err != nil {
		log.Panic("Failed to start proxy server: ", err)
	}
}

// QueryHandler is a custom handler for processing MySQL queries
type QueryHandler struct {
	client *client.Conn
}

// NewQueryHandler creates a new QueryHandler instance
func NewQueryHandler(client *client.Conn) *QueryHandler {
	return &QueryHandler{
		client: client,
	}
}

// splitQuery splits a query string into multiple query statements
func (h *QueryHandler) splitQuery(query string) ([]string, error) {
	// Return stored procedures as is without splitting
	if regexp.MustCompile(`(?i)BEGIN\s+[\s\S]*?END`).MatchString(query) {
		return []string{query}, nil
	}

	// Split regular SQL queries into multiple statements
	queries, err := sqlparser.SplitStatementToPieces(query)
	if err != nil {
		log.Printf("Failed to split query: %v", err)
		return nil, err
	}

	return queries, nil
}

// HandleQuery processes queries from clients
func (h *QueryHandler) HandleQuery(query string) (*mysql.Result, error) {
	log.Printf("Handling query: %s", query)

	// Split query into multiple statements
	queries, err := h.splitQuery(query)
	if err != nil {
		log.Printf("Failed to split query: %v", err)
		return nil, err
	}

	// Execute each query in order
	var res *mysql.Result
	for _, q := range queries {
		trimQuery := strings.Trim(q, " ")
		if trimQuery == "" || trimQuery == "\n" {
			continue
		}

		// Execute query
		res, err = h.client.Execute(q)
		if err != nil {
			log.Printf("Query execution error: %v", err)
			return nil, err
		}
	}

	return res, nil
}

// UseDB processes database selection commands
func (h *QueryHandler) UseDB(dbName string) error {
	log.Printf("Using database: %s", dbName)

	return h.client.UseDB(dbName)
}

// HandleFieldList processes FIELD_LIST commands
func (h *QueryHandler) HandleFieldList(table string, wildcard string) ([]*mysql.Field, error) {
	log.Printf("Handle field list request for table: %s, wildcard: %s", table, wildcard)
	return nil, fmt.Errorf("field list not supported")
}

// HandleStmtPrepare processes Prepare statements
func (h *QueryHandler) HandleStmtPrepare(query string) (int, int, interface{}, error) {
	log.Printf("Prepare statement request: %s", query)
	return 0, 0, nil, fmt.Errorf("prepared statements not supported")
}

// HandleStmtExecute processes the execution of Prepare statements
func (h *QueryHandler) HandleStmtExecute(context interface{}, query string, args []interface{}) (*mysql.Result, error) {
	log.Printf("Execute prepared statement: %s with args: %v", query, args)
	return nil, fmt.Errorf("prepared statement execution not supported")
}

// HandleStmtClose processes the closing of Prepare statements
func (h *QueryHandler) HandleStmtClose(context interface{}) error {
	log.Printf("Closing prepared statement")
	return nil
}

// HandleOtherCommand processes other commands
func (h *QueryHandler) HandleOtherCommand(cmd byte, data []byte) error {
	log.Printf("Handling other command: %d", cmd)
	return fmt.Errorf("command not supported")
}
