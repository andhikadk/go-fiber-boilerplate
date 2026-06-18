package routes

import (
	"go-fiber-boilerplate/config"
	"go-fiber-boilerplate/internal/cache"
	"go-fiber-boilerplate/internal/database"
	"go-fiber-boilerplate/internal/handlers"
	"go-fiber-boilerplate/internal/middleware"
	"go-fiber-boilerplate/internal/services"
	"go-fiber-boilerplate/pkg/utils"

	scalar "github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, _ *cache.Client) {
	emailService := services.NewNoopEmailService()
	if config.AppConfig.SMTPHost == "" {
		utils.Log("Routes").Warn("SMTP Host not configured, email service disabled")
	} else {
		utils.Log("Routes").Warn("SMTP Host configured, but template only wires the no-op EmailService skeleton")
	}
	_ = services.NewNoopStorageService()

	authService := services.NewAuthService(database.GetDB(), emailService)
	userService := services.NewUserService(database.GetDB())
	resourceService := services.NewResourceService(database.GetDB())

	authHandler := handlers.NewAuth(authService)
	userHandler := handlers.NewUser(userService)
	resourceHandler := handlers.NewResource(resourceService)

	app.Get("/health", handlers.HealthCheck)

	if !config.AppConfig.IsProduction() {
		app.Static("/swagger.json", "./docs/swagger.json")
		app.Get("/docs", func(c *fiber.Ctx) error {
			scheme := "http"
			if c.Protocol() == "https" {
				scheme = "https"
			}
			specURL := scheme + "://" + c.Hostname() + "/swagger.json"
			html, err := scalar.ApiReferenceHTML(&scalar.Options{
				SpecURL:     specURL,
				Theme:       scalar.ThemeDefault,
				Layout:      scalar.LayoutModern,
				DarkMode:    true,
				ShowSidebar: true,
			})
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
			}
			c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
			return c.SendString(html)
		})
	}

	api := app.Group("/api")

	authGroup := api.Group("/auth")
	authGroup.Use(middleware.NewAuthLimiter())
	{
		authGroup.Post("/register", authHandler.Register)
		authGroup.Post("/login", authHandler.Login)
		authGroup.Post("/refresh", authHandler.RefreshToken)
		authGroup.Post("/forgot-password", authHandler.ForgotPassword)
		authGroup.Post("/reset-password", authHandler.ResetPassword)
	}

	userGroup := api.Group("/user")
	userGroup.Use(middleware.AuthMiddleware())
	{
		userGroup.Get("/profile", userHandler.GetProfile)
		userGroup.Put("/profile", userHandler.UpdateProfile)
		userGroup.Post("/change-password", userHandler.ChangePassword)
	}

	resourcesGroup := api.Group("/resources")
	resourcesGroup.Use(middleware.AuthMiddleware())
	{
		resourcesGroup.Get("/", resourceHandler.ListResources)
		resourcesGroup.Post("/", resourceHandler.CreateResource)
		resourcesGroup.Get("/:id", resourceHandler.GetResource)
		resourcesGroup.Put("/:id", resourceHandler.UpdateResource)
		resourcesGroup.Delete("/:id", resourceHandler.DeleteResource)
	}

	app.Use(func(c *fiber.Ctx) error {
		return utils.NotFoundResponse(c, "endpoint not found")
	})
}
