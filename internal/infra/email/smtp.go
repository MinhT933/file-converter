// SMTP email provider implementation viết tắt của Simple Mail Transfer Protocol là giao thức chuẩn để gửi email giữa các máy chủ email.
package email

import (
	"context"
	"fmt"
	"log"

	"gopkg.in/gomail.v2"
)

type SMTPProvider struct {
	host, user, pass, from string
	port                   int
}

func NewSMTPProvider(host string, port int, user, pass, from string) *SMTPProvider {
	return &SMTPProvider{
		host: host,
		user: user,
		port: port,
		pass: pass,
		from: from,
	}
}

func (p *SMTPProvider) SendConversionEmail(ctx context.Context, to, fileURL string) error {
	log.Printf("Đang gửi email đến: %s", to)
	// m = mail
	m := gomail.NewMessage()
	m.SetHeader("From", p.from)
	m.SetHeader("To", to)

	// viết email
	m.SetHeader("Subject", "File của bạn đã convert xong 🎉")

	body := fmt.Sprintf(
		"Chào bạn,\n\nFile của bạn đã được convert thành công!\nBạn có thể tải về ở đây:\n%s\n\nCảm ơn bạn đã sử dụng dịch vụ.",
		fileURL,
	)

	m.SetBody("text/plain", body)

	// tạo dialer để kết nối với smtp server
	dialer := gomail.NewDialer(p.host, p.port, p.user, p.pass)

	err := dialer.DialAndSend(m)
	if err != nil {
		log.Printf("❌ Lỗi gửi email: %v", err)
		return fmt.Errorf("failed to send conversion email: %w", err)
	}
	log.Println("✅ Gửi email conversion thành công!")
	return nil
}

func (p *SMTPProvider) SendWelcomeEmail(ctx context.Context, to string, userName string) error {
	log.Printf("Đang gửi email chào mừng đến: %s", to)
	m := gomail.NewMessage()
	m.SetHeader("From", p.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Chào mừng bạn đến với dịch vụ của chúng tôi!")

	body := "Chào bạn,\n\nCảm ơn bạn đã đăng ký sử dụng dịch vụ của chúng tôi!\nChúng tôi hy vọng bạn sẽ có trải nghiệm tuyệt vời.\n\nTrân trọng,\nĐội ngũ hỗ trợ."
	m.SetBody("text/plain", body)

	dialer := gomail.NewDialer(p.host, p.port, p.user, p.pass)

	err := dialer.DialAndSend(m)
	if err != nil {
		log.Printf("❌ Lỗi gửi email chào mừng: %v", err)
		return fmt.Errorf("failed to send welcome email: %w", err)
	}
	log.Println("✅ Gửi email chào mừng thành công!")
	return nil
}

func (p *SMTPProvider) SendPasswordResetEmail(ctx context.Context, to, resetLink string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", p.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Reset mật khẩu của bạn 🔐")

	body := fmt.Sprintf(
		"Chào bạn,\n\nBạn đã yêu cầu reset mật khẩu.\nVui lòng click vào link sau để đặt lại mật khẩu:\n%s\n\nNếu bạn không yêu cầu thay đổi này, vui lòng bỏ qua email này.",
		resetLink,
	)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(p.host, p.port, p.user, p.pass)
	return d.DialAndSend(m)
}
