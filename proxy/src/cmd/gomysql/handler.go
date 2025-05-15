package main

import (
	"fmt"
	"log/slog"
	"regexp"
	"strings"

	"github.com/go-mysql-org/go-mysql/client"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/xwb1989/sqlparser"
)

var (
	// storedProcRegex matches stored procedure definitions
	storedProcRegex = regexp.MustCompile(`(?i)BEGIN\s+[\s\S]*?END`)
)

// QueryHandler processes MySQL queries and commands
type QueryHandler struct {
	client *client.Conn
}

// NewQueryHandler creates a new QueryHandler instance
func NewQueryHandler(client *client.Conn) *QueryHandler {
	if client == nil {
		slog.Error("MySQL client connection cannot be nil")
		panic("MySQL client connection cannot be nil")
	}
	return &QueryHandler{
		client: client,
	}
}

// splitQuery splits a query string into multiple query statements
func (h *QueryHandler) splitQuery(query string) ([]string, error) {
	if query == "" {
		return nil, nil
	}

	// Return stored procedures as is without splitting
	if storedProcRegex.MatchString(query) {
		return []string{query}, nil
	}

	// Split regular SQL queries into multiple statements
	queries, err := sqlparser.SplitStatementToPieces(query)
	if err != nil {
		return nil, fmt.Errorf("failed to split query: %w", err)
	}

	return queries, nil
}

// HandleQuery processes queries from clients
func (h *QueryHandler) HandleQuery(query string) (*mysql.Result, error) {
	slog.Info("Handling query", "query", query)

	// Split query into multiple statements
	queries, err := h.splitQuery(query)
	if err != nil {
		return nil, err
	}

	if len(queries) == 0 {
		return &mysql.Result{}, nil
	}

	// Execute each query in order
	var res *mysql.Result
	for _, q := range queries {
		trimQuery := strings.TrimSpace(q)
		if trimQuery == "" {
			continue
		}

		// Execute query
		res, err = h.client.Execute(trimQuery)
		if err != nil {
			return nil, fmt.Errorf("query execution error: %w", err)
		}
	}

	return res, nil
}

// UseDB processes database selection commands
func (h *QueryHandler) UseDB(dbName string) error {
	if dbName == "" {
		return fmt.Errorf("database name cannot be empty")
	}

	slog.Info("Using database", "db", dbName)
	if err := h.client.UseDB(dbName); err != nil {
		return fmt.Errorf("failed to use database '%s': %w", dbName, err)
	}

	return nil
}

// HandleFieldList processes FIELD_LIST commands
func (h *QueryHandler) HandleFieldList(table string, wildcard string) ([]*mysql.Field, error) {
	slog.Info("Handle field list request", "table", table, "wildcard", wildcard)
	return nil, fmt.Errorf("field list not supported")
}

// HandleStmtPrepare processes Prepare statements
func (h *QueryHandler) HandleStmtPrepare(query string) (int, int, interface{}, error) {
	slog.Info("Prepare statement request", "query", query)
	return 0, 0, nil, fmt.Errorf("prepared statements not supported")
}

// HandleStmtExecute processes the execution of Prepare statements
func (h *QueryHandler) HandleStmtExecute(context interface{}, query string, args []interface{}) (*mysql.Result, error) {
	slog.Info("Execute prepared statement", "query", query, "args", args)
	return nil, fmt.Errorf("prepared statement execution not supported")
}

// HandleStmtClose processes the closing of Prepare statements
func (h *QueryHandler) HandleStmtClose(context interface{}) error {
	slog.Info("Closing prepared statement")
	return nil
}

// HandleOtherCommand processes other commands
func (h *QueryHandler) HandleOtherCommand(cmd byte, data []byte) error {
	slog.Info("Handling other command", "cmd", cmd)
	return fmt.Errorf("command not supported")
}
