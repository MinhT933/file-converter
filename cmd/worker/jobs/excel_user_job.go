package jobs

import (
	"context"
	"fmt"
	"time"

	"github.com/MinhT933/file-converter/internal/db"
	"github.com/MinhT933/file-converter/internal/excel"
	"github.com/jackc/pgx/v5"
)

type (
	ExcelUserJob struct{}

	User struct {
		ID int  `excel:"ID"`
		Name string `excel:"Name"`
		Email string `excel:"Email"`
		Age int `excel:"Age"`
		Active bool `excel:"Active"`
		Joined time.Time `excel:"Joined" time_format:"2006-01-02"`
	}
)

func (ExcelUserJob) Parse(ctx context.Context, path string) (<-chan excel.ExcelRow, <-chan error) {
	return excel.ReadExcelStream(ctx, path, "Sheet1")
}

func (ExcelUserJob) Transform(row []string) (User, error) {
  var u User 

  headerMap := excel.BuildHeaderIndex([]string{"ID", "Name", "Email", "Age", "Active", "Joined"})
  if err := excel.MapRowByHeader(row, headerMap, &u); err != nil {
	return User{}, err
  }
  return u, nil
}

func (ExcelUserJob) Validate(data User) []error {
	var errs []error
	if data.ID <= 0 {
		errs = append(errs, fmt.Errorf("invalid ID"))
	}
	if data.Name == "" {
		errs = append(errs, fmt.Errorf("name is required"))
	}
	return errs
}

func (ExcelUserJob) InsertBatch(ctx context.Context, data []User) error {
	if len(data) == 0 {
		return nil
	}

	rows := make([][]interface{}, len(data))
	for i, user := range data {
		rows[i] = []interface{}{user.ID, user.Name, user.Email, user.Age, user.Active, user.Joined}
	}

	_, err := db.Pool.CopyFrom(
		ctx,
		pgx.Identifier{"users"},
		[]string{"id", "name", "email", "age", "active", "joined"},
		pgx.CopyFromRows(rows),
	
	)
	return err


}

func (ExcelUserJob) ReportError(row []string, errs []error) {
	fmt.Println("Row failed:", row, "Errors:", errs)
}
