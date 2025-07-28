package services

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"errors"
	"strconv"

	"github.com/MinhT933/file-converter/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthService struct {
	userRepo domain.UserRepository
}

// hàm dựng factory func cho AuthService
// đảm bảo nhất quán
// trong việc khởi tạo AuthService với UserRepository
// và có thể mở rộng trong tương lai nếu cần
// Ví dụ: nếu cần thêm các dịch vụ khác như EmailService, có thể thêm vào

func NewAuthService(userRepo domain.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func (s *AuthService) HandleProviderLogin(ctx context.Context, providerData domain.ProviderData) (*domain.AuthResult, error) {
	var ErrUserNotFound = errors.New("user not found")
	// Xử lý đăng nhập với data nhà cung cấp
	provider := extractProviderName(providerData.ProviderID)

	existingUser, err := s.userRepo.FindByEmail(ctx, providerData.Email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			log.Print("User does not exist, creating new user")
			return s.createOAuthUser(ctx, providerData, provider)
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if existingUser.Provider != nil && strings.EqualFold(*existingUser.Provider, provider) {
		//user tồn tại
		return s.LoginExistingUser(ctx, existingUser)
	}

	return &domain.AuthResult{
		RequiresLinking: true,
		ExistingUser:    existingUser,
		ProviderData:    &providerData,
	}, nil
}

func (s *AuthService) LoginExistingUser(ctx context.Context, user *domain.User) (*domain.AuthResult, error) {
	// Sửa: Pass user thay vì user.UserID
	sessionToken, err := generateJWT(user.UserID)
	if err != nil {
		log.Printf("Error creating session for existing user: %v", err)
		return nil, fmt.Errorf("failed to create session: %w", err)
	}
	return &domain.AuthResult{
		SessionToken: sessionToken,
		User:         user,
	}, nil
}

func (s *AuthService) createOAuthUser(ctx context.Context, providerData domain.ProviderData, provider string) (*domain.AuthResult, error) {
	// Tạo người dùng mới với thông tin từ nhà cung cấp
	newUser := domain.User{
		UserID:        uuid.New().String(),
		Email:         providerData.Email,
		Name:          providerData.DisplayName,
		Provider:      &provider,              // Pointer
		ProviderID:    &providerData.UID,      // Pointer, sử dụng UID
		AvatarURL:     &providerData.PhotoURL, // Pointer, sử dụng PhotoURL
		EmailVerified: providerData.EmailVerified,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	log.Printf("Creating new user: %+v", newUser)

	err := s.userRepo.Create(ctx, newUser)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return nil, err
	}

	token, err := generateJWT(newUser.UserID)
	if err != nil {
		log.Printf("Error generating JWT for new user: %v", err)
		return nil, err
	}

	return &domain.AuthResult{
		SessionToken: token,
		User:         &newUser,
	}, nil
}

var errMissingJWTSecret = errors.New("JWT_SECRET_KEY environment variable is not set")

func generateJWT(userID string) (string, error) {
	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		return "", errMissingJWTSecret
	}

	expMinutes := 1
	if v := os.Getenv("JWT_EXP_MINUTES"); v != "" {
		if m, err := strconv.Atoi(v); err == nil {
			expMinutes = m
		}
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Duration(expMinutes) * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func extractProviderName(providerID string) string {
	if idx := strings.Index(providerID, "."); idx > 0 {
		return providerID[:idx]
	}
	return providerID
}
