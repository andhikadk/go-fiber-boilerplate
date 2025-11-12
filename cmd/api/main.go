package main

import (
	"flag"
	"fmt"
	"log"

	"go-fiber-boilerplate/assets"
	"go-fiber-boilerplate/config"
	"go-fiber-boilerplate/internal/database"
	"go-fiber-boilerplate/internal/middleware"
	"go-fiber-boilerplate/internal/routes"
	"go-fiber-boilerplate/internal/utils"

	_ "go-fiber-boilerplate/docs" // Import generated docs

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/gorm"
)

//	@title						Go Fiber REST API Boilerplate
//	@version					1.0
//	@description				A production-ready REST API boilerplate built with Go Fiber, GORM, JWT authentication, and PostgreSQL/SQLite
//	@termsOfService				http://swagger.io/terms/
//	@contact.name				API Support
//	@contact.url				https://github.com/yourusername/go-fiber-boilerplate
//	@contact.email				support@example.com
//	@license.name				MIT
//	@license.url				https://opensource.org/licenses/MIT
//	@host						localhost:4000
//	@BasePath					/
//	@schemes					http https
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Type "Bearer" followed by a space and JWT token.

func main() {
	// Parse command line flags
	migrateCmd := flag.String("migrate", "", "Run migrations (use: -migrate=auto or -migrate=sql)")
	seedCmd := flag.Bool("seed", false, "Seed database with sample data")
	statusCmd := flag.Bool("status", false, "Show migration status")
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	if err := utils.InitLogger(); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	// Initialize database
	db, err := database.Initialize(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Handle migration commands
	if *migrateCmd != "" {
		if *migrateCmd == "sql" {
			log.Println("Running SQL migrations from embedded files...")
			if err := database.MigrateFromFS(db, assets.MigrationsFS); err != nil {
				log.Fatalf("Migration failed: %v", err)
			}
		} else {
			// Default to AutoMigrate for development
			if err := database.Migrate(db, cfg); err != nil {
				log.Fatalf("Migration failed: %v", err)
			}
		}
		log.Println("Migrations completed successfully")
		return
	}

	// Handle seed command
	if *seedCmd {
		log.Println("Seeding database...")
		if err := database.SeedFromFS(db, assets.MigrationsFS); err != nil {
			log.Fatalf("Seeding failed: %v", err)
		}
		log.Println("Seeding completed successfully")
		return
	}

	// Handle status command
	if *statusCmd {
		showMigrationStatus(db)
		return
	}

	// Run normal migrations (AutoMigrate for dev, SQL for production)
	if err := database.Migrate(db, cfg); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:           cfg.AppName,
		ReadTimeout:       cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
		EnablePrintRoutes: cfg.IsDevelopment(),
	})

	// Setup global middleware
	setupMiddleware(app, cfg)

	// Setup routes
	routes.SetupRoutes(app)

	// Start server
	startServer(app, cfg)
}

// showMigrationStatus displays migration status
func showMigrationStatus(db *gorm.DB) {
	fmt.Println("\n=== Migration Status ===")

	migrator := database.NewMigrator(db)
	migrations, err := migrator.GetAppliedMigrations()
	if err != nil {
		log.Fatalf("Failed to get migration status: %v", err)
	}

	if len(migrations) == 0 {
		fmt.Println("No migrations applied yet")
	} else {
		fmt.Println("Applied migrations:")
		for _, m := range migrations {
			fmt.Printf("  ✓ %s\n", m)
		}
	}

	seeder := database.NewSeeder(db)
	seeds, err := seeder.GetAppliedSeeds()

	fmt.Println("\nApplied seeds:")
	if err != nil {
		// Table doesn't exist yet, no seeds applied
		fmt.Println("  No seeds applied yet")
	} else if len(seeds) == 0 {
		fmt.Println("  No seeds applied yet")
	} else {
		for _, s := range seeds {
			fmt.Printf("  ✓ %s\n", s)
		}
	}
	fmt.Println()
}

// setupMiddleware configures global middleware
func setupMiddleware(app *fiber.App, cfg *config.Config) {
	// Logger middleware
	app.Use(fiberLogger.New(fiberLogger.Config{
		Format: "[${time}] ${status} - ${method} ${path} (${latency})\n",
	}))

	// Recovery middleware (panic recovery)
	app.Use(recover.New())

	// CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.CORSAllowedOrigins,
		AllowMethods: cfg.CORSAllowedMethods,
		AllowHeaders: cfg.CORSAllowedHeaders,
	}))

	// Security middleware (Helmet)
	app.Use(helmet.New())

	// Compression middleware
	app.Use(compress.New(compress.Config{
		Level: compress.LevelDefault,
	}))

	// Error handling middleware
	app.Use(middleware.ErrorHandlingMiddleware())
}

// startServer starts the Fiber server
func startServer(app *fiber.App, cfg *config.Config) {
	address := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Starting %s on %s [%s mode]", cfg.AppName, address, cfg.Env)

	if err := app.Listen(address); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
