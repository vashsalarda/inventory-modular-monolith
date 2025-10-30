package pos

import (
	"github.com/gofiber/fiber/v2"

	"inventory-modular-monolith/internal/database"
	inventoryRepo "inventory-modular-monolith/internal/modules/inventory/repository"
	store_repo "inventory-modular-monolith/internal/modules/merchant/repository"
	store_service "inventory-modular-monolith/internal/modules/merchant/service"
	"inventory-modular-monolith/internal/modules/pos/handler"
	"inventory-modular-monolith/internal/modules/pos/repository"
	"inventory-modular-monolith/internal/modules/pos/service"
)

func RegisterRoutes(router fiber.Router, db *database.MongoDB) {
	repo := repository.NewSaleRepository(db.Database)
	productRepo := inventoryRepo.NewProductRepository(db.Database)
	svc := service.NewSaleService(repo, productRepo)
	store_repo := store_repo.NewStoreRepository(db.Database)
	store_svc := store_service.NewStoreService(store_repo)
	h := handler.NewSaleHandler(svc, store_svc)

	pos := router.Group("/pos")
	pos.Post("/sales", h.CreateSale)
	pos.Get("/stores/:storeId/sales", h.GetSalesByStore)
}