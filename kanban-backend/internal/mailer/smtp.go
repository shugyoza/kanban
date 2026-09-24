package mailer

import (
	"context"
	"fmt"
	"net/smtp"
	"net/url"
	"os"
	"strings"
	"time"
)

type SMTPMailer struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

func NewSMTPMailerFromEnv() *SMTPMailer {
	return &SMTPMailer{
		Host:     envOrDefault("SMTP_HOST", "localhost"),
		Port:     envOrDefault("SMTP_PORT", "1025"),
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
		From:     envOrDefault("SMTP_FROM", "no-reply@kanban.local"),
	}
}

func (m *SMTPMailer) SendInvitation(ctx context.Context, recipient string, token string, registerURL string, expiresAt time.Time) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	inviteURL := registerURL + "?email=" + recipient + "&token=" + url.QueryEscape(token)
	body := fmt.Sprintf("You have been invited to join Kanban.\n\nRegister your account here:\n%s\n\nThis invitation expires at %s UTC.\n", inviteURL, expiresAt.UTC().Format(time.RFC3339))
	message := strings.Join([]string{
		"From: " + m.From,
		"To: " + recipient,
		"Subject: Your Kanban invitation",
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
		"",
		body,
	}, "\r\n")

	var auth smtp.Auth
	if m.Username != "" {
		auth = smtp.PlainAuth("", m.Username, m.Password, m.Host)
	}

	return smtp.SendMail(m.Host+":"+m.Port, auth, m.From, []string{recipient}, []byte(message))
}

func envOrDefault(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}