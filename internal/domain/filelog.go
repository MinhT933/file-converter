package domain

import "time"

type FileLog struct {
    LogID        string       `json:"log_id" db:"log_id"`
    ConversionID int       `json:"conversion_id" db:"conversion_id"`
    UserID       int       `json:"user_id" db:"user_id"`
    Action       string    `json:"action" db:"action"`
    LoggedAt     time.Time `json:"logged_at" db:"logged_at"`
    
    Conversion *Conversion `json:"conversion,omitempty" db:"-"`
    User       *User       `json:"user,omitempty" db:"-"`
}