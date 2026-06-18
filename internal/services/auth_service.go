package services

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"go-fiber-boilerplate/config"
	"go-fiber-boilerplate/internal/dto"
	"go-fiber-boilerplate/internal/models"
	"go-fiber-boilerplate/pkg/jwt"
	"go-fiber-boilerplate/pkg/utils"
	"gorm.io/gorm"
)

var (
	ErrEmailAlreadyRegistered = errors.New("email already registered")
	ErrInvalidCredentials     = errors.New("invalid email or password")
	ErrInactiveAccount        = errors.New("user account is inactive")
	ErrInvalidRefreshToken    = errors.New("invalid refresh token")
	ErrPasswordResetDisabled  = errors.New("password reset email service is not configured")
	ErrInvalidResetToken      = errors.New("invalid or expired reset token")
)

type AuthService interface {
	Register(req *dto.RegisterRequest) (*models.User, error)
	Login(req *dto.LoginRequest) (*dto.LoginResponse, error)
	RefreshToken(refreshTokenString string) (string, error)
	ForgotPassword(email string) error
	ResetPassword(token, newPassword string) error
}

type authService struct {
	db           *gorm.DB
	emailService EmailService
}

func NewAuthService(db *gorm.DB, emailService EmailService) AuthService {
	return &authService{db: db, emailService: emailService}
}

func (s *authService) Register(req *dto.RegisterRequest) (*models.User, error) {
	var existingUser models.User
	if err := s.db.Where("email = ?", strings.ToLower(req.Email)).First(&existingUser).Error; err == nil {
		return nil, ErrEmailAlreadyRegistered
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	role := "user"
	tx := s.db.Begin()
	user := &models.User{
		Email:               strings.ToLower(req.Email),
		Password:            &hashedPassword,
		PasswordIsSetByUser: true,
		Role:                &role,
		IsActive:            true,
	}
	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	profile := &models.UserProfile{
		UserID:    user.ID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}
	if err := tx.Create(profile).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	user.Profile = profile
	return user, nil
}

func (s *authService) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {
	var user models.User
	if err := s.db.Where("email = ?", strings.ToLower(req.Email)).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}
	if !user.IsActive {
		return nil, ErrInactiveAccount
	}
	if user.Password == nil {
		return nil, ErrInvalidCredentials
	}
	if err := utils.VerifyPassword(req.Password, *user.Password); err != nil {
		return nil, ErrInvalidCredentials
	}

	tm := jwt.NewTokenManager(config.AppConfig.JWTSecret)
	role := roleString(user.Role)
	accessToken, err := tm.GenerateAccessToken(user.ID, user.Email, role, config.AppConfig.JWTExpiry)
	if err != nil {
		return nil, err
	}
	refreshToken, err := tm.GenerateRefreshToken(user.ID, user.Email, config.AppConfig.JWTRefreshExpiry)
	if err != nil {
		return nil, err
	}
	return &dto.LoginResponse{
		Token:        accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(config.AppConfig.JWTExpiry.Seconds()),
		Role:         role,
	}, nil
}

func (s *authService) RefreshToken(refreshTokenString string) (string, error) {
	tm := jwt.NewTokenManager(config.AppConfig.JWTSecret)
	claims, err := tm.ValidateRefreshToken(refreshTokenString)
	if err != nil {
		return "", ErrInvalidRefreshToken
	}

	var user models.User
	if err := s.db.First(&user, claims.UserID).Error; err != nil {
		return "", err
	}
	if !user.IsActive {
		return "", ErrInactiveAccount
	}
	return tm.GenerateAccessToken(user.ID, user.Email, roleString(user.Role), config.AppConfig.JWTExpiry)
}

func (s *authService) ForgotPassword(email string) error {
	if s.emailService == nil || !s.emailService.Enabled() {
		return ErrPasswordResetDisabled
	}

	var user models.User
	if err := s.db.Where("email = ?", strings.ToLower(email)).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	token := utils.RandomString(32)
	reset := &models.PasswordReset{
		UserID:    user.ID,
		TokenHash: hashToken(token),
		ExpiresAt: time.Now().Add(30 * time.Minute),
	}
	if err := s.db.Create(reset).Error; err != nil {
		return err
	}
	return s.emailService.SendPasswordReset(user.Email, token)
}

func (s *authService) ResetPassword(token, newPassword string) error {
	var reset models.PasswordReset
	if err := s.db.Where("token_hash = ? AND used_at IS NULL AND expires_at > ?", hashToken(token), time.Now()).First(&reset).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrInvalidResetToken
		}
		return err
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	now := time.Now()
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.User{}).Where("id = ?", reset.UserID).Updates(map[string]interface{}{
			"password":                hashedPassword,
			"password_is_set_by_user": true,
			"updated_at":              now,
		}).Error; err != nil {
			return err
		}
		return tx.Model(&reset).Update("used_at", now).Error
	})
}

func roleString(role *string) string {
	if role == nil {
		return ""
	}
	return *role
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
