package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"fleetify/internal/config"
)

// CORS  config
func CORS() fiber.Handler {
	return func(c *fiber.Ctx) error {
		cfg := config.AppConfig.CORS

		origin := c.Get("Origin")
		allowedOrigins := strings.Split(cfg.AllowedOrigins, ",")
		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if strings.TrimSpace(allowedOrigin) == origin || origin == "" {
				allowed = true
				break
			}
		}

		if allowed {
			c.Set("Access-Control-Allow-Origin", origin)
		}

		c.Set("Access-Control-Allow-Credentials", "true")
		c.Set("Access-Control-Allow-Methods", cfg.AllowedMethods)
		c.Set("Access-Control-Allow-Headers", cfg.AllowedHeaders)

		if c.Method() == "OPTIONS" {
			return c.SendStatus(fiber.StatusNoContent)
		}

		return c.Next()
	}
}

