package htmlwk

import (
	"context"
	"io"
	"os"
	"path/filepath"

	wkhtml "github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/MinhT933/file-converter/internal/converter"
)

type HTML2PDF struct{}

func (HTML2PDF) Convert(ctx context.Context, r io.Reader, w io.Writer) error {
	// 1. Tạo thư mục tạm cho job
	tmpDir, err := os.MkdirTemp("", "job-*")
	if err != nil { return err }
	defer os.RemoveAll(tmpDir) // dọn sau khi xong

	// 2. Ghi toàn bộ HTML ra file index.html
	htmlPath := filepath.Join(tmpDir, "index.html")
	if err := writeAll(htmlPath, r); err != nil {
		return err
	}
	

	// 3. Khởi tạo wkhtmltopdf
	pdfg, err := wkhtml.NewPDFGenerator()
	if err != nil { return err }

	// 4. Thêm trang bằng ĐƯỜNG DẪN FILE
	page := wkhtml.NewPage(htmlPath)
	page.EnableLocalFileAccess.Set(true) // cho đọc file://

	pdfg.Dpi.Set(300)
	pdfg.AddPage(page)

	if err := pdfg.Create(); err != nil {
        return err              // trả lỗi rõ ràng (hết treo)
    }

	// 6. Ghi PDF ra writer
	_, err = w.Write(pdfg.Bytes())
	return err
}

func writeAll(path string, src io.Reader) error {
	f, err := os.Create(path)
	if err != nil { return err }
	defer f.Close()
	_, err = io.Copy(f, src)
	return err
}

func init() {
	converter.Registry["html_pdf"] = HTML2PDF{}
}
