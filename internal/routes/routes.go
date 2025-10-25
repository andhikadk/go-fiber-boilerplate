package routes

import (
	"go-fiber-boilerplate/internal/handlers"
	"go-fiber-boilerplate/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes configures all application routes
func SetupRoutes(app *fiber.App) {
	// Health check routes (public)
	app.Get("/health", handlers.HealthCheck)
	app.Get("/ready", handlers.ReadinessCheck)

	// Auth routes (public)
	authGroup := app.Group("/auth")
	authGroup.Post("/register", handlers.Register)
	authGroup.Post("/login", handlers.Login)
	authGroup.Post("/refresh", handlers.RefreshToken)

	// Protected routes (require authentication)
	// User routes
	userGroup := app.Group("/user")
	userGroup.Use(middleware.AuthMiddleware())
	{
		userGroup.Get("/profile", handlers.GetProfile)
		userGroup.Put("/profile", handlers.UpdateProfile)
		userGroup.Post("/change-password", handlers.ChangePassword)
	}

	// API routes
	apiGroup := app.Group("/api")
	apiGroup.Use(middleware.AuthMiddleware())
	{
		// Books routes
		booksGroup := apiGroup.Group("/books")
		{
			booksGroup.Get("/", handlers.GetBooks)
			booksGroup.Get("/search", handlers.SearchBooks)
			booksGroup.Get("/:id", handlers.GetBook)
			booksGroup.Post("/", handlers.CreateBook)
			booksGroup.Put("/:id", handlers.UpdateBook)
			booksGroup.Delete("/:id", handlers.DeleteBook)
		}
	}

	// 404 handler
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  fiber.StatusNotFound,
			"message": "endpoint not found",
		})
	})
}
