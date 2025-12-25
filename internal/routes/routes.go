package routes

import (
	"fleetify/internal/handlers"
	"fleetify/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v1")

	api.Get("/health", handlers.HealthCheck)
	api.Get("/health/db", handlers.HealthCheckDB)

	auth := api.Group("/auth")
	auth.Post("/register", handlers.Register)
	auth.Post("/login", handlers.Login)
	auth.Post("/token", middleware.Auth(), handlers.GetToken)
	auth.Post("/logout", middleware.Auth(), handlers.Logout)

	users := api.Group("/users", middleware.Auth())
	users.Get("/", handlers.GetUsers)

	user := api.Group("/user", middleware.Auth())
	user.Get("/:uname", handlers.GetUserByUsername)

	userAdmin := api.Group("/user", middleware.Auth(), middleware.Admin())
	userAdmin.Put("/:uname", handlers.UpdateUserByUsername)
	userAdmin.Put("/:uname/password", handlers.ChangePassword)
	userAdmin.Delete("/:uname", handlers.DeleteUser)
}
