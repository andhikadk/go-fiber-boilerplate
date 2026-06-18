package services

import (
	"errors"
	"time"

	"go-fiber-boilerplate/internal/dto"
	"go-fiber-boilerplate/internal/models"
	"go-fiber-boilerplate/pkg/utils"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidPassword = errors.New("invalid current password")
	ErrNoPasswordSet   = errors.New("this account does not have a password")
)

type UserService interface {
	GetUserByID(id uint) (*models.User, error)
	GetUserResponse(user *models.User) *dto.UserResponse
	UpdateProfile(userID uint, req *dto.UpdateProfileRequest) (*models.User, error)
	ChangePassword(userID uint, oldPassword, newPassword string) error
}

type userService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) UserService {
	return &userService{db: db}
}

func (s *userService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := s.db.Preload("Profile").First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (s *userService) GetUserResponse(user *models.User) *dto.UserResponse {
	role := roleString(user.Role)
	resp := &dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Role:      role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	if user.Profile != nil {
		resp.Profile = &dto.UserProfileResponse{
			FirstName: user.Profile.FirstName,
			LastName:  user.Profile.LastName,
		}
	}
	return resp
}

func (s *userService) UpdateProfile(userID uint, req *dto.UpdateProfileRequest) (*models.User, error) {
	var user models.User
	if err := s.db.Preload("Profile").First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	if user.Profile == nil {
		user.Profile = &models.UserProfile{UserID: user.ID}
	}
	user.Profile.FirstName = req.FirstName
	user.Profile.LastName = req.LastName
	user.Profile.UpdatedAt = time.Now()

	if err := s.db.Save(user.Profile).Error; err != nil {
		return nil, err
	}
	return s.GetUserByID(userID)
}

func (s *userService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	user, err := s.GetUserByID(userID)
	if err != nil {
		return err
	}
	if user.Password == nil {
		return ErrNoPasswordSet
	}
	if err := utils.VerifyPassword(oldPassword, *user.Password); err != nil {
		return ErrInvalidPassword
	}
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}
	return s.db.Model(user).Updates(map[string]interface{}{
		"password":   hashedPassword,
		"updated_at": time.Now(),
	}).Error
}
