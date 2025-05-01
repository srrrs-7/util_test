package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"regexp"
	"strings"
	"time"

	"github.com/go-mysql-org/go-mysql/client"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/xwb1989/sqlparser"
)

var (
	poolMap = make(map[string]*client.Pool)
)

// QueryHandler is a custom handler for processing simple MySQL queries
type QueryHandler struct {
	pool *client.Pool
}

// initConnection establishes the initial connection to the target database
func (h *QueryHandler) initConnection(dbName string) error {
	// Establish a new connection
	connPool, err := client.NewPoolWithOptions(
		TARGET_ADDR,
		TARGET_USER,
		TARGET_PASS,
		dbName,
		client.WithLogger(slog.Default()),
		client.WithNewPoolPingTimeout(CONN_TIMEOUT*time.Second),
		client.WithPoolLimits(MIN_OPEN_CONNS, MAX_IDLE_CONNS, MAX_OPEN_CONNS),
	)
	if err != nil {
		log.Printf("Failed to create connection pool: %v", err)
		return err
	}

	poolMap[dbName] = connPool
	h.pool = connPool

	return nil
}

func (h *QueryHandler) splitQuery(query string) ([]string, error) {
	// stored procedure regex
	if regexp.MustCompile(`(?i)BEGIN\s+[\s\S]*?END`).MatchString(query) {
		return []string{query}, nil
	}

	queries, err := sqlparser.SplitStatementToPieces(query)
	if err != nil {
		log.Printf("Failed to split query: %v", err)
		return nil, err
	}

	return queries, nil
}

// HandleQuery processes queries from clients
func (h *QueryHandler) HandleQuery(query string) (*mysql.Result, error) {
	log.Printf("Handle query: %s", query)

	queries, err := h.splitQuery(query)
	if err != nil {
		log.Printf("Failed to split query: %v", err)
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(CONN_TIMEOUT)*time.Second)
	defer cancel()
	conn, err := h.pool.GetConn(ctx)
	if err != nil {
		log.Printf("Failed to get connection: %v", err)
		return nil, err
	}

	var res *mysql.Result
	for _, q := range queries {
		if q = strings.Trim(q, " "); q == "\n" || q == "" {
			continue
		}

		res, err = conn.Execute(query)
		h.pool.PutConn(conn)
		if err != nil {
			log.Printf("Failed to execute query: %v", err)
			return nil, err
		}
	}

	return res, nil
}

// UseDB processes database selection commands
func (h *QueryHandler) UseDB(dbName string) error {
	log.Println("UseDB:", dbName)

	if poolMap[dbName] != nil {
		return nil
	}

	if err := h.initConnection(dbName); err != nil {
		log.Printf("Failed to switch database: %v", err)
		return err
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
