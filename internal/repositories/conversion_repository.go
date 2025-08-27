package repositories

import (
	"context"
	"database/sql"
	"log"

	"github.com/MinhT933/file-converter/internal/domain"
)

type ConversionRepo struct {
	db *sql.DB
}

func NewConversionRepository(db *sql.DB) *ConversionRepo {
	return &ConversionRepo{
		db: db,
	}
}

func (r *ConversionRepo) Create(ctx context.Context, conversion *domain.Conversion) (string, error) {
	query := `INSERT INTO conversions (id, user_id, original_filename, converted_filename, status, created_at, expires_at, updated_at)
	Values($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.ExecContext(ctx, query,
		conversion.ConversionID, conversion.UserID, conversion.OriginalFilename, conversion.ConvertedFilename, conversion.Status,
		conversion.CreatedAt, conversion.ExpiresAt, conversion.UpdatedAt)
	if err != nil {
		log.Println("Error inserting conversion:", err)
		return "", err
	}
	return conversion.ConversionID, nil
}

func (r *ConversionRepo) FindByUserID(ctx context.Context, userID string) ([]domain.Conversion, error) {
	query := `SELECT id, user_id, original_file, converted_file, status, created_at, updated_at
		FROM conversions WHERE user_id = $1`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversions []domain.Conversion
	for rows.Next() {
		var conv domain.Conversion
		if err := rows.Scan(&conv.ConversionID, &conv.UserID, &conv.OriginalFilename,
			&conv.ConvertedFilename, &conv.CreatedAt, &conv.ExpiresAt); err != nil {
			return nil, err
		}
		conversions = append(conversions, conv)
	}

	return conversions, nil
}

func (r *ConversionRepo) FindByID(ctx context.Context, conversionID string) (*domain.Conversion, error) {
	return nil, nil // Placeholder for actual implementation
}

func (r *ConversionRepo) UpdateConversionStatus(ctx context.Context, conversionID string, status string) error {
	query := `UPDATE conversions SET status = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, conversionID)
	if err != nil {
		log.Println("🤔🤔🤔 Error updating conversion status:", err)
		return err
	}
	return nil
}
