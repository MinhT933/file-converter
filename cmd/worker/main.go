package main

import (
	"context"
	"encoding/json"
	"log"
	"path/filepath"

	"github.com/hibiken/asynq"
	"github.com/joho/godotenv"

	"github.com/MinhT933/file-converter/internal/config"
	"github.com/MinhT933/file-converter/internal/converter"
	"github.com/MinhT933/file-converter/internal/infra/email"
	"github.com/MinhT933/file-converter/internal/tasks"
)

func main() {
	log.Println("[DEBUG] Bắt đầu khởi động worker ✅")
	godotenv.Load()
	cfg := config.Load()

	// client để enqueue tiếp email task
	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPass,
	})

	defer client.Close()

	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: cfg.RedisAddr, Password: cfg.RedisPass},
		asynq.Config{Concurrency: 10, Queues: map[string]int{"default": 10}},
	)

	emailSvc := email.NewSMTPProvider(
		cfg.SMTPHost, cfg.SMTPPort,
		cfg.SMTPUser, cfg.SMTPPass,
		cfg.EmailFrom,
	)

	mux := asynq.NewServeMux()

	// Handler cho file:convert
	mux.HandleFunc(tasks.TypeConvertFile, func(ctx context.Context, t *asynq.Task) error {
		var p struct {
			Email      string `json:"email"`
			InputPath  string `json:"input_path"`
			OutputPath string `json:"output_path"`
		}
		if err := json.Unmarshal(t.Payload(), &p); err != nil {
			return err
		}

		// 1) Thực hiện convert (ví dụ HTML→PDF, DOCX→PDF,…)
		if err := converter.ConvertFile(p.InputPath, p.OutputPath); err != nil {
			return err
		}

		// 2) Tạo link download
		downloadURL := "https://127.0.0.1:8080/downloads/" + filepath.Base(p.OutputPath)

		// 3) Enqueue task gửi email
		emailPayload, _ := json.Marshal(map[string]string{
			"email":    p.Email,
			"file_url": downloadURL,
		})

		_, err := client.Enqueue(asynq.NewTask(tasks.TypeEmailNotify, emailPayload))
		return err

	})

	// Handler cho email:notify
	mux.HandleFunc(tasks.TypeEmailNotify, func(ctx context.Context, t *asynq.Task) error {
		log.Printf("zô email notify")
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[ERROR] Panic trong xử lý task email:notify: %v", r)
			}
		}()
		var p struct {
			Email   string `json:"email"`
			FileURL string `json:"file_url"`
		}
		if err := json.Unmarshal(t.Payload(), &p); err != nil {
			log.Printf("Error unmarshalling task payload: %v", err)
			return err
		}
		return emailSvc.SendConversionEmail(ctx, p.Email, p.FileURL)
	})

	if err := srv.Run(mux); err != nil {
		log.Fatalf("asynq server error: %v", err)
	}
}
