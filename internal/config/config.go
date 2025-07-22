package config

import (
<<<<<<< HEAD
=======
	"fmt"
>>>>>>> e63175e (feat(auth): implement user authentication with email/password and social login)
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
<<<<<<< HEAD
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
=======
	AppName          string
	PortHTTP         string
	RedisAddr        string
	RedisPass        string
	MaxUploadMB      int
	SMTPHost         string
	SMTPPort         int
	SMTPUser         string
	SMTPPass         string
	EmailFrom        string
	FirebaseCredFile string
	FirebaseCredJSON string
>>>>>>> e63175e (feat(auth): implement user authentication with email/password and social login)
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
	// ✅ Debug biến môi trường
	credFile := os.Getenv("FIREBASE_CRED_FILE")
	fmt.Printf("🔍 FIREBASE_CRED_FILE: '%s'\n", credFile)

	// ✅ Kiểm tra file tồn tại
	if credFile != "" {
		if _, err := os.Stat(credFile); os.IsNotExist(err) {
			fmt.Printf("❌ File không tồn tại: %s\n", credFile)
		} else {
			fmt.Printf("✅ File tồn tại: %s\n", credFile)
		}
	}

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
		AppName:          getEnv("APP_NAME", "converter"),
		PortHTTP:         getEnv("PORT_HTTP", "8080"),
		RedisAddr:        getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPass:        os.Getenv("REDIS_PASSWORD"), // nếu muốn default thì dùng getEnv
		MaxUploadMB:      mb,
		SMTPHost:         getEnv("SMTP_HOST", ""),
		SMTPPort:         smtpPort,
		SMTPUser:         getEnv("SMTP_USER", ""),
		SMTPPass:         getEnv("SMTP_PASS", ""),
		EmailFrom:        getEnv("EMAIL_FROM", ""),
		FirebaseCredFile: os.Getenv("FIREBASE_CRED_FILE"),
		FirebaseCredJSON: os.Getenv("FIREBASE_CRED_JSON"),
	}
}
