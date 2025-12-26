package middleware

import (
	"fleetify/pkg/jwt"

	"github.com/gofiber/fiber/v2"
)

func Admin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, ok := c.Locals("user").(*jwt.Claims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   true,
				"message": "Unauthorized",
			})
		}

		if claims.Role != "ADMIN" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":   true,
				"message": "Admin access required",
			})
		}

		return c.Next()
	}
}

func ItemModifyAccess() fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, ok := c.Locals("user").(*jwt.Claims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   true,
				"message": "Unauthorized",
			})
		}

		allowedRoles := map[string]bool{
			"ADMIN":     true,
			"MANAGER":   true,
			"SUPPLIERS": true,
		}

		if !allowedRoles[claims.Role] {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":   true,
				"message": "Access denied. Only ADMIN, MANAGER, and SUPPLIERS can modify items",
			})
		}

		return c.Next()
	}
}
