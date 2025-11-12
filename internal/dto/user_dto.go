package dto

import (
	"errors"
	"strings"
	"time"
)

// UpdateProfileRequest is the request body for updating user profile
type UpdateProfileRequest struct {
	Name string `json:"name" example:"John Doe Updated"`
}

// Validate validates the UpdateProfileRequest
func (r *UpdateProfileRequest) Validate() error {
	if strings.TrimSpace(r.Name) == "" {
		return errors.New("name is required and cannot be empty")
	}
	if len(r.Name) < 2 {
		return errors.New("name must be at least 2 characters")
	}
	if len(r.Name) > 255 {
		return errors.New("name cannot exceed 255 characters")
	}
	return nil
}

// ChangePasswordRequest is the request body for changing password
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" example:"oldpassword123"`
	NewPassword string `json:"new_password" example:"newpassword123"`
}

// Validate validates the ChangePasswordRequest
func (r *ChangePasswordRequest) Validate() error {
	// Validate Old Password
	if r.OldPassword == "" {
		return errors.New("old_password is required and cannot be empty")
	}

	// Validate New Password
	if r.NewPassword == "" {
		return errors.New("new_password is required and cannot be empty")
	}
	if len(r.NewPassword) < 6 {
		return errors.New("new_password must be at least 6 characters")
	}
	if len(r.NewPassword) > 255 {
		return errors.New("new_password cannot exceed 255 characters")
	}

	// Check if passwords are different
	if r.OldPassword == r.NewPassword {
		return errors.New("new password must be different from old password")
	}

	return nil
}

// UserResponse is the response for user data (public information only)
type UserResponse struct {
	ID        uint      `json:"id" example:"1"`
	Name      string    `json:"name" example:"John Doe"`
	Email     string    `json:"email" example:"john@example.com"`
	Role      string    `json:"role" example:"user"`
	IsActive  bool      `json:"is_active" example:"true"`
	CreatedAt time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}
