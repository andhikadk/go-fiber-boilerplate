package models

import (
	"time"

	"gorm.io/gorm"
)

type Resource struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"type:varchar(120);not null;index" json:"name"`
	Description *string        `gorm:"type:text" json:"description,omitempty"`
	Status      string         `gorm:"type:varchar(40);not null;default:'active';index" json:"status"`
	CreatedByID uint           `gorm:"not null;index" json:"created_by_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	CreatedBy *User `gorm:"foreignKey:CreatedByID" json:"created_by,omitempty"`
}

func (Resource) TableName() string {
	return "resources"
}
