package handlers

import (
	"go-fiber-boilerplate/internal/database"
	"go-fiber-boilerplate/internal/dto"
	"go-fiber-boilerplate/internal/middleware"
	"go-fiber-boilerplate/internal/services"
	"go-fiber-boilerplate/internal/utils"
	pkgUtils "go-fiber-boilerplate/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// Register godoc
//
//	@Summary		Register a new user
//	@Description	Register a new user with name, email, and password
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.RegisterRequest	true	"Registration request"
//	@Success		201		{object}	models.APIResponse	"User registered successfully"
//	@Failure		400		{object}	models.APIResponse	"Invalid request or validation error"
//	@Failure		409		{object}	models.APIResponse	"Email already registered"
//	@Router			/auth/register [post]
func Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		utils.ErrorLogger.Printf("[Register] Failed to parse request body: %v", err)
		return pkgUtils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate request using DTO's self-validation
	if err := req.Validate(); err != nil {
		utils.ErrorLogger.Printf("[Register] Validation failed: %v", err)
		return pkgUtils.BadRequestResponse(c, err.Error())
	}

	// Register user
	authService := services.NewAuthService(database.GetDB())
	user, err := authService.Register(&req)
	if err != nil {
		utils.ErrorLogger.Printf("[Register] Registration failed for %s: %v", req.Email, err)
		return pkgUtils.ConflictResponse(c, err.Error())
	}

	utils.InfoLogger.Printf("[Register] User registered successfully: %s (ID: %d)", user.Email, user.ID)
	return pkgUtils.CreatedResponse(c, "User registered successfully", user.GetPublicUser())
}

// Login godoc
//
//	@Summary		User login
//	@Description	Authenticate user and return access and refresh tokens
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.LoginRequest							true	"Login credentials"
//	@Success		200		{object}	models.APIResponse{data=dto.LoginResponse}	"Login successful"
//	@Failure		400		{object}	models.APIResponse							"Invalid request or validation error"
//	@Failure		401		{object}	models.APIResponse							"Invalid credentials or inactive account"
//	@Router			/auth/login [post]
func Login(c *fiber.Ctx) error {
	var req dto.LoginRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		utils.ErrorLogger.Printf("[Login] Failed to parse request body: %v", err)
		return pkgUtils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate request using DTO's self-validation
	if err := req.Validate(); err != nil {
		utils.ErrorLogger.Printf("[Login] Validation failed: %v", err)
		return pkgUtils.BadRequestResponse(c, err.Error())
	}

	// Authenticate user
	authService := services.NewAuthService(database.GetDB())
	loginResp, err := authService.Login(&req)
	if err != nil {
		utils.ErrorLogger.Printf("[Login] Authentication failed for %s: %v", req.Email, err)
		return pkgUtils.UnauthorizedResponse(c, err.Error())
	}

	utils.InfoLogger.Printf("[Login] User logged in successfully: %s", req.Email)
	return pkgUtils.SuccessResponse(c, fiber.StatusOK, "Login successful", loginResp)
}

// RefreshToken godoc
//
//	@Summary		Refresh access token
//	@Description	Generate a new access token using a valid refresh token
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.RefreshTokenRequest								true	"Refresh token"
//	@Success		200		{object}	models.APIResponse{data=dto.RefreshTokenResponse}	"Token refreshed successfully"
//	@Failure		400		{object}	models.APIResponse									"Invalid request body"
//	@Failure		401		{object}	models.APIResponse									"Invalid or expired refresh token"
//	@Router			/auth/refresh [post]
func RefreshToken(c *fiber.Ctx) error {
	var req dto.RefreshTokenRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		utils.ErrorLogger.Printf("[RefreshToken] Failed to parse request body: %v", err)
		return pkgUtils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate request using DTO's self-validation
	if err := req.Validate(); err != nil {
		utils.ErrorLogger.Printf("[RefreshToken] Validation failed: %v", err)
		return pkgUtils.BadRequestResponse(c, err.Error())
	}

	// Refresh token
	authService := services.NewAuthService(database.GetDB())
	newAccessToken, err := authService.RefreshToken(req.RefreshToken)
	if err != nil {
		utils.ErrorLogger.Printf("[RefreshToken] Token refresh failed: %v", err)
		return pkgUtils.UnauthorizedResponse(c, err.Error())
	}

	utils.InfoLogger.Printf("[RefreshToken] Token refreshed successfully")
	return pkgUtils.SuccessResponse(c, fiber.StatusOK, "Token refreshed successfully", dto.RefreshTokenResponse{
		Token: newAccessToken,
	})
}

