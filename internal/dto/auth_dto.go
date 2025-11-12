package dto

import (
	"errors"
	"regexp"
	"strings"
)

// RegisterRequest is the request body for user registration
type RegisterRequest struct {
	Name     string `json:"name" example:"John Doe"`
	Email    string `json:"email" example:"john@example.com"`
	Password string `json:"password" example:"password123"`
}

// Validate validates the RegisterRequest
func (r *RegisterRequest) Validate() error {
	// Validate Name
	if strings.TrimSpace(r.Name) == "" {
		return errors.New("name is required and cannot be empty")
	}
	if len(r.Name) < 2 {
		return errors.New("name must be at least 2 characters")
	}
	if len(r.Name) > 255 {
		return errors.New("name cannot exceed 255 characters")
	}

	// Validate Email
	if strings.TrimSpace(r.Email) == "" {
		return errors.New("email is required and cannot be empty")
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(r.Email) {
		return errors.New("invalid email format")
	}

	// Validate Password
	if r.Password == "" {
		return errors.New("password is required and cannot be empty")
	}
	if len(r.Password) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	if len(r.Password) > 255 {
		return errors.New("password cannot exceed 255 characters")
	}

	return nil
}

// LoginRequest is the request body for user login
type LoginRequest struct {
	Email    string `json:"email" example:"john@example.com"`
	Password string `json:"password" example:"password123"`
}

// Validate validates the LoginRequest
func (r *LoginRequest) Validate() error {
	// Validate Email
	if strings.TrimSpace(r.Email) == "" {
		return errors.New("email is required and cannot be empty")
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(r.Email) {
		return errors.New("invalid email format")
	}

	// Validate Password
	if r.Password == "" {
		return errors.New("password is required and cannot be empty")
	}

	return nil
}

// LoginResponse is the response for successful login
type LoginResponse struct {
	Token        string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	ExpiresIn    int64  `json:"expires_in" example:"900"`
}

// RefreshTokenRequest is the request body for refreshing access token
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// Validate validates the RefreshTokenRequest
func (r *RefreshTokenRequest) Validate() error {
	if strings.TrimSpace(r.RefreshToken) == "" {
		return errors.New("refresh_token is required and cannot be empty")
	}
	return nil
}

// RefreshTokenResponse is the response for successful token refresh
type RefreshTokenResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}
