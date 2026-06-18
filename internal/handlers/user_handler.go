package handlers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"go-fiber-boilerplate/internal/dto"
	"go-fiber-boilerplate/internal/middleware"
	"go-fiber-boilerplate/internal/services"
	"go-fiber-boilerplate/pkg/utils"
)

type User struct {
	userService services.UserService
}

func NewUser(userService services.UserService) *User {
	return &User{userService: userService}
}

// GetProfile godoc
//
//	@Summary		Get user profile
//	@Tags			Users
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	models.APIResponse	"Profile retrieved successfully"
//	@Failure		401	{object}	models.APIResponse	"Unauthorized"
//	@Failure		404	{object}	models.APIResponse	"User not found"
//	@Router			/user/profile [get]
func (h *User) GetProfile(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "Invalid user")
	}
	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			return utils.NotFoundResponse(c, "User not found")
		}
		utils.LogCtx(c.UserContext(), "User").Error("Get profile failed", "user_id", userID, "error", err)
		return utils.InternalErrorResponse(c, "Failed to get profile")
	}
	return utils.SuccessResponse(c, fiber.StatusOK, "Profile retrieved successfully", h.userService.GetUserResponse(user))
}

// UpdateProfile godoc
//
//	@Summary		Update user profile
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		dto.UpdateProfileRequest	true	"Profile update data"
//	@Success		200		{object}	models.APIResponse			"Profile updated successfully"
//	@Failure		400		{object}	models.APIResponse			"Invalid request"
//	@Router			/user/profile [put]
func (h *User) UpdateProfile(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "Invalid user")
	}
	var req dto.UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}
	if err := req.Validate(); err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}
	user, err := h.userService.UpdateProfile(userID, &req)
	if err != nil {
		utils.LogCtx(c.UserContext(), "User").Error("Update profile failed", "user_id", userID, "error", err)
		return utils.InternalErrorResponse(c, "Failed to update profile")
	}
	return utils.SuccessResponse(c, fiber.StatusOK, "Profile updated successfully", h.userService.GetUserResponse(user))
}

// ChangePassword godoc
//
//	@Summary		Change user password
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		dto.ChangePasswordRequest	true	"Password change data"
//	@Success		200		{object}	models.APIResponse			"Password changed successfully"
//	@Failure		400		{object}	models.APIResponse			"Invalid request"
//	@Failure		401		{object}	models.APIResponse			"Invalid current password"
//	@Router			/user/change-password [post]
func (h *User) ChangePassword(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "Invalid user")
	}
	var req dto.ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}
	if err := req.Validate(); err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}
	if err := h.userService.ChangePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		if errors.Is(err, services.ErrInvalidPassword) || errors.Is(err, services.ErrNoPasswordSet) {
			return utils.UnauthorizedResponse(c, err.Error())
		}
		utils.LogCtx(c.UserContext(), "User").Error("Change password failed", "user_id", userID, "error", err)
		return utils.InternalErrorResponse(c, "Failed to change password")
	}
	return utils.SuccessResponse(c, fiber.StatusOK, "Password changed successfully", nil)
}
