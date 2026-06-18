package mailer

import (
	"fmt"
	"io"
	"time"

	"go-fiber-boilerplate/pkg/utils"
	"gopkg.in/gomail.v2"
)

type SMTPClient struct {
	host      string
	port      int
	username  string
	password  string
	fromName  string
	fromEmail string
}

func NewSMTPClient(host string, port int, username, password, fromName, fromEmail string) Mailer {
	return &SMTPClient{
		host:      host,
		port:      port,
		username:  username,
		password:  password,
		fromName:  fromName,
		fromEmail: fromEmail,
	}
}

func (s *SMTPClient) SendEmail(msg *EmailMessage) error {
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(s.fromEmail, s.fromName))
	m.SetHeader("To", msg.To...)
	m.SetHeader("Subject", msg.Subject)

	if msg.HTMLBody != "" {
		m.SetBody("text/html", msg.HTMLBody)
	}
	if msg.TextBody != "" {
		m.AddAlternative("text/plain", msg.TextBody)
	}

	for _, att := range msg.Attachments {
		filename := att.Filename
		content := att.Content
		m.Attach(filename, gomail.SetCopyFunc(func(w io.Writer) error {
			_, err := w.Write(content)
			return err
		}))
	}

	s.logRequest(msg)
	start := time.Now()
	dialer := gomail.NewDialer(s.host, s.port, s.username, s.password)
	if err := dialer.DialAndSend(m); err != nil {
		s.logError(msg, err, time.Since(start))
		return fmt.Errorf("failed to send email: %w", err)
	}

	s.logResponse(msg, time.Since(start))
	return nil
}

func (s *SMTPClient) SendEmailAsync(msg *EmailMessage) error {
	return s.SendEmail(msg)
}

func (s *SMTPClient) logRequest(msg *EmailMessage) {
	utils.Log("SMTP").Info("Request",
		"host", s.host,
		"port", s.port,
		"from", s.fromEmail,
		"to", msg.To,
		"subject", msg.Subject,
		"html_size", len(msg.HTMLBody),
		"text_size", len(msg.TextBody),
		"attachments", len(msg.Attachments),
	)
}

func (s *SMTPClient) logResponse(msg *EmailMessage, elapsed time.Duration) {
	utils.Log("SMTP").Info("Response",
		"to", msg.To,
		"subject", msg.Subject,
		"status", "sent",
		"duration_ms", elapsed.Milliseconds(),
	)
}

func (s *SMTPClient) logError(msg *EmailMessage, err error, elapsed time.Duration) {
	utils.Log("SMTP").Error("Request failed",
		"host", s.host,
		"port", s.port,
		"to", msg.To,
		"subject", msg.Subject,
		"duration_ms", elapsed.Milliseconds(),
		"error", err,
	)
}
