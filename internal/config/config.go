package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/MinhT933/file-converter/internal/db"
	"github.com/spf13/viper"
)

type (
	Config struct {
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
		DB          *db.DatabaseConfig
		FirebaseCredFile string
		FirebaseCredJSON string
		Env         string
		CORSOrigins string
		Nats        struct {
			URL string
		}
		
	}
	ConfigWorker struct {
		App struct {
			Name       string `mapstructure:"name"`
			Env        string `mapstructure:"env"`
			WorkerType string `mapstructure:"workerType"`
		}

		Nats struct {
			URL     string `mapstructure:"url"`
			Subject string `mapstructure:"subject"`
		} `mapstructure:"nats"`

		DB *db.DatabaseConfig `mapstructure:"db"`

	}

)

// getEnv đọc biến môi trường, nếu không có thì trả về fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func Load() *Config {
	// load file .env nếu có
	// _ = godotenv.Load()

	// parse SMTP_PORT
	smtpPort, err := strconv.Atoi(getEnv("SMTP_PORT", "25"))
	if err != nil {
		log.Printf("⚠️  Invalid SMTP_PORT, using default 25")
		smtpPort = 25
	}

	return &Config{
		AppName:     getEnv("APP_NAME", "File Converter"),
		PortHTTP:    getEnv("PORT_HTTP", "8081"),
		RedisAddr:   getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPass:   getEnv("REDIS_PASS", ""),
		MaxUploadMB: 10,
		SMTPHost:    getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:    smtpPort,
		SMTPUser:    getEnv("SMTP_USER", ""),
		SMTPPass:    getEnv("SMTP_PASS", ""),
		EmailFrom:   getEnv("EMAIL_FROM", ""),
		DB:          db.LoadDatabase(),
		FirebaseCredFile: getEnv("FIREBASE_CRED_FILE", ""),
		FirebaseCredJSON: getEnv("FIREBASE_CRED_JSON", ""),
		Env:         getEnv("ENV", "dev"),
		CORSOrigins: getEnv("CORS_ORIGINS", "http://127.0.0.1:8081,http://localhost:8081,https://localhost:3000,http://localhost:3000,https://api-convert-file.minht.io.vn"),
		Nats:        struct{URL string}{URL: getEnv("NATS_URL", "nats://localhost:4222")},
	}
}

func LoadConfigWorker() ConfigWorker {
	var cfg ConfigWorker

	// Cho phép đọc ENV
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Provide a sensible default so worker doesn't crash when env is missing
	viper.SetDefault("app.workerType", "excel")

	// Bind env cụ thể (nếu muốn rõ ràng)
	_ = viper.BindEnv("nats.url", "NATS_URL")
	_ = viper.BindEnv("nats.stream", "NATS_STREAM")
	_ = viper.BindEnv("nats.subject", "NATS_SUBJECT")
	_ = viper.BindEnv("app.workerType", "WORKER_TYPE")
	_ = viper.BindEnv("app.env", "ENV")

	cfg.DB = db.LoadDatabase()

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Error unmarshalling config: %v", err)
	}

	// Warn if caller intentionally left WORKER_TYPE empty
	if cfg.App.WorkerType == "" {
		log.Printf("⚠️  WORKER_TYPE not set; using default 'excel'")
		cfg.App.WorkerType = viper.GetString("app.workerType")
	}
	return cfg
}
