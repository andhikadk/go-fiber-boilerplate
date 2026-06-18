package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"go-fiber-boilerplate/pkg/utils"
)

type Config struct {
	Port         string
	Env          string
	AppName      string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration

	DBDriver          string
	DBHost            string
	DBPort            string
	DBUser            string
	DBPassword        string
	DBName            string
	DBSSLMode         string
	DBMaxOpenConns    int
	DBMaxIdleConns    int
	DBConnMaxLifetime time.Duration

	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int
	CacheEnabled  bool
	CacheTTL      time.Duration

	JWTSecret        string
	JWTExpiry        time.Duration
	JWTRefreshExpiry time.Duration

	CORSAllowedOrigins string
	CORSAllowedMethods string
	CORSAllowedHeaders string

	LogLevel         string
	LogHTTPBody      bool
	LogBodyMaxBytes  int
	LogRetentionDays int
	LogHealthSampleN int
	LogQuiet         bool

	SentryDSN              string
	SentryLogLevel         string
	SentryTracesSampleRate float64

	SMTPHost      string
	SMTPPort      int
	SMTPUser      string
	SMTPPassword  string
	SMTPFromName  string
	SMTPFromEmail string
}

var AppConfig *Config

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		utils.Log("Config").Info("Using system environment variables")
	}

	cfg := &Config{
		Port:         getEnv("PORT", "4000"),
		Env:          getEnv("ENV", "development"),
		AppName:      getEnv("APP_NAME", "Go Fiber Boilerplate API"),
		ReadTimeout:  parseDuration(getEnv("READ_TIMEOUT", "10s")),
		WriteTimeout: parseDuration(getEnv("WRITE_TIMEOUT", "10s")),
		IdleTimeout:  parseDuration(getEnv("IDLE_TIMEOUT", "60s")),

		DBDriver:          getEnv("DB_DRIVER", "postgres"),
		DBHost:            getEnv("DB_HOST", "localhost"),
		DBPort:            getEnv("DB_PORT", "5432"),
		DBUser:            getEnv("DB_USER", "postgres"),
		DBPassword:        getEnv("DB_PASSWORD", "postgres"),
		DBName:            getEnv("DB_NAME", "fiber_boilerplate"),
		DBSSLMode:         getEnv("DB_SSL_MODE", "disable"),
		DBMaxOpenConns:    parseInt(getEnv("DB_MAX_OPEN_CONNS", "25")),
		DBMaxIdleConns:    parseInt(getEnv("DB_MAX_IDLE_CONNS", "5")),
		DBConnMaxLifetime: parseDuration(getEnv("DB_CONN_MAX_LIFETIME", "5m")),

		RedisHost:     getEnv("REDIS_HOST", ""),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       parseInt(getEnv("REDIS_DB", "0")),
		CacheEnabled:  parseBool(getEnv("CACHE_ENABLED", "true")),
		CacheTTL:      parseDuration(getEnv("CACHE_TTL", "5m")),

		JWTSecret:        getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-this-in-production"),
		JWTExpiry:        parseDuration(getEnv("JWT_EXPIRY", "15m")),
		JWTRefreshExpiry: parseDuration(getEnv("JWT_REFRESH_EXPIRY", "168h")),

		CORSAllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:4000,http://localhost:8080"),
		CORSAllowedMethods: getEnv("CORS_ALLOWED_METHODS", "GET,POST,PUT,PATCH,DELETE,OPTIONS"),
		CORSAllowedHeaders: getEnv("CORS_ALLOWED_HEADERS", "Content-Type,Authorization"),

		LogLevel:         getEnv("LOG_LEVEL", "info"),
		LogHTTPBody:      parseBool(getEnv("LOG_HTTP_BODY", "true")),
		LogBodyMaxBytes:  parseInt(getEnv("LOG_BODY_MAX_BYTES", "8192")),
		LogRetentionDays: parseInt(getEnv("LOG_RETENTION_DAYS", "14")),
		LogHealthSampleN: parseInt(getEnv("LOG_HEALTH_SAMPLE_N", "20")),
		LogQuiet:         parseBool(getEnv("LOG_QUIET", "true")),

		SentryDSN:              getEnv("SENTRY_DSN", ""),
		SentryLogLevel:         getEnv("SENTRY_LOG_LEVEL", "info"),
		SentryTracesSampleRate: parseFloat(getEnv("SENTRY_TRACES_SAMPLE_RATE", "0")),

		SMTPHost:      getEnv("SMTP_HOST", ""),
		SMTPPort:      parseInt(getEnv("SMTP_PORT", "587")),
		SMTPUser:      getEnv("SMTP_USER", ""),
		SMTPPassword:  getEnv("SMTP_PASSWORD", ""),
		SMTPFromName:  getEnv("SMTP_FROM_NAME", "Go Fiber Boilerplate"),
		SMTPFromEmail: getEnv("SMTP_FROM_EMAIL", ""),
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	AppConfig = cfg
	return cfg, nil
}

func (c *Config) Validate() error {
	if c.JWTSecret == "your-super-secret-jwt-key-change-this-in-production" {
		return fmt.Errorf("JWT_SECRET must be changed from the default value")
	}
	if len(c.JWTSecret) < 32 {
		return fmt.Errorf("JWT_SECRET must be at least 32 characters long")
	}
	if c.DBDriver != "postgres" && c.DBDriver != "sqlite" {
		return fmt.Errorf("DB_DRIVER must be either 'postgres' or 'sqlite'")
	}
	if c.DBMaxIdleConns > c.DBMaxOpenConns && c.DBMaxOpenConns > 0 {
		return fmt.Errorf("DB_MAX_IDLE_CONNS (%d) must not exceed DB_MAX_OPEN_CONNS (%d)", c.DBMaxIdleConns, c.DBMaxOpenConns)
	}
	if c.IsProduction() && c.DBSSLMode == "disable" {
		utils.Log("Config").Warn("DB_SSL_MODE is disabled in production")
	}
	return nil
}

func (c *Config) IsDevelopment() bool {
	return c.Env == "development"
}

func (c *Config) IsProduction() bool {
	return c.Env == "production"
}

func (c *Config) RedisEnabled() bool {
	return strings.TrimSpace(c.RedisHost) != ""
}

func (c *Config) RedisAddr() string {
	return strings.TrimSpace(c.RedisHost) + ":" + strings.TrimSpace(c.RedisPort)
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func parseDuration(s string) time.Duration {
	duration, err := time.ParseDuration(s)
	if err != nil {
		return 10 * time.Second
	}
	return duration
}

func parseInt(s string) int {
	value, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil {
		return 0
	}
	return value
}

func parseFloat(s string) float64 {
	value, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	if err != nil {
		return 0
	}
	return value
}

func parseBool(s string) bool {
	value, err := strconv.ParseBool(strings.TrimSpace(s))
	if err != nil {
		return false
	}
	return value
}
