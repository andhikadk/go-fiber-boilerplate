package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"syscall"

	"go-fiber-boilerplate/assets"
	"go-fiber-boilerplate/config"
	"go-fiber-boilerplate/internal/cache"
	"go-fiber-boilerplate/internal/database"
	"go-fiber-boilerplate/internal/middleware"
	"go-fiber-boilerplate/internal/routes"
	"go-fiber-boilerplate/pkg/utils"

	_ "go-fiber-boilerplate/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	fiberRecover "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"gorm.io/gorm"
)

//	@title						Go Fiber Boilerplate API
//	@version					2.0
//	@description				A production-ready REST API template built with Fiber, GORM, PostgreSQL, JWT, SQL migrations, and Scalar docs.
//	@host						localhost:4000
//	@BasePath					/api
//	@schemes					http https
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Type "Bearer" followed by a space and JWT token.

func main() {
	migrateCmd := flag.String("migrate", "", "Run SQL migrations (use: -migrate=run, -migrate=fresh, or -migrate=status)")
	seedCmd := flag.Bool("seed", false, "Seed database with sample data")
	flag.Parse()

	utils.InitLogger()

	cfg, err := config.LoadConfig()
	if err != nil {
		utils.Log("App").Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}
	utils.SetLogLevel(cfg.LogLevel)
	utils.SetQuiet(cfg.LogQuiet)
	utils.CleanupOldLogs("logs/app", cfg.LogRetentionDays)

	db, err := database.Initialize(cfg)
	if err != nil {
		utils.Log("App").Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer database.Close()

	if *migrateCmd != "" {
		handleMigrationCommand(db, cfg, *migrateCmd)
		return
	}

	if *seedCmd {
		utils.Log("App").Info("Seeding database")
		if err := database.SeedFromFS(db, assets.SeedsFS); err != nil {
			utils.Log("App").Error("Seeding failed", "error", err)
			os.Exit(1)
		}
		utils.Log("App").Info("Seeding completed successfully")
		return
	}

	utils.Log("App").Info("Checking and running pending migrations")
	if err := database.MigrateFromFS(db, assets.MigrationsFS); err != nil {
		utils.Log("App").Error("Failed to run migrations", "error", err)
		os.Exit(1)
	}

	app := fiber.New(fiber.Config{
		AppName:           cfg.AppName,
		BodyLimit:         10 * 1024 * 1024,
		ReadTimeout:       cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
		EnablePrintRoutes: cfg.IsDevelopment(),
	})

	cacheClient := cache.New(cfg)
	defer cacheClient.Close()
	middleware.InitLimiterStorage(cache.NewLimiterStorage(cfg))

	setupMiddleware(app, cfg)
	routes.SetupRoutes(app, cacheClient)
	startServer(app, cfg)
}

func handleMigrationCommand(db *gorm.DB, cfg *config.Config, cmd string) {
	switch cmd {
	case "fresh":
		if !cfg.IsDevelopment() {
			utils.Log("App").Error("migrate=fresh is only allowed in development mode", "env", cfg.Env)
			os.Exit(1)
		}
		migrator := database.NewMigrator(db)
		if err := migrator.FreshMigrate(); err != nil {
			utils.Log("App").Error("Fresh migration failed", "error", err)
			os.Exit(1)
		}
		if err := database.MigrateFromFS(db, assets.MigrationsFS); err != nil {
			utils.Log("App").Error("Migration failed", "error", err)
			os.Exit(1)
		}
	case "status":
		showMigrationStatus(db)
	default:
		if err := database.MigrateFromFS(db, assets.MigrationsFS); err != nil {
			utils.Log("App").Error("Migration failed", "error", err)
			os.Exit(1)
		}
	}
}

func showMigrationStatus(db *gorm.DB) {
	fmt.Println("\n=== Migration Status ===")
	migrator := database.NewMigrator(db)
	migrations, err := migrator.GetAppliedMigrations()
	if err != nil {
		utils.Log("App").Error("Failed to get migration status", "error", err)
		os.Exit(1)
	}

	if len(migrations) == 0 {
		fmt.Println("No migrations applied yet")
	} else {
		fmt.Println("Applied migrations:")
		for _, migration := range migrations {
			fmt.Printf("  %s\n", migration)
		}
	}

	fmt.Println("\nAvailable seeds:")
	entries, err := assets.SeedsFS.ReadDir("migrations/seeds")
	if err != nil || len(entries) == 0 {
		fmt.Println("  No seed files found")
	} else {
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
				fmt.Printf("  %s\n", entry.Name())
			}
		}
	}
	fmt.Println()
}

func setupMiddleware(app *fiber.App, cfg *config.Config) {
	app.Use(requestid.New())
	app.Use(middleware.RequestContext())
	app.Use(fiberRecover.New(fiberRecover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			utils.Log("Recover").Error("panic recovered",
				"request_id", c.Locals("requestid"),
				"method", c.Method(),
				"path", c.Path(),
				"panic", fmt.Sprintf("%v", e),
				"stack", string(debug.Stack()),
			)
		},
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORSAllowedOrigins,
		AllowMethods:     cfg.CORSAllowedMethods,
		AllowHeaders:     cfg.CORSAllowedHeaders,
		AllowCredentials: true,
	}))
	app.Use(helmet.New())
	app.Use(middleware.NewGlobalLimiter())
	app.Use(compress.New(compress.Config{Level: compress.LevelDefault}))
	app.Use(middleware.AccessLog(cfg.LogHTTPBody, cfg.LogBodyMaxBytes, cfg.LogHealthSampleN))
	app.Use(middleware.ErrorHandlingMiddleware())
}

func startServer(app *fiber.App, cfg *config.Config) {
	address := fmt.Sprintf(":%s", cfg.Port)
	utils.Log("App").Info("Starting server", "app_name", cfg.AppName, "address", address, "mode", cfg.Env)

	go func() {
		if err := app.Listen(address); err != nil {
			utils.Log("App").Error("Server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	utils.Log("App").Info("Shutting down server gracefully")
	if err := app.Shutdown(); err != nil {
		utils.Log("App").Error("Server shutdown error", "error", err)
	}
	utils.Log("App").Info("Server stopped")
}
