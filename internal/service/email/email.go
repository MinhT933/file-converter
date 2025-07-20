package email

import (
	"fmt"
    "gopkg.in/gomail.v2"
    "log"
)

type Service struct {
	host, user, pass, from string 
	port                   int
}

func NewService(host string, port int, user, pass, from string) * Service {
    return &Service{host: host, port: port, user: user, pass: pass, from: from}
} 

func (s *Service) SendConversionEmail(to, fileURL string) error {
    log.Printf("email đã được gửi đi")
    m := gomail.NewMessage()
    m.SetHeader("From", s.from)
    m.SetHeader("To", to)
    m.SetHeader("Subject", "File của bạn đã convert xong 🎉")
    body := fmt.Sprintf(
        "Chào bạn,\n\nFile của bạn đã được convert thành công!\nBạn có thể tải về ở đây:\n%s\n\nCảm ơn bạn đã sử dụng dịch vụ.",
        fileURL,
    )
    m.SetBody("text/plain", body)
  
    d := gomail.NewDialer(s.host, s.port, s.user, s.pass)

      err := d.DialAndSend(m)
    if err != nil {
        log.Printf("❌ Lỗi gửi email: %v", err)
    } else {
        log.Println("✅ Gửi email thành công!")
    }
    return err
}