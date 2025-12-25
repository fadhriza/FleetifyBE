package main

import (
	"log"

	"fleetify/internal/config"
	"fleetify/internal/database"
	"fleetify/internal/middleware"
	"fleetify/internal/routes"
	"fleetify/pkg/errors"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Load configuration
	if err := config.LoadConfig(); err != nil {
		errors.LogError("Failed to load configuration", err)
		log.Fatal("Failed to load configuration:", err)
	}

	// Connect to database
	if err := database.Connect(); err != nil {
		errors.LogError("Failed to connect to database", err)
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "Fleetify API",
		ServerHeader: "Fleetify",
		ErrorHandler: customErrorHandler,
	})

	// Middleware
	app.Use(middleware.Recover())
	app.Use(middleware.Logger())
	app.Use(middleware.CORS())

	// Setup routes
	routes.SetupRoutes(app)

	// Start server
	addr := config.AppConfig.Server.Host + ":" + config.AppConfig.Server.Port
	log.Printf("🚀 Server starting on %s", addr)
	log.Fatal(app.Listen(addr))
}

func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	// Log the error with full stack trace
	errors.LogError("HTTP Error", err)

	return c.Status(code).JSON(fiber.Map{
		"error":   true,
		"message": err.Error(),
		"code":    code,
	})
}
