package database

import (
	"embed"

	"go-fiber-boilerplate/config"
	"go-fiber-boilerplate/pkg/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Initialize(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(
		cfg.GetDialector(),
		&gorm.Config{Logger: logger.Default.LogMode(cfg.GetGormLogLevel())},
	)
	if err != nil {
		utils.Log("Database").Error("Failed to connect to database", "error", err)
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		utils.Log("Database").Error("Failed to get underlying sql.DB", "error", err)
		return nil, err
	}
	sqlDB.SetMaxOpenConns(cfg.DBMaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.DBMaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.DBConnMaxLifetime)

	utils.Log("Database").Info("Database connection established", "max_open", cfg.DBMaxOpenConns, "max_idle", cfg.DBMaxIdleConns)
	DB = db
	return db, nil
}

func MigrateFromFS(db *gorm.DB, migrations embed.FS) error {
	migrator := NewMigrator(db)
	return migrator.RunMigrationsFromFS(migrations)
}

func SeedFromFS(db *gorm.DB, seeds embed.FS) error {
	seeder := NewSeeder(db)
	return seeder.SeedFromFS(seeds)
}

func Close() error {
	if DB == nil {
		return nil
	}
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func GetDB() *gorm.DB {
	return DB
}
