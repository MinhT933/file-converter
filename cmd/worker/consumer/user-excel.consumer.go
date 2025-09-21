package consumer

import (
	"context"
	"encoding/json"
	"os"

	"github.com/MinhT933/file-converter/cmd/worker/jobs"
	"github.com/MinhT933/file-converter/internal/contextx"
	"github.com/MinhT933/file-converter/internal/template"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
)

type (
	UserExcelConsumer struct{}

	ExcelPayload struct {
		Path string `json:"path"`
	}
)

func NewUserExcelConsumer() *UserExcelConsumer {
	return &UserExcelConsumer{}
}

func (excel *UserExcelConsumer) Config() jetstream.ConsumerConfig {
	return jetstream.ConsumerConfig{
		Durable:       "user-excel-consumer",
		FilterSubject: "import.excel.user",
	}
}

func (excel UserExcelConsumer) Parse(data []byte) (ExcelPayload, error) {
	var payload ExcelPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		return ExcelPayload{}, err
	}
	return payload, nil
}

func (excel UserExcelConsumer) Handle(ctx context.Context, payload ExcelPayload) error {
	log := contextx.Logger(ctx)

	// Defensive checks: ensure path is provided and file exists
	if payload.Path == "" {
		log.Error("empty path in payload")
		return nil
	}

	// Normalize path: remove any leading ./ and use as-is for worker
	// (worker runs on same filesystem as server in this setup)
	cleanPath := payload.Path
	if len(cleanPath) >= 2 && cleanPath[:2] == "./" {
		cleanPath = cleanPath[2:]
	}

	if _, err := os.Stat(cleanPath); err != nil {
		log.Error("file not found for import", zap.String("path", cleanPath), zap.Error(err))
		return nil
	}

	job := jobs.ExcelUserJob{}
	results := template.RunImportWorkflow(ctx, job, cleanPath, 500)

	for res := range results {
		// Process the Excel file located at payload.Path
		if len(res.Errs) > 0 {
		log.Warn("Row failed", zap.Any("row", res.Row), zap.Any("errs", res.Errs))
			continue
	
		}
		log.Info("Row imported", zap.Any("user", res.Data))
	}
		return nil

}