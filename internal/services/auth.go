package services

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/MinhT933/file-converter/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthService struct {
	userRepo domain.UserRepository
}

// hàm dựng factory func cho AuthService
func NewAuthService(userRepo domain.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func (s *AuthService) HandleProviderLogin(ctx context.Context, providerData domain.ProviderData) (*domain.AuthResult, error) {
	// Xử lý đăng nhập với data nhà cung cấp
	provider := extractProviderName(providerData.ProviderID)

	existingUser, err := s.userRepo.FindByEmail(ctx, providerData.Email)
	if err != nil {
		log.Print("User does not exist, creating new user")
		return s.createOAuthUser(ctx, providerData, provider)
	}

	if existingUser.Provider != nil && *existingUser.Provider == provider {
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

	if err != nil {
		log.Printf("Error creating session: %v", err)
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

var jwtSecret = []byte(os.Getenv("JWT_SECRET_KEY"))

func generateJWT(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(1 * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func extractProviderName(providerID string) string {
	if idx := strings.Index(providerID, "."); idx > 0 {
		return providerID[:idx]
	}
	return providerID
}
