package mailer

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/gogaruda/apperror"
	"github.com/gogaruda/auth/internal/config"
	"gopkg.in/gomail.v2"
)

type Mailer interface {
	Send(to, subject, body string) error
	GenerateRandomToken(n int) (string, error)
}

type mailer struct {
	cfg config.EmailConfig
}

func NewMailer(c config.EmailConfig) Mailer {
	return &mailer{cfg: c}
}

func (mr *mailer) Send(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", mr.cfg.MailFromAddress)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(
		mr.cfg.MailHost,
		mr.cfg.MailPort,
		mr.cfg.MailUsername,
		mr.cfg.MailPassword,
	)

	return d.DialAndSend(m)
}

func (mr *mailer) GenerateRandomToken(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", apperror.New(apperror.CodeInternalError, "gagal membuat token", err)
	}

	return hex.EncodeToString(bytes), nil
}
