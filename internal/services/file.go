package services

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/MinhT933/file-converter/internal/domain"
)

type FileService struct {
	conversionRepo domain.ConversionRepository
}

func NewFileService(conversionRepo domain.ConversionRepository) *FileService {
	return &FileService{
		conversionRepo: conversionRepo,
	}
}
func (s *FileService) SaveConvertedFile(ctx context.Context, UserID string, conversion *domain.Conversion) (string, string, error) {
	fullPath, err := saveFileToStorage(UserID, conversion.OriginalFilename, conversion.ConvertedFilename)
	if err != nil {
		return "", "", fmt.Errorf("failed to save file: %w", err)
	}

	_, err = s.conversionRepo.Create(ctx, conversion)
	if err != nil {
		return "", "", fmt.Errorf("failed to save conversion record: %w", err)
	}

	return fullPath, conversion.ConversionID, nil
}

func saveFileToStorage(userID, originalName, convertedName string) (string, error) {
	// tạo thư mục lưu trữ nếu chưa tồn tại
	outputDir := "./storage/" + userID
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	fullPath := filepath.Join(outputDir, convertedName)
	if err := os.WriteFile(fullPath, []byte{}, 0644); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	return fullPath, nil
}

func (s *FileService) UpdateConversionStatus(ctx context.Context, conversionID string, status string) error {
	return s.conversionRepo.UpdateConversionStatus(ctx, conversionID, status)
}
