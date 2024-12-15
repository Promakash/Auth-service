package infra

import (
	"fmt"
	"net/smtp"
)

type SMTPConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type SMTPService interface {
	SendMail(email []string, subject, body string) error
}

type MailService struct {
	SMTPConfig
	auth smtp.Auth
}

func NewEmailSender(cfg SMTPConfig) SMTPService {
	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
	return &MailService{
		SMTPConfig: cfg,
		auth:       auth,
	}
}

func (s *MailService) SendMail(email []string, subject, body string) error {
	msg := []byte(fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		s.Username, email, subject, body,
	))

	//auth disabled for docker tests
	err := smtp.SendMail(
		fmt.Sprintf("%s:%s", s.Host, s.Port),
		nil,
		s.Username,
		email,
		msg,
	)

	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}
