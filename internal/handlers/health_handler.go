package handlers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"go-fiber-boilerplate/internal/database"
	"go-fiber-boilerplate/pkg/utils"
)

// HealthCheck godoc
//
//	@Summary		Health check
//	@Description	Check API and database health
//	@Tags			Health
//	@Produce		json
//	@Success		200	{object}	models.APIResponse
//	@Failure		503	{object}	models.APIResponse
//	@Router			/health [get]
func HealthCheck(c *fiber.Ctx) error {
	start := time.Now()
	db := database.GetDB()
	if db == nil {
		return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "Database is not reachable")
	}

	var pingResult int
	if err := db.Raw("SELECT 1").Scan(&pingResult).Error; err != nil || pingResult != 1 {
		if err != nil {
			utils.LogCtx(c.UserContext(), "Health").Error("Database health check failed", "error", err)
		}
		return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "Database is not reachable")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "API is running", fiber.Map{
		"response_time": fmt.Sprintf("%.2fms", float64(time.Since(start))/float64(time.Millisecond)),
	})
}
