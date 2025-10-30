package inventory

import (
	"github.com/gofiber/fiber/v2"

	"inventory-modular-monolith/internal/database"
	"inventory-modular-monolith/internal/modules/inventory/handler"
	"inventory-modular-monolith/internal/modules/inventory/repository"
	"inventory-modular-monolith/internal/modules/inventory/service"
	store_repo "inventory-modular-monolith/internal/modules/merchant/repository"
	store_service "inventory-modular-monolith/internal/modules/merchant/service"
)

func RegisterRoutes(router fiber.Router, db *database.MongoDB) {
	repo := repository.NewProductRepository(db.Database)
	svc := service.NewProductService(repo)
	store_repo := store_repo.NewStoreRepository(db.Database)
	store_svc := store_service.NewStoreService(store_repo)
	h := handler.NewProductHandler(svc, store_svc)

	inventory := router.Group("/inventory")
	inventory.Get("/products", h.GetAllProducts)
	inventory.Post("/products", h.CreateProduct)
	inventory.Get("/products/:id", h.GetProduct)
	inventory.Get("/stores/:storeId/products", h.GetProductsByStore)
}