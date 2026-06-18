package handlers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"go-fiber-boilerplate/internal/dto"
	"go-fiber-boilerplate/internal/services"
	"go-fiber-boilerplate/pkg/utils"
)

type Auth struct {
	authService services.AuthService
}

func NewAuth(authService services.AuthService) *Auth {
	return &Auth{authService: authService}
}

// Register godoc
//
//	@Summary		Register a new user
//	@Description	Register a new user with email, password, and profile name
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.RegisterRequest	true	"Registration request"
//	@Success		201		{object}	models.APIResponse	"User registered successfully"
//	@Failure		400		{object}	models.APIResponse	"Invalid request"
//	@Failure		409		{object}	models.APIResponse	"Email already registered"
//	@Router			/auth/register [post]
func (h *Auth) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		utils.LogCtx(c.UserContext(), "Auth").Error("Failed to parse request body", "error", err)
		return utils.BadRequestResponse(c, "Invalid request body")
	}
	if err := req.Validate(); err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}
	user, err := h.authService.Register(&req)
	if err != nil {
		if errors.Is(err, services.ErrEmailAlreadyRegistered) {
			return utils.ConflictResponse(c, "Email already registered")
		}
		utils.LogCtx(c.UserContext(), "Auth").Error("Registration failed", "error", err)
		return utils.InternalErrorResponse(c, "Failed to register user")
	}
	utils.LogCtx(c.UserContext(), "Auth").Info("User registered", "user_id", user.ID, "email", user.Email)
	return utils.CreatedResponse(c, "User registered successfully", user)
}

// Login godoc
//
//	@Summary		User login
//	@Description	Authenticate user and return access and refresh tokens
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.LoginRequest	true	"Login credentials"
//	@Success		200		{object}	models.APIResponse	"Login successful"
//	@Failure		400		{object}	models.APIResponse	"Invalid request"
//	@Failure		401		{object}	models.APIResponse	"Invalid credentials"
//	@Router			/auth/login [post]
func (h *Auth) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}
	if err := req.Validate(); err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}
	resp, err := h.authService.Login(&req)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) || errors.Is(err, services.ErrInactiveAccount) {
			return utils.UnauthorizedResponse(c, err.Error())
		}
		utils.LogCtx(c.UserContext(), "Auth").Error("Login failed", "email", req.Email, "error", err)
		return utils.InternalErrorResponse(c, "Failed to login")
	}
	return utils.SuccessResponse(c, fiber.StatusOK, "Login successful", resp)
}

// RefreshToken godoc
//
//	@Summary		Refresh access token
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.RefreshTokenRequest	true	"Refresh token"
//	@Success		200		{object}	models.APIResponse		"Token refreshed successfully"
//	@Failure		401		{object}	models.APIResponse		"Invalid refresh token"
//	@Router			/auth/refresh [post]
func (h *Auth) RefreshToken(c *fiber.Ctx) error {
	var req dto.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}
	if err := req.Validate(); err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}
	token, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		return utils.UnauthorizedResponse(c, err.Error())
	}
	return utils.SuccessResponse(c, fiber.StatusOK, "Token refreshed successfully", dto.RefreshTokenResponse{Token: token})
}

// ForgotPassword godoc
//
//	@Summary		Request password reset
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.ForgotPasswordRequest	true	"Email address"
//	@Success		200		{object}	models.APIResponse			"Reset email sent if email exists"
//	@Router			/auth/forgot-password [post]
func (h *Auth) ForgotPassword(c *fiber.Ctx) error {
	var req dto.ForgotPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}
	if err := req.Validate(); err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}
	if err := h.authService.ForgotPassword(req.Email); err != nil {
		if errors.Is(err, services.ErrPasswordResetDisabled) {
			return utils.SuccessResponse(c, fiber.StatusOK, "If the email exists, a reset link will be sent", nil)
		}
		utils.LogCtx(c.UserContext(), "Auth").Error("Forgot password failed", "error", err)
		return utils.InternalErrorResponse(c, "Failed to process request")
	}
	return utils.SuccessResponse(c, fiber.StatusOK, "If the email exists, a reset link will be sent", nil)
}

// ResetPassword godoc
//
//	@Summary		Reset password
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.ResetPasswordRequest	true	"Reset token and new password"
//	@Success		200		{object}	models.APIResponse			"Password reset successful"
//	@Failure		400		{object}	models.APIResponse			"Invalid token or request"
//	@Router			/auth/reset-password [post]
func (h *Auth) ResetPassword(c *fiber.Ctx) error {
	var req dto.ResetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}
	if err := req.Validate(); err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}
	if err := h.authService.ResetPassword(req.Token, req.NewPassword); err != nil {
		if errors.Is(err, services.ErrInvalidResetToken) {
			return utils.BadRequestResponse(c, "Invalid or expired reset token")
		}
		utils.LogCtx(c.UserContext(), "Auth").Error("Reset password failed", "error", err)
		return utils.InternalErrorResponse(c, "Failed to reset password")
	}
	return utils.SuccessResponse(c, fiber.StatusOK, "Password reset successful", nil)
}
