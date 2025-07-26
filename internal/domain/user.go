package domain

import (
  "github.com/google/uuid"
  "time"
)
type User struct {
	UserID  string  `json:"user_id" db:"user_id"`
	Email   string  `json:"email" db:"email" validate:"required,email"`
	Name    string  `json:"name" db:"name" validate:"required"`
	AvatarURL *string  `json:"avatar_url" db:"avatar_url"`
    Provider *string  `json:"provider" db:"provider"`
	EmailVerified bool   `json:"email_verified" db:"email_verified"`
	CreatedAt    time.Time `db:"created_at"`
	IsPremium     bool   `json:"is_premium" db:"is_premium"`
	UpdatedAt     time.Time  `db:"updated_at"`
	ProviderID    *string  `json:"provider_id" db:"provider_id"`

	Roles   []Role `json:"roles" db:"roles"`
	Conversions []Conversion `json:"conversions" db:"conversions"`
	FileLogs []FileLog `json:"file_logs" db:"file_logs"`
}

type Role  struct {
	RoleID   string `json:"role_id" db:"role_id"`
	RoleName string `json:"role_name" db:"role_name"`
}

type UserRole struct {
    UserID int `json:"user_id" db:"user_id"`
    RoleID int `json:"role_id" db:"role_id"`
}

func (u *User) GenerateID(){
	u.UserID = uuid.New().String()
}
