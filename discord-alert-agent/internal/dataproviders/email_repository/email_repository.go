package email_repository

import (
	"fmt"
	"net/smtp"

	"github.com/rs/zerolog"
)

// EmailRepository описывает интерфейс для отправки email‑алертов.
type EmailRepository interface {
	SendEmail(to, subject, body string) error
}

type emailRepository struct {
	smtpHost string
	smtpPort int
	username string
	password string
	from     string
	logger   *zerolog.Logger
}

// NewEmailRepository создаёт экземпляр emailRepository и принимает логгер.
func NewEmailRepository(logger *zerolog.Logger, smtpHost string, smtpPort int, username, password, from string) (EmailRepository, error) {
	return &emailRepository{
		smtpHost: smtpHost,
		smtpPort: smtpPort,
		username: username,
		password: password,
		from:     from,
		logger:   logger,
	}, nil
}

// SendEmail отправляет письмо по заданным параметрам через SMTP и логирует результат.
func (r *emailRepository) SendEmail(to, subject, body string) error {
	addr := fmt.Sprintf("%s:%d", r.smtpHost, r.smtpPort)
	auth := smtp.PlainAuth("", r.username, r.password, r.smtpHost)
	// Исправленный формат письма с заголовком From
	msg := []byte("From: " + r.from + "\r\n" + // Добавляем From
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")
	err := smtp.SendMail(addr, auth, r.from, []string{to}, msg)
	if err != nil {
		r.logger.Error().Err(err).Msgf("Failed to send email to %s", to)
		return err
	}
	r.logger.Info().Msgf("Successfully sent email to %s", to)
	return nil
}
