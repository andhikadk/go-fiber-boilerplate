package middleware

import (
	"github.com/gofiber/fiber/v2"
	"go-fiber-boilerplate/pkg/utils"
)

func ErrorHandlingMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := c.Next(); err != nil {
			return handleError(c, err)
		}
		return nil
	}
}

func handleError(c *fiber.Ctx, err error) error {
	if e, ok := err.(*fiber.Error); ok {
		return utils.ErrorResponse(c, e.Code, e.Message)
	}
	utils.LogCtx(c.UserContext(), "Error").Error("Unhandled request error", "error", err)
	return utils.InternalErrorResponse(c, "Internal Server Error")
}
