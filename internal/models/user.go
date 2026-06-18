package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID                  uint           `gorm:"primaryKey" json:"id"`
	Email               string         `gorm:"uniqueIndex;size:255;not null" json:"email"`
	Password            *string        `gorm:"size:255" json:"-"`
	PasswordIsSetByUser bool           `gorm:"not null;default:false" json:"password_is_set_by_user"`
	Role                *string        `gorm:"size:50;index" json:"role"`
	IsActive            bool           `gorm:"not null;default:true" json:"is_active"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"-"`

	Profile *UserProfile `gorm:"foreignKey:UserID" json:"profile,omitempty"`
}

func (User) TableName() string {
	return "users"
}

type UserProfile struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"uniqueIndex;not null" json:"user_id"`
	FirstName string    `gorm:"type:varchar(120);not null" json:"first_name"`
	LastName  *string   `gorm:"type:varchar(120)" json:"last_name,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (UserProfile) TableName() string {
	return "user_profiles"
}
