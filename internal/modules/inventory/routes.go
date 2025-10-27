package inventory

import (
	"github.com/gofiber/fiber/v2"

	"inventory-modular-monolith/internal/database"
	"inventory-modular-monolith/internal/modules/inventory/handler"
	"inventory-modular-monolith/internal/modules/inventory/repository"
	"inventory-modular-monolith/internal/modules/inventory/service"
)

func RegisterRoutes(router fiber.Router, db *database.MongoDB) {
	repo := repository.NewProductRepository(db.Database)
	svc := service.NewProductService(repo)
	h := handler.NewProductHandler(svc)

	inventory := router.Group("/inventory")
	inventory.Get("/products", h.GetAllProducts)
	inventory.Post("/products", h.CreateProduct)
	inventory.Get("/products/:id", h.GetProduct)
	inventory.Get("/stores/:storeId/products", h.GetProductsByStore)
}