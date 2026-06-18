package dto

import "time"

type UpdateProfileRequest struct {
	FirstName string  `json:"first_name" validate:"required,min=2,max=120" example:"John"`
	LastName  *string `json:"last_name" validate:"omitempty,max=120" example:"Doe"`
}

func (r *UpdateProfileRequest) Validate() error {
	return validate.Struct(r)
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required,min=1" example:"oldpassword123"`
	NewPassword string `json:"new_password" validate:"required,min=8,max=255,nefield=OldPassword" example:"newpassword123"`
}

func (r *ChangePasswordRequest) Validate() error {
	return validate.Struct(r)
}

type UserProfileResponse struct {
	FirstName string  `json:"first_name" example:"John"`
	LastName  *string `json:"last_name,omitempty" example:"Doe"`
}

type UserResponse struct {
	ID        uint                 `json:"id" example:"1"`
	Email     string               `json:"email" example:"john@example.com"`
	Role      string               `json:"role" example:"user"`
	IsActive  bool                 `json:"is_active" example:"true"`
	CreatedAt time.Time            `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt time.Time            `json:"updated_at" example:"2024-01-01T00:00:00Z"`
	Profile   *UserProfileResponse `json:"profile,omitempty"`
}
