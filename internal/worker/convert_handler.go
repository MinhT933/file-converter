package worker

import (
	"context"
	"encoding/json"
	"github.com/hibiken/asynq"
)

type PayloadConvert struct {
	Src string `json:"src"` // đường dẫn file gốc
}

func ConvertHandler() asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		var p PayloadConvert
		if err := json.Unmarshal(t.Payload(), &p); err != nil {
			return err
		}

		return nil
	}
}
