package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
	db *gorm.DB
}

// User構造体にgormタグを追加
type User struct {
	ID        int       `gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"size:100;not null"`
	Email     string    `gorm:"size:100;not null;unique"`
	Points    int       `gorm:"default:0"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func main() {
	// DB configuration
	config := DBConfig{
		Host:     "gomysql",
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
	}
	log.Printf("Found %d users in total", len(users))
	for _, u := range users {
		log.Printf("- ID=%d, Name=%s", u.ID, u.Name)
	}

	if err = db.TransferPoints(1, 2, 10); err != nil {
		log.Printf("Failed to transfer points: %v", err)
	}
}

// gormでDB接続
func NewDBHandler(config DBConfig) (*DBHandler, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		config.User, config.Password, config.Host, config.Port, config.DBName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("Database connection error: %w", err)
	}
	return &DBHandler{db: db}, nil
}

func (h *DBHandler) Close() error {
	sqlDB, err := h.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Create tables with AutoMigrate
func (h *DBHandler) EnsureTablesExist() error {
	if err := h.db.AutoMigrate(&User{}); err != nil {
		return fmt.Errorf("failed to migrate table: %w", err)
	}
	log.Println("Table users migrated or already exists")
	return nil
}

// CreateUser with gorm
func (h *DBHandler) CreateUser(name, email string) (int, error) {
	user := User{Name: name, Email: email}
	if err := h.db.Create(&user).Error; err != nil {
		return 0, fmt.Errorf("User creation error: %w", err)
	}
	return user.ID, nil
}

// GetUserByID with gorm
func (h *DBHandler) GetUserByID(id int) (*User, error) {
	var user User
	if err := h.db.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("User with ID %d does not exist", id)
		}
		return nil, fmt.Errorf("User retrieval error: %w", err)
	}
	return &user, nil
}

// GetAllUsers with gorm
func (h *DBHandler) GetAllUsers() ([]User, error) {
	var users []User
	if err := h.db.Find(&users).Error; err != nil {
		return nil, fmt.Errorf("User list retrieval error: %w", err)
	}
	return users, nil
}

// TransferPoints with gorm transaction
func (h *DBHandler) TransferPoints(fromUserID, toUserID, points int) error {
	return h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&User{}).Where("id = ?", fromUserID).UpdateColumn("points", gorm.Expr("points - ?", points)).Error; err != nil {
			return fmt.Errorf("Point deduction error: %w", err)
		}
		if err := tx.Model(&User{}).Where("id = ?", toUserID).UpdateColumn("points", gorm.Expr("points + ?", points)).Error; err != nil {
			return fmt.Errorf("Point addition error: %w", err)
		}
		return nil
	})
}
