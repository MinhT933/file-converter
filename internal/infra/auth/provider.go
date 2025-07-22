// provider.go
// Package auth provides authentication services for the application.
package auth

import (
	"context"
)

// Provider định nghĩa interface cho authentication services
type Provider interface {
	// VerifyToken xác thực JWT token và trả về thông tin user
	VerifyToken(ctx context.Context, token string) (*User, error)

	// GetUser lấy thông tin user theo UID
	GetUser(ctx context.Context, uid string) (*User, error)

	// CreateUser tạo user mới
	CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error)

	// UpdateUser cập nhật thông tin user
	UpdateUser(ctx context.Context, uid string, req *UpdateUserRequest) (*User, error)

	// SetCustomClaims đặt custom claims cho user
	SetCustomClaims(ctx context.Context, uid string, claims map[string]interface{}) error

	Login(ctx context.Context, email, password string) (*User, string, error)

	SocialLogin(ctx context.Context, provider, accessToken, idToken string) (*User, string, error)

	GoogleLogin(ctx context.Context, accessToken, idToken string) (*User, string, error)

	FacebookLogin(ctx context.Context, accessToken string) (*User, string, error)
}

// User đại diện cho một authenticated user
type User struct {
	UID           string                 `json:"uid"`
	Email         string                 `json:"email"`
	DisplayName   string                 `json:"display_name"`
	PhotoURL      string                 `json:"photo_url"`
	EmailVerified bool                   `json:"email_verified"`
	PhoneNumber   string                 `json:"phone_number"`
	Disabled      bool                   `json:"disabled"`
	CustomClaims  map[string]interface{} `json:"custom_claims"`
	CreatedAt     int64                  `json:"created_at"`
	LastLoginAt   int64                  `json:"last_login_at"`
}

// CreateUserRequest chứa dữ liệu để tạo user mới
type CreateUserRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=6"`
	DisplayName string `json:"display_name"`
	PhotoURL    string `json:"photo_url"`
	PhoneNumber string `json:"phone_number"`
}

// UpdateUserRequest chứa dữ liệu để cập nhật user
type UpdateUserRequest struct {
	Email       *string `json:"email,omitempty"`
	Password    *string `json:"password,omitempty"`
	DisplayName *string `json:"display_name,omitempty"`
	PhotoURL    *string `json:"photo_url,omitempty"`
	PhoneNumber *string `json:"phone_number,omitempty"`
	Disabled    *bool   `json:"disabled,omitempty"`
}

// AuthError đại diện cho các lỗi authentication
type AuthError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (e *AuthError) Error() string {
	if e.Details != "" {
		return e.Code + ": " + e.Message + " (" + e.Details + ")"
	}
	return e.Code + ": " + e.Message
}

// Common error codes
const (
	ErrCodeInvalidToken  = "INVALID_TOKEN"
	ErrCodeUserNotFound  = "USER_NOT_FOUND"
	ErrCodeEmailExists   = "EMAIL_EXISTS"
	ErrCodeWeakPassword  = "WEAK_PASSWORD"
	ErrCodeInvalidEmail  = "INVALID_EMAIL"
	ErrCodeUnauthorized  = "UNAUTHORIZED"
	ErrCodeInternalError = "INTERNAL_ERROR"
	ErrCodeTokenExpired  = "TOKEN_EXPIRED"
)

// NewAuthError tạo một AuthError mới
func NewAuthError(code, message string) *AuthError {
	return &AuthError{
		Code:    code,
		Message: message,
	}
}

// NewAuthErrorWithDetails tạo một AuthError với details
func NewAuthErrorWithDetails(code, message, details string) *AuthError {
	return &AuthError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// IsAuthError kiểm tra xem error có phải là AuthError không
func IsAuthError(err error) bool {
	_, ok := err.(*AuthError)
	return ok
}

// GetAuthErrorCode lấy error code từ error
func GetAuthErrorCode(err error) string {
	if authErr, ok := err.(*AuthError); ok {
		return authErr.Code
	}
	return ErrCodeInternalError
}
