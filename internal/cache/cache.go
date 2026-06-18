package cache

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	redisstore "github.com/gofiber/storage/redis/v2"
	"github.com/redis/go-redis/v9"
	"go-fiber-boilerplate/config"
	"go-fiber-boilerplate/pkg/utils"
)

type Client struct {
	rdb     *redis.Client
	enabled bool
	ttl     time.Duration
}

func New(cfg *config.Config) *Client {
	if !cfg.RedisEnabled() || !cfg.CacheEnabled {
		utils.Log("Cache").Info("Cache disabled", "redis_configured", cfg.RedisEnabled(), "cache_enabled", cfg.CacheEnabled)
		return &Client{enabled: false, ttl: cfg.CacheTTL}
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr(),
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		utils.Log("Cache").Warn("Redis unreachable, cache disabled", "addr", cfg.RedisAddr(), "error", err)
		_ = rdb.Close()
		return &Client{enabled: false, ttl: cfg.CacheTTL}
	}

	utils.Log("Cache").Info("Cache enabled", "addr", cfg.RedisAddr())
	return &Client{rdb: rdb, enabled: true, ttl: cfg.CacheTTL}
}

func (c *Client) Enabled() bool {
	return c != nil && c.enabled
}

func (c *Client) TTL() time.Duration {
	if c == nil {
		return 0
	}
	return c.ttl
}

func (c *Client) GetJSON(ctx context.Context, key string, dest interface{}) bool {
	if !c.Enabled() {
		return false
	}
	data, err := c.rdb.Get(ctx, key).Bytes()
	if err != nil {
		if err != redis.Nil {
			utils.LogCtx(ctx, "Cache").Warn("Get failed", "key", key, "error", err)
		}
		return false
	}
	if err := json.Unmarshal(data, dest); err != nil {
		utils.LogCtx(ctx, "Cache").Warn("Unmarshal failed", "key", key, "error", err)
		return false
	}
	return true
}

func (c *Client) SetJSON(ctx context.Context, key string, val interface{}, ttl time.Duration) {
	if !c.Enabled() {
		return
	}
	data, err := json.Marshal(val)
	if err != nil {
		utils.LogCtx(ctx, "Cache").Warn("Marshal failed", "key", key, "error", err)
		return
	}
	if ttl <= 0 {
		ttl = c.ttl
	}
	if err := c.rdb.Set(ctx, key, data, ttl).Err(); err != nil {
		utils.LogCtx(ctx, "Cache").Warn("Set failed", "key", key, "error", err)
	}
}

func (c *Client) Delete(ctx context.Context, keys ...string) {
	if !c.Enabled() || len(keys) == 0 {
		return
	}
	if err := c.rdb.Del(ctx, keys...).Err(); err != nil {
		utils.LogCtx(ctx, "Cache").Warn("Delete failed", "keys", keys, "error", err)
	}
}

func (c *Client) Close() {
	if c.Enabled() {
		_ = c.rdb.Close()
	}
}

func NewLimiterStorage(cfg *config.Config) fiber.Storage {
	if !cfg.RedisEnabled() {
		return nil
	}
	port := 6379
	if p, err := strconv.Atoi(strings.TrimSpace(cfg.RedisPort)); err == nil {
		port = p
	}
	return redisstore.New(redisstore.Config{
		Host:     cfg.RedisHost,
		Port:     port,
		Password: cfg.RedisPassword,
		Database: cfg.RedisDB,
		Reset:    false,
	})
}
