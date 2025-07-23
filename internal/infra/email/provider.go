package email

type Provider interface {
	SendConversionEmail(to, fileURL string) error
	SendWelcomeEmail(to, name string) error
	SendPasswordResetEmail(to, resetLink string) error
}

type EmailRequest struct {
	To      string
	Subject string
	Body    string
	IsHTML  bool
}
