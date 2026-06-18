package handlers

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go-fiber-boilerplate/internal/dto"
	"go-fiber-boilerplate/internal/middleware"
	"go-fiber-boilerplate/internal/services"
	"go-fiber-boilerplate/pkg/utils"
)

type Resource struct {
	resourceService services.ResourceService
}

func NewResource(resourceService services.ResourceService) *Resource {
	return &Resource{resourceService: resourceService}
}

// ListResources godoc
//
//	@Summary		List resources
//	@Tags			Resources
//	@Produce		json
//	@Security		BearerAuth
//	@Param			page	query		int	false	"Page number"
//	@Param			limit	query		int	false	"Items per page"
//	@Success		200		{object}	models.PaginatedResponse
//	@Router			/resources [get]
func (h *Resource) ListResources(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	resources, total, err := h.resourceService.ListResources(page, limit)
	if err != nil {
		utils.LogCtx(c.UserContext(), "Resource").Error("List resources failed", "error", err)
		return utils.InternalErrorResponse(c, "Failed to list resources")
	}
	return utils.PaginatedResponse(c, "Resources retrieved successfully", resources, page, limit, total)
}

// GetResource godoc
//
//	@Summary		Get resource
//	@Tags			Resources
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Resource ID"
//	@Success		200	{object}	models.APIResponse
//	@Failure		404	{object}	models.APIResponse
//	@Router			/resources/{id} [get]
func (h *Resource) GetResource(c *fiber.Ctx) error {
	id, err := parseIDParam(c, "id")
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid resource ID")
	}
	resource, err := h.resourceService.GetResource(id)
	if err != nil {
		if errors.Is(err, services.ErrResourceNotFound) {
			return utils.NotFoundResponse(c, "Resource not found")
		}
		utils.LogCtx(c.UserContext(), "Resource").Error("Get resource failed", "id", id, "error", err)
		return utils.InternalErrorResponse(c, "Failed to get resource")
	}
	return utils.SuccessResponse(c, fiber.StatusOK, "Resource retrieved successfully", resource)
}

// CreateResource godoc
//
//	@Summary		Create resource
//	@Tags			Resources
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		dto.CreateResourceRequest	true	"Resource creation data"
//	@Success		201		{object}	models.APIResponse
//	@Router			/resources [post]
func (h *Resource) CreateResource(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "Invalid user")
	}
	var req dto.CreateResourceRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}
	if err := req.Validate(); err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}
	resource, err := h.resourceService.CreateResource(userID, &req)
	if err != nil {
		utils.LogCtx(c.UserContext(), "Resource").Error("Create resource failed", "error", err)
		return utils.InternalErrorResponse(c, "Failed to create resource")
	}
	return utils.CreatedResponse(c, "Resource created successfully", resource)
}

// UpdateResource godoc
//
//	@Summary		Update resource
//	@Tags			Resources
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int							true	"Resource ID"
//	@Param			request	body		dto.UpdateResourceRequest	true	"Resource update data"
//	@Success		200		{object}	models.APIResponse
//	@Router			/resources/{id} [put]
func (h *Resource) UpdateResource(c *fiber.Ctx) error {
	id, err := parseIDParam(c, "id")
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid resource ID")
	}
	var req dto.UpdateResourceRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}
	if err := req.Validate(); err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}
	resource, err := h.resourceService.UpdateResource(id, &req)
	if err != nil {
		if errors.Is(err, services.ErrResourceNotFound) {
			return utils.NotFoundResponse(c, "Resource not found")
		}
		utils.LogCtx(c.UserContext(), "Resource").Error("Update resource failed", "id", id, "error", err)
		return utils.InternalErrorResponse(c, "Failed to update resource")
	}
	return utils.SuccessResponse(c, fiber.StatusOK, "Resource updated successfully", resource)
}

// DeleteResource godoc
//
//	@Summary		Delete resource
//	@Tags			Resources
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Resource ID"
//	@Success		200	{object}	models.APIResponse
//	@Router			/resources/{id} [delete]
func (h *Resource) DeleteResource(c *fiber.Ctx) error {
	id, err := parseIDParam(c, "id")
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid resource ID")
	}
	if err := h.resourceService.DeleteResource(id); err != nil {
		if errors.Is(err, services.ErrResourceNotFound) {
			return utils.NotFoundResponse(c, "Resource not found")
		}
		utils.LogCtx(c.UserContext(), "Resource").Error("Delete resource failed", "id", id, "error", err)
		return utils.InternalErrorResponse(c, "Failed to delete resource")
	}
	return utils.SuccessResponse(c, fiber.StatusOK, "Resource deleted successfully", nil)
}

func parseIDParam(c *fiber.Ctx, name string) (uint, error) {
	id, err := strconv.ParseUint(c.Params(name), 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}
