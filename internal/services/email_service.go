package services

import (
	"fmt"
	"html"
	"strings"

	"go-fiber-boilerplate/pkg/mailer"
	"go-fiber-boilerplate/pkg/utils"
)

type EmailService interface {
	Enabled() bool
	SendPasswordReset(email, token string) error
}

type noopEmailService struct{}

func NewNoopEmailService() EmailService {
	return noopEmailService{}
}

func (noopEmailService) Enabled() bool {
	return false
}

func (noopEmailService) SendPasswordReset(email, _ string) error {
	utils.Log("Email").Warn("Password reset email skipped because email service is disabled", "email", email)
	return nil
}

type smtpEmailService struct {
	mailer           mailer.Mailer
	appName          string
	passwordResetURL string
}

func NewEmailService(m mailer.Mailer, appName, passwordResetURL string) EmailService {
	return &smtpEmailService{
		mailer:           m,
		appName:          appName,
		passwordResetURL: passwordResetURL,
	}
}

func (s *smtpEmailService) Enabled() bool {
	return s != nil && s.mailer != nil
}

func (s *smtpEmailService) SendPasswordReset(email, token string) error {
	resetURL := s.buildPasswordResetURL(token)
	msg := &mailer.EmailMessage{
		To:       []string{email},
		Subject:  fmt.Sprintf("Reset your %s password", s.appName),
		HTMLBody: s.passwordResetHTML(resetURL),
		TextBody: s.passwordResetText(resetURL),
	}
	return s.mailer.SendEmail(msg)
}

func (s *smtpEmailService) buildPasswordResetURL(token string) string {
	if strings.Contains(s.passwordResetURL, "{token}") {
		return strings.ReplaceAll(s.passwordResetURL, "{token}", token)
	}
	separator := "?"
	if strings.Contains(s.passwordResetURL, "?") {
		separator = "&"
	}
	return s.passwordResetURL + separator + "token=" + token
}

func (s *smtpEmailService) passwordResetHTML(resetURL string) string {
	appName := html.EscapeString(s.appName)
	escapedURL := html.EscapeString(resetURL)
	return fmt.Sprintf(`<!doctype html>
<html>
<body style="font-family: Arial, sans-serif; color: #111827; line-height: 1.5;">
  <h2>Reset your %s password</h2>
  <p>We received a request to reset your password.</p>
  <p>
    <a href="%s" style="display: inline-block; padding: 10px 16px; background: #111827; color: #ffffff; text-decoration: none; border-radius: 6px;">
      Reset password
    </a>
  </p>
  <p>If the button does not work, copy and paste this link into your browser:</p>
  <p><a href="%s">%s</a></p>
  <p>This link will expire soon. If you did not request a password reset, you can ignore this email.</p>
</body>
</html>`, appName, escapedURL, escapedURL, escapedURL)
}

func (s *smtpEmailService) passwordResetText(resetURL string) string {
	return fmt.Sprintf(`Reset your %s password

We received a request to reset your password.

Open this link to reset your password:
%s

This link will expire soon. If you did not request a password reset, you can ignore this email.
`, s.appName, resetURL)
}
