package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"go-fiber-boilerplate/config"
	"go-fiber-boilerplate/pkg/jwt"
	"go-fiber-boilerplate/pkg/utils"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return utils.UnauthorizedResponse(c, "missing authorization header")
		}

		token, err := jwt.ExtractTokenFromHeader(authHeader)
		if err != nil {
			return utils.UnauthorizedResponse(c, "invalid authorization header format")
		}

		tm := jwt.NewTokenManager(config.AppConfig.JWTSecret)
		claims, err := tm.ValidateAccessToken(token)
		if err != nil {
			return utils.UnauthorizedResponse(c, "invalid or expired token")
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("role", claims.Role)
		return c.Next()
	}
}

func OptionalAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Next()
		}

		token, err := jwt.ExtractTokenFromHeader(authHeader)
		if err != nil {
			return c.Next()
		}

		tm := jwt.NewTokenManager(config.AppConfig.JWTSecret)
		claims, err := tm.ValidateAccessToken(token)
		if err != nil {
			return c.Next()
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("role", claims.Role)
		return c.Next()
	}
}

func AdminMiddleware() fiber.Handler {
	return RequireRoles("admin")
}

func RequireRoles(roles ...string) fiber.Handler {
	allowed := make(map[string]bool, len(roles))
	for _, role := range roles {
		allowed[strings.ToLower(role)] = true
	}
	return func(c *fiber.Ctx) error {
		role, _ := c.Locals("role").(string)
		if !allowed[strings.ToLower(role)] {
			return utils.ForbiddenResponse(c, "access denied")
		}
		return c.Next()
	}
}

func GetUserIDFromContext(c *fiber.Ctx) (uint, error) {
	userID := c.Locals("user_id")
	if userID == nil {
		return 0, fiber.ErrUnauthorized
	}
	id, ok := userID.(uint)
	if !ok {
		return 0, fiber.ErrUnauthorized
	}
	return id, nil
}

func GetEmailFromContext(c *fiber.Ctx) string {
	email := c.Locals("email")
	if email == nil {
		return ""
	}
	s, ok := email.(string)
	if !ok {
		return ""
	}
	return s
}
