package dto

import "time"

type CreateResourceRequest struct {
	Name        string  `json:"name" validate:"required,min=2,max=120" example:"Example Resource"`
	Description *string `json:"description" validate:"omitempty,max=2000" example:"A reusable sample resource"`
	Status      string  `json:"status" validate:"omitempty,oneof=active inactive archived" example:"active"`
}

func (r *CreateResourceRequest) Validate() error {
	return validate.Struct(r)
}

type UpdateResourceRequest struct {
	Name        *string `json:"name" validate:"omitempty,min=2,max=120" example:"Updated Resource"`
	Description *string `json:"description" validate:"omitempty,max=2000" example:"Updated description"`
	Status      *string `json:"status" validate:"omitempty,oneof=active inactive archived" example:"inactive"`
}

func (r *UpdateResourceRequest) Validate() error {
	return validate.Struct(r)
}

type ResourceResponse struct {
	ID          uint      `json:"id" example:"1"`
	Name        string    `json:"name" example:"Example Resource"`
	Description *string   `json:"description,omitempty" example:"A reusable sample resource"`
	Status      string    `json:"status" example:"active"`
	CreatedByID uint      `json:"created_by_id" example:"1"`
	CreatedAt   time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt   time.Time `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}
