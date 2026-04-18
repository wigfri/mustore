package mail

import (
	"fmt"

	"github.com/wigfri/mustore/services/config"
	"gopkg.in/gomail.v2"
)

type Sender struct {
	cfg *config.Config
}

func NewSender(cfg *config.Config) *Sender {
	return &Sender{cfg: cfg}
}

func (s *Sender) SendVerification(to, code string) error {
	subject := "Код подтверждения регистрации"
	body := fmt.Sprintf("Ваш код подтверждения: %s\n\nКод действителен 24 часа.", code)
	return s.send(to, subject, body)
}

func (s *Sender) SendLoginCode(to, code string) error {
	subject := "Код для входа"
	body := fmt.Sprintf("Ваш код для входа: %s\n\nКод действителен 10 минут.", code)
	return s.send(to, subject, body)
}

func (s *Sender) send(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.cfg.MailFrom())
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain; charset=UTF-8", body)

	d := gomail.NewDialer(s.cfg.SMTPHost(), s.cfg.SMTPPort(), s.cfg.SMTPUser(), s.cfg.SMTPPassword())
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("smtp send: %w", err)
	}
	return nil
}
