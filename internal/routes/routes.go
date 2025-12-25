package routes

import (
	"github.com/gofiber/fiber/v2"
	"fleetify/internal/handlers"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v1")

	api.Get("/health", handlers.HealthCheck)
	api.Get("/health/db", handlers.HealthCheckDB)
}

