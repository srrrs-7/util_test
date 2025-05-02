package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"proxy/config"
	"sync"
	"time"

	"github.com/go-mysql-org/go-mysql/client"
	"github.com/go-mysql-org/go-mysql/mysql"
)

var (
	poolMap map[string]*client.Pool
)

// QueryHandler is a custom handler for processing simple MySQL queries
type QueryHandler struct {
	config config.Config
	pool   *client.Pool
	txConn *client.Conn
	mu     sync.Mutex
}

// initConnection establishes the initial connection to the target database
func (h *QueryHandler) initConnection(dbName string) error {
	log.Printf("Initializing connection to database: %s", dbName)

	if poolMap[dbName] != nil {
		return nil
	}

	// Establish a new connection
	connPool, err := client.NewPoolWithOptions(
		TARGET_ADDR,
		TARGET_USER,
		TARGET_PASS,
		dbName,
		client.WithLogger(slog.Default()),
		client.WithNewPoolPingTimeout(
			time.Duration(h.config.Setting.ConnLifetime)*time.Second,
		),
		client.WithPoolLimits(
			h.config.Setting.MinOpenConns,
			h.config.Setting.MaxIdleConns,
			h.config.Setting.MaxOpenConns,
		),
	)
	if err != nil {
		log.Printf("Failed to create connection pool: %v", err)
		return err
	}

	poolMap[dbName] = connPool
	h.pool = connPool

	return nil
}

// HandleQuery processes queries from clients
func (h *QueryHandler) HandleQuery(query string) (*mysql.Result, error) {
	log.Printf("Handle query: %s", query)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(h.config.Setting.ConnLifetime)*time.Second,
	)
	defer cancel()
	conn, err := h.pool.GetConn(ctx)
	if err != nil {
		log.Printf("Failed to get connection: %v", err)
		return nil, err
	}

	res, err := conn.ExecuteMultiple(query, func(result *mysql.Result, err error) {
		if err != nil {
			log.Printf("Error executing query: %v", err)
			return
		}
	})
	defer h.pool.PutConn(conn)
	if err != nil {
		log.Printf("Failed to execute query: %v", err)
		return nil, err
	}

	return res, nil
}

// UseDB processes database selection commands
func (h *QueryHandler) UseDB(dbName string) error {
	log.Printf("Switching to database: %s", dbName)

	h.mu.Lock()
	defer h.mu.Unlock()

	if err := h.initConnection(dbName); err != nil {
		log.Printf("Failed to switch database: %v", err)
		return err
	}

	if poolMap[dbName] == nil {
		h.pool = poolMap[dbName]
		return nil
	}

	return nil
}

// HandleFieldList processes FIELD_LIST commands
func (h *QueryHandler) HandleFieldList(table string, wildcard string) ([]*mysql.Field, error) {
	log.Printf("Handle field list: %s, %s", table, wildcard)
	return nil, fmt.Errorf("not supported")
}

// HandleStmtPrepare processes Prepare statements
func (h *QueryHandler) HandleStmtPrepare(query string) (int, int, interface{}, error) {
	log.Printf("Prepare statement: %s", query)
	return 0, 0, nil, fmt.Errorf("not supported")
}

// HandleStmtExecute processes the execution of Prepare statements
func (h *QueryHandler) HandleStmtExecute(context interface{}, query string, args []interface{}) (*mysql.Result, error) {
	log.Printf("Execute statement: %s with args: %v", query, args)
	return nil, fmt.Errorf("not supported")
}

// HandleStmtClose processes the closing of Prepare statements
func (h *QueryHandler) HandleStmtClose(context interface{}) error {
	return nil
}

// HandleOtherCommand processes other commands
func (h *QueryHandler) HandleOtherCommand(cmd byte, data []byte) error {
	log.Printf("Handle other command: %d", cmd)
	return fmt.Errorf("not supported now")
}
