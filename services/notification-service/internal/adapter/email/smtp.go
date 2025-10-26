package email

import (
	"context"
	"fmt"
	"net/smtp"

	"github.com/squ1ky/avigo-c2c-marketplace/services/notification-service/internal/config"
	"github.com/squ1ky/avigo-c2c-marketplace/services/notification-service/internal/domain"
)

type SMTPSender struct {
	host     string
	port     int
	username string
	password string
	from     string
	fromName string
}

func NewSMTPSender(cfg *config.EmailConfig) *SMTPSender {
	return &SMTPSender{
		host:     cfg.SMTPHost,
		port:     cfg.SMTPPort,
		username: cfg.SMTPUsername,
		password: cfg.SMTPPassword,
		from:     cfg.FromAddress,
		fromName: cfg.FromName,
	}
}

func (s *SMTPSender) Send(ctx context.Context, msg *domain.EmailMessage) error {
	if msg == nil {
		return fmt.Errorf("email message cannot be nil")
	}

	if msg.To == "" {
		return fmt.Errorf("recipient email cannot be empty")
	}

	select {
	case <-ctx.Done():
		return fmt.Errorf("email send cancelled: %w", ctx.Err())
	default:
	}

	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	auth := smtp.PlainAuth("", s.username, s.password, s.host)

	mimeMessage := s.buildMIMEMessage(msg)

	err := smtp.SendMail(addr, auth, s.from, []string{msg.To}, []byte(mimeMessage))
	if err != nil {
		return fmt.Errorf("failed to send email via SMTP to %s: %w", msg.To, err)
	}

	return nil
}

func (s *SMTPSender) buildMIMEMessage(msg *domain.EmailMessage) string {
	from := fmt.Sprintf("%s <%s>", s.fromName, s.from)

	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = msg.To
	headers["Subject"] = msg.Subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + msg.HTMLBody

	return message
}
