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
	SMTPHost    string
	SMTPPort    int
	SMTPUser    string
	SMTPPass    string
	EmailFrom   string
	DB          *DBconfig
}

// getEnv đọc biến môi trường, nếu không có thì trả về fallback
func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}

func Load() *Config {
	// load file .env nếu có
	_ = godotenv.Load()

	// parse MAX_UPLOAD_MB
	mb, err := strconv.Atoi(getEnv("MAX_UPLOAD_MB", "10"))
	if err != nil {
		mb = 10
	}

	// parse SMTP_PORT
	smtpPort, err := strconv.Atoi(getEnv("SMTP_PORT", "25"))
	if err != nil {
		smtpPort = 25
	}

	return &Config{
		AppName:     getEnv("APP_NAME", "converter"),
		PortHTTP:    getEnv("PORT_HTTP", "8080"),
		RedisAddr:   getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPass:   os.Getenv("REDIS_PASSWORD"), // nếu muốn default thì dùng getEnv
		MaxUploadMB: mb,
		SMTPHost:    getEnv("SMTP_HOST", ""),
		SMTPPort:    smtpPort,
		SMTPUser:    getEnv("SMTP_USER", ""),
		SMTPPass:    getEnv("SMTP_PASS", ""),
		EmailFrom:   getEnv("EMAIL_FROM", ""),

		DB: LoadDB(),
	}
}
