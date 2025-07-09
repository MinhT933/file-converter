package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName     string
	PortHTTP    string
	RedisAddr   string
	RedisPass   string
	MaxUploadMB int
}

func Load() *Config {
	_ = godotenv.Load()

	mb, _ := strconv.Atoi(getEnv("MAX_UPLOAD_MB", "10"))
	return &Config{
		AppName:     getEnv("APP_NAME", "converter"),
		PortHTTP:    getEnv("PORT_HTTP", "8080"),
		RedisAddr:   getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPass:   os.Getenv("REDIS_PASSWORD"),
		MaxUploadMB: mb,
	}
}

func getEnv(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
