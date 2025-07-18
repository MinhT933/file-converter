// internal/converter/converter.go
package converter

import (
	"context"
	"os"
	"path/filepath"
   "errors"
   "io"

)
var Registry = map[string]Converter{}
var ErrUnsupportedFormat = errors.New("unsupported file format")
// ConvertFile là hàm tiện lợi để dùng plugin theo đuôi 


type Converter interface {
	Convert(ctx context.Context, r io.Reader, w io.Writer) error
}

func ConvertFile(inputPath, outputPath string) error {
	ext := filepath.Ext(inputPath)
	slug := ""

	switch ext {
	case ".html":
		slug = "html_pdf"
	// thêm các định dạng khác nếu cần
	default:
		return ErrUnsupportedFormat // bạn định nghĩa lỗi này nếu muốn
	}

	converter, ok := Registry[slug]
	if !ok {
		return ErrUnsupportedFormat
	}

	in, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer out.Close()

	return converter.Convert(context.Background(), in, out)
}