// GetProfile godoc
//
//	@Summary		Get user profile
//	@Description	Retrieve the authenticated user's profile information
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	models.APIResponse{data=dto.UserResponse}	"Profile retrieved successfully"
//	@Failure		401	{object}	models.APIResponse							"Unauthorized or invalid token"
//	@Failure		404	{object}	models.APIResponse							"User not found"
//	@Router			/user/profile [get]
func GetProfile(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		utils.ErrorLogger.Printf("[GetProfile] Failed to get user ID from context: %v", err)
		return pkgUtils.UnauthorizedResponse(c, "Invalid user")
	}

	// Get user
	authService := services.NewAuthService(database.GetDB())
	user, err := authService.GetUserByID(userID)
	if err != nil {
		utils.ErrorLogger.Printf("[GetProfile] Failed to get user profile (ID: %d): %v", userID, err)
		return pkgUtils.NotFoundResponse(c, err.Error())
	}

	utils.InfoLogger.Printf("[GetProfile] Profile retrieved successfully (ID: %d)", userID)
	return pkgUtils.SuccessResponse(c, fiber.StatusOK, "Profile retrieved successfully", user.GetPublicUser())
}

// UpdateProfile godoc
//
//	@Summary		Update user profile
//	@Description	Update the authenticated user's profile information
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		dto.UpdateProfileRequest					true	"Profile update data"
//	@Success		200		{object}	models.APIResponse{data=dto.UserResponse}	"Profile updated successfully"
//	@Failure		400		{object}	models.APIResponse							"Invalid request or validation error"
//	@Failure		401		{object}	models.APIResponse							"Unauthorized or invalid token"
//	@Failure		500		{object}	models.APIResponse							"Failed to update profile"
//	@Router			/user/profile [put]
func UpdateProfile(c *fiber.Ctx) error {
	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		utils.ErrorLogger.Printf("[UpdateProfile] Failed to get user ID from context: %v", err)
		return pkgUtils.UnauthorizedResponse(c, "Invalid user")
	}

	var req dto.UpdateProfileRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		utils.ErrorLogger.Printf("[UpdateProfile] Failed to parse request body: %v", err)
		return pkgUtils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate request using DTO's self-validation
	if err := req.Validate(); err != nil {
		utils.ErrorLogger.Printf("[UpdateProfile] Validation failed: %v", err)
		return pkgUtils.BadRequestResponse(c, err.Error())
	}

	// Update user
	authService := services.NewAuthService(database.GetDB())
	user, err := authService.UpdateUser(userID, req.Name)
	if err != nil {
		utils.ErrorLogger.Printf("[UpdateProfile] Failed to update profile (ID: %d): %v", userID, err)
		return pkgUtils.InternalErrorResponse(c, "Failed to update profile")
	}

	utils.InfoLogger.Printf("[UpdateProfile] Profile updated successfully (ID: %d)", userID)
	return pkgUtils.SuccessResponse(c, fiber.StatusOK, "Profile updated successfully", user.GetPublicUser())
}

// ChangePassword godoc
//
//	@Summary		Change user password
//	@Description	Change the authenticated user's password
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		dto.ChangePasswordRequest	true	"Password change data"
//	@Success		200		{object}	models.APIResponse			"Password changed successfully"
//	@Failure		400		{object}	models.APIResponse			"Invalid request or validation error"
//	@Failure		401		{object}	models.APIResponse			"Unauthorized or invalid old password"
//	@Router			/user/change-password [post]
func ChangePassword(c *fiber.Ctx) error {
	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		utils.ErrorLogger.Printf("[ChangePassword] Failed to get user ID from context: %v", err)
		return pkgUtils.UnauthorizedResponse(c, "Invalid user")
	}

	var req dto.ChangePasswordRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		utils.ErrorLogger.Printf("[ChangePassword] Failed to parse request body: %v", err)
		return pkgUtils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate request using DTO's self-validation
	if err := req.Validate(); err != nil {
		utils.ErrorLogger.Printf("[ChangePassword] Validation failed: %v", err)
		return pkgUtils.BadRequestResponse(c, err.Error())
	}

	// Change password
	authService := services.NewAuthService(database.GetDB())
	if err := authService.ChangePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		utils.ErrorLogger.Printf("[ChangePassword] Failed to change password (ID: %d): %v", userID, err)
		return pkgUtils.UnauthorizedResponse(c, err.Error())
	}

	utils.InfoLogger.Printf("[ChangePassword] Password changed successfully (ID: %d)", userID)
	return pkgUtils.SuccessResponse(c, fiber.StatusOK, "Password changed successfully", nil)
}
