package excel

import (
	"context"
	"fmt"
	"os"

	"github.com/MinhT933/file-converter/internal/contextx"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
)

type ExcelRow []string

func ReadExcelStream(ctx context.Context, path string, sheetName string) (<-chan ExcelRow, <-chan error) {
	log := contextx.Logger(ctx)
	output := make(chan ExcelRow)
	errCh := make(chan error, 1)


	file, err := os.Open(path)
	if err != nil {
		log.Error("Failed to open file", zap.String("path", path), zap.Error(err))
		errCh <- fmt.Errorf("❌ failed to open file: %w", err)
		return output, errCh
	}

	go func() {
		defer close(output)
		defer close(errCh)
		defer file.Close()

		// Simulate reading Excel rows
        f, err := excelize.OpenReader(file)
		if err != nil {
			errCh <- err
			return
		}
		defer f.Close()

		rows, err := f.Rows(sheetName)
		if err != nil {
			errCh <- fmt.Errorf("❌ failed to get rows: %w", err)
			return
		}
		firstRow := true
		for rows.Next() {
			cols, err := rows.Columns()
			if err != nil {
				errCh <- fmt.Errorf("❌ failed to read row: %w", err)
				return
			}
			if firstRow {
				firstRow = false
				return
			}
			output <- cols
		}
	}()

	return output, errCh
}