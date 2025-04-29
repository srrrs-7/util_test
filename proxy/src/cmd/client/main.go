package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Database connection information
type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

// DB handler
type DBHandler struct {
	db *sql.DB
}

func main() {
	// DB configuration
	config := DBConfig{
		Host:     "proxy",
		Port:     8080,
		User:     "root",
		Password: "root",
		DBName:   "test",
	}

	// Initialize DB handler
	db, err := NewDBHandler(config)
	if err != nil {
		log.Fatalf("DB handler initialization error: %v", err)
	}
	defer db.Close()

	// Ensure tables exist
	if err := db.EnsureTablesExist(); err != nil {
		log.Fatalf("テーブル確認エラー: %v", err)
	}

	// Example: Create user
	userID, err := db.CreateUser("Taro Tanaka", "tanaka@example.com")
	if err != nil {
		log.Printf("Failed to create user: %v", err)
	} else {
		log.Printf("User created. ID: %d", userID)
	}

	// Example: Get user
	user, err := db.GetUserByID(1)
	if err != nil {
		log.Printf("Failed to get user: %v", err)
	} else {
		log.Printf("User: ID=%d, Name=%s, Email=%s", user.ID, user.Name, user.Email)
	}

	// Example: Get all users
	users, err := db.GetAllUsers()
	if err != nil {
		log.Printf("Failed to get user list: %v", err)
	} else {
		log.Printf("Found %d users in total", len(users))
		for _, u := range users {
			log.Printf("- ID=%d, Name=%s", u.ID, u.Name)
		}
	}
}

// Create a new DB handler
func NewDBHandler(config DBConfig) (*DBHandler, error) {
	// Create connection string
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		config.User, config.Password, config.Host, config.Port, config.DBName)

	// Connect to DB
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("Database connection error: %w", err)
	}

	// Connection test
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("Database connection check error: %w", err)
	}

	// Connection pool settings
	db.SetMaxOpenConns(25)                 // Maximum number of connections
	db.SetMaxIdleConns(25)                 // Maximum number of idle connections
	db.SetConnMaxLifetime(5 * time.Minute) // Maximum lifetime of a connection

	return &DBHandler{db: db}, nil
}

// Close DB handler
func (h *DBHandler) Close() error {
	return h.db.Close()
}

// Sample data structure
type User struct {
	ID        int
	Name      string
	Email     string
	CreatedAt time.Time
}

func (h *DBHandler) EnsureTablesExist() error {
	// users テーブルを作成
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		email VARCHAR(100) NOT NULL UNIQUE,
		points INT DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)
	`

	_, err := h.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	log.Println("Table users created or already exists")
	return nil
}

// Get user by ID
func (h *DBHandler) GetUserByID(id int) (*User, error) {
	var user User
	query := "SELECT id, name, email, created_at FROM users WHERE id = ?"

	err := h.db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("User with ID %d does not exist", id)
		}
		return nil, fmt.Errorf("User retrieval error: %w", err)
	}

	return &user, nil
}

// Create a new user
func (h *DBHandler) CreateUser(name, email string) (int, error) {
	query := "INSERT INTO users (name, email, created_at) VALUES (?, ?, ?)"

	result, err := h.db.Exec(query, name, email, time.Now())
	if err != nil {
		return 0, fmt.Errorf("User creation error: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("Insert ID retrieval error: %w", err)
	}

	return int(id), nil
}

// Get all users
func (h *DBHandler) GetAllUsers() ([]User, error) {
	query := "SELECT id, name, email, created_at FROM users"

	rows, err := h.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("User list retrieval error: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt); err != nil {
			return nil, fmt.Errorf("User data reading error: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Row iteration error: %w", err)
	}

	return users, nil
}

// Transaction example
func (h *DBHandler) TransferPoints(fromUserID, toUserID, points int) error {
	// Start transaction
	tx, err := h.db.Begin()
	if err != nil {
		return fmt.Errorf("Transaction start error: %w", err)
	}

	// Ensure rollback (disabled upon commit)
	defer tx.Rollback()

	// Deduct points from fromUser
	_, err = tx.Exec("UPDATE users SET points = points - ? WHERE id = ?", points, fromUserID)
	if err != nil {
		return fmt.Errorf("Point deduction error: %w", err)
	}

	// Add points to toUser
	_, err = tx.Exec("UPDATE users SET points = points + ? WHERE id = ?", points, toUserID)
	if err != nil {
		return fmt.Errorf("Point addition error: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("Transaction commit error: %w", err)
	}

	return nil
}
