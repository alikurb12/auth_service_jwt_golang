package email

import (
	"fmt"
	"net/smtp"

	"github.com/alikurb12/auth_service_jwt_golang/pkg/config"
)

type Sender struct {
	cfg config.EmailConfig
}

func NewSender(cfg config.EmailConfig) *Sender {
	return &Sender{cfg: cfg}
}

func (s *Sender) SendVerification(toEmail, token string) error {
	link := fmt.Sprintf("%s/api/auth/verify-email?token=%s", s.cfg.AppURL, token)
	subject := "Подтвердите ваш email"
	body := fmt.Sprintf(`
		Привет!

		Для подтверждения email перейдите по ссылке:
		%s

		Ссылка действительна 24 часа.

		Если вы не регистрировались — просто проигнорируйте это письмо.
		`,
		link,
	)

	return s.send(toEmail, subject, body)
}

func (s *Sender) send(toEmail, subject, body string) error {
	if s.cfg.SMTPUser == "" || s.cfg.SMTPPassword == "" {
		fmt.Printf("[EMAIL STUB] To: %s\nSubject: %s\n%s\n", toEmail, subject, body)
		return nil
	}
	auth := smtp.PlainAuth("", s.cfg.SMTPUser, s.cfg.SMTPPassword, s.cfg.SMTPHost)
	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		s.cfg.FromAddress, toEmail, subject, body,
	)

	addr := fmt.Sprintf("%s:%d", s.cfg.SMTPHost, s.cfg.SMTPPort)
	if err := smtp.SendMail(addr, auth, s.cfg.FromAddress, []string{toEmail}, []byte(msg)); err != nil {
		return fmt.Errorf("Send email to %s: %w error", toEmail, err)
	}

	return nil
}
