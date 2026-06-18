package dto

import "github.com/go-playground/validator/v10"

var validate = validator.New()

type RegisterRequest struct {
	Email     string  `json:"email" validate:"required,email" example:"john@example.com"`
	Password  string  `json:"password" validate:"required,min=8,max=255" example:"password123"`
	FirstName string  `json:"first_name" validate:"required,min=2,max=120" example:"John"`
	LastName  *string `json:"last_name" validate:"omitempty,max=120" example:"Doe"`
}

func (r *RegisterRequest) Validate() error {
	return validate.Struct(r)
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email" example:"john@example.com"`
	Password string `json:"password" validate:"required,min=1" example:"password123"`
}

func (r *LoginRequest) Validate() error {
	return validate.Struct(r)
}

type LoginResponse struct {
	Token        string `json:"token" example:"eyJhbGciOi..."`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOi..."`
	ExpiresIn    int64  `json:"expires_in" example:"900"`
	Role         string `json:"role" example:"user"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required" example:"eyJhbGciOi..."`
}

func (r *RefreshTokenRequest) Validate() error {
	return validate.Struct(r)
}

type RefreshTokenResponse struct {
	Token string `json:"token" example:"eyJhbGciOi..."`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email" example:"john@example.com"`
}

func (r *ForgotPasswordRequest) Validate() error {
	return validate.Struct(r)
}

type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required" example:"abc123"`
	NewPassword string `json:"new_password" validate:"required,min=8,max=255" example:"newpassword123"`
}

func (r *ResetPasswordRequest) Validate() error {
	return validate.Struct(r)
}
