package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/go-mysql-org/go-mysql/client"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/xwb1989/sqlparser"
)

var storedProcRegex = regexp.MustCompile(`(?i)BEGIN\s+[\s\S]*?END`)

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
	if storedProcRegex.MatchString(query) {
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
		trimQuery := strings.TrimSpace(q)
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
