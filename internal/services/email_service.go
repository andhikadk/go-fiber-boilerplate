package services

import "go-fiber-boilerplate/pkg/utils"

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

func (noopEmailService) SendPasswordReset(email, token string) error {
	utils.Log("Email").Warn("Password reset email skipped because email service is disabled", "email", email)
	return nil
}
