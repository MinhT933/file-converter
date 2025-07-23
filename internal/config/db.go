package config

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
)

type DBconfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
	MaxIdle  int    `json:"maxidle"`
	MaxOpen  int    `json:"maxopen"`
}

func LoadDB() *DBconfig {
	atoi := func(key, def string) int {
		n, _ := strconv.Atoi(getEnv(key, def))
		return n
	}

	return &DBconfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     atoi("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", ""),
		Name:     getEnv("DB_NAME", "file_converter"),
		MaxIdle:  atoi("DB_MAX_IDLE", "5"),
		MaxOpen:  atoi("DB_MAX_OPEN", "100"),
	}
}

// TestConnection - Test database connection and return *sql.DB
func (c *DBconfig) TestConnection() (*sql.DB, error) {
	// Build connection string
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.Name)

	// Open database connection
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxIdleConns(c.MaxIdle)
	db.SetMaxOpenConns(c.MaxOpen)
	db.SetConnMaxLifetime(time.Hour)

	// Test the connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// 🎉 Success message
	log.Printf("✅ Database connected successfully!")
	log.Printf("📊 Connection details: %s:%d/%s (MaxIdle:%d, MaxOpen:%d)",
		c.Host, c.Port, c.Name, c.MaxIdle, c.MaxOpen)

	return db, nil
}

// ConnectDB - Load config and establish database connection
func ConnectDB() (*sql.DB, error) {
	config := LoadDB()

	log.Printf("🔄 Connecting to database %s:%d/%s...",
		config.Host, config.Port, config.Name)

	db, err := config.TestConnection()
	if err != nil {
		log.Printf("❌ Database connection failed: %v", err)
		return nil, err
	}

	return db, nil
}
