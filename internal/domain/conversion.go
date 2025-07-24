package domain

import "time"

type ConversionPreset struct {
    PresetID   string       `json:"preset_id" db:"preset_id"`
    UserID     string       `json:"user_id" db:"user_id"`
    Name       string    `json:"name" db:"name"`
    Parameters string    `json:"parameters" db:"parameters"`
    CreatedAt  time.Time `json:"created_at" db:"created_at"`
    
    User *User `json:"user,omitempty" db:"-"`
}

type Conversion struct {
    ConversionID     string       `json:"conversion_id" db:"conversion_id"`
    UserID           string       `json:"user_id" db:"user_id"`
    OriginalFilename string    `json:"original_filename" db:"original_filename"`
    ConvertedFilename string   `json:"converted_filename" db:"converted_filename"`
    CreatedAt        time.Time `json:"created_at" db:"created_at"`
    ExpiresAt        time.Time `json:"expires_at" db:"expires_at"`
    
    User     *User       `json:"user,omitempty" db:"-"`
    FileLogs []FileLog   `json:"file_logs,omitempty" db:"-"`
}