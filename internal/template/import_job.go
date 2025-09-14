package template

import (
	"context"

	"github.com/MinhT933/file-converter/internal/excel"
)

type ImportJob[T any] interface {
	Parse(ctx context.Context, path string) (<-chan excel.ExcelRow, <-chan error)
	Transform(row []string) (T, error)
	Validate(data T) []error
	InsertBatch(ctx context.Context, data []T) error
	ReportError(row []string, errs []error)
}

type ImportResult[T any] struct {
	Row   []string
	Data  *T
	Errs  []error
	Stage string // parse | transform | validate | insert | success
}
