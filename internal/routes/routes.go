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

	roles := api.Group("/roles", middleware.Auth(), middleware.Admin())
	roles.Get("/", handlers.GetRoles)
	roles.Get("/:oid", handlers.GetRoleByOID)
	roles.Post("/", handlers.CreateRole)
	roles.Put("/:oid", handlers.UpdateRole)
	roles.Delete("/:oid", handlers.DeleteRole)

	items := api.Group("/items", middleware.Auth())
	items.Get("/", handlers.GetItems)
	items.Get("/:id", handlers.GetItemById)
	items.Post("/", middleware.ItemModifyAccess(), handlers.CreateItem)
	items.Put("/:id", middleware.ItemModifyAccess(), handlers.UpdateItem)
	items.Delete("/:id", middleware.ItemModifyAccess(), handlers.DeleteItem)

	suppliers := api.Group("/suppliers", middleware.Auth())
	suppliers.Get("/", handlers.GetSuppliers)
	suppliers.Get("/:id", handlers.GetSupplierById)
	suppliers.Post("/", handlers.CreateSupplier)
	suppliers.Put("/:id", handlers.UpdateSupplier)
	suppliers.Delete("/:id", handlers.DeleteSupplier)

	purchasings := api.Group("/purchasings", middleware.Auth())
	purchasings.Get("/", handlers.GetPurchasings)
	purchasings.Get("/:id", handlers.GetPurchasingById)
	purchasings.Post("/", handlers.CreatePurchasing)
	purchasings.Put("/:id", handlers.UpdatePurchasing)
	purchasings.Delete("/:id", handlers.DeletePurchasing)

	purchasingDetails := api.Group("/purchasing-details", middleware.Auth())
	purchasingDetails.Get("/", handlers.GetPurchasingDetails)
	purchasingDetails.Get("/purchasing/:purchasing_id", handlers.GetPurchasingDetailsByPurchasingId)
	purchasingDetails.Get("/:id", handlers.GetPurchasingDetailById)
	purchasingDetails.Post("/", handlers.CreatePurchasingDetail)
	purchasingDetails.Put("/:id", handlers.UpdatePurchasingDetail)
	purchasingDetails.Delete("/:id", handlers.DeletePurchasingDetail)
}
