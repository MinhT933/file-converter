package main

import (
	"fmt"
	"os"
	"strconv"
)

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

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

func main() {
	fmt.Println("🐳 Testing DB Config in Container")
	fmt.Println("=================================")

	config := LoadDB()
	fmt.Printf("  Host:     '%s'\n", config.Host)
	fmt.Printf("  Port:     %d\n", config.Port)
	fmt.Printf("  User:     '%s'\n", config.User)
	fmt.Printf("  Password: '%s'\n", config.Password)
	fmt.Printf("  Name:     '%s'\n", config.Name)
	fmt.Printf("  MaxIdle:  %d\n", config.MaxIdle)
	fmt.Printf("  MaxOpen:  %d\n", config.MaxOpen)

	fmt.Println("\n🔍 Raw Environment Variables:")
	envVars := []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_MAX_IDLE", "DB_MAX_OPEN"}
	for _, env := range envVars {
		value := os.Getenv(env)
		if value == "" {
			value = "[NOT SET]"
		}
		fmt.Printf("  %s=%s\n", env, value)
	}
}
