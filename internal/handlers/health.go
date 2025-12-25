package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"fleetify/internal/database"
)

// HealthCheckServer
func HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "ok",
		"message": "Server is running",
	})
}

// HealthCheckDB
func HealthCheckDB(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbStatus := "ok"
	dbMessage := "Database connection is healthy"

	if err := database.HealthCheck(ctx); err != nil {
		dbStatus = "error"
		dbMessage = err.Error()
	}

	return c.JSON(fiber.Map{
		"status":  "ok",
		"message": "Server is running",
		"database": fiber.Map{
			"status":  dbStatus,
			"message": dbMessage,
		},
	})
}
