package middleware

import (
	"net"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"go-fiber-boilerplate/pkg/utils"
)

var limiterStorage fiber.Storage

func InitLimiterStorage(store fiber.Storage) {
	limiterStorage = store
}

func NewGlobalLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Storage:    limiterStorage,
		Max:        200,
		Expiration: time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			ip := c.IP()
			parsed := net.ParseIP(ip)
			if parsed != nil && parsed.To4() == nil {
				mask := net.CIDRMask(64, 128)
				return "global:" + parsed.Mask(mask).String() + "/64"
			}
			return "global:" + ip
		},
		LimitReached: func(c *fiber.Ctx) error {
			utils.Log("Security").Warn("Global rate limit exceeded", "ip", c.IP(), "path", c.Path())
			return utils.TooManyRequestsResponse(c, "Too many requests, please slow down")
		},
	})
}

func NewAuthLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Storage:    limiterStorage,
		Max:        10,
		Expiration: time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return "auth:" + c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			utils.Log("Auth").Warn("Auth rate limit exceeded", "ip", c.IP(), "path", c.Path())
			return utils.TooManyRequestsResponse(c, "Too many login/registration attempts, please try again in 1 minute")
		},
	})
}

func NewPublicLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Storage:    limiterStorage,
		Max:        60,
		Expiration: time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return "public:" + c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			utils.Log("Public").Warn("Public rate limit exceeded", "ip", c.IP(), "path", c.Path())
			return utils.TooManyRequestsResponse(c, "Too many requests, please slow down")
		},
	})
}
