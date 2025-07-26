package domain

type ProviderData struct {
    UID           string `json:"uid" validate:"required"`
    Email         string `json:"email" validate:"required,email"`
    DisplayName   string `json:"display_name" validate:"required"`
    PhotoURL      string `json:"photo_url"`
    ProviderID    string `json:"provider_id" validate:"required"` // "google.com"
    EmailVerified bool   `json:"email_verified"`
}

type AuthResult struct {
    SessionToken     string        `json:"session_token,omitempty"`
    User            *User          `json:"user,omitempty"`
    RequiresLinking bool           `json:"requires_linking,omitempty"`
    ExistingUser    *User          `json:"existing_user,omitempty"`
    ProviderData    *ProviderData  `json:"provider_data,omitempty"`
}