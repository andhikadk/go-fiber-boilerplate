package config

import (
	"fmt"
	"net/url"
	"os"

	"go-fiber-boilerplate/pkg/utils"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func (c *Config) GetDatabaseURL() string {
	switch c.DBDriver {
	case "postgres":
		u := &url.URL{
			Scheme: "postgres",
			User:   url.UserPassword(c.DBUser, c.DBPassword),
			Host:   fmt.Sprintf("%s:%s", c.DBHost, c.DBPort),
			Path:   c.DBName,
		}
		q := u.Query()
		q.Set("sslmode", c.DBSSLMode)
		u.RawQuery = q.Encode()
		return u.String()
	case "sqlite":
		return c.DBName
	default:
		utils.Log("Config").Error("Unsupported database driver", "driver", c.DBDriver)
		os.Exit(1)
		return ""
	}
}

func (c *Config) GetDialector() gorm.Dialector {
	switch c.DBDriver {
	case "postgres":
		return postgres.Open(c.GetDatabaseURL())
	case "sqlite":
		return sqlite.Open(c.GetDatabaseURL())
	default:
		utils.Log("Config").Error("Unsupported database driver", "driver", c.DBDriver)
		os.Exit(1)
		return nil
	}
}

func (c *Config) GetGormLogLevel() logger.LogLevel {
	switch c.LogLevel {
	case "debug":
		return logger.Info
	case "info":
		return logger.Warn
	case "error":
		return logger.Error
	default:
		return logger.Silent
	}
}
