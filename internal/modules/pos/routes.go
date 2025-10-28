package pos

import (
	"github.com/gofiber/fiber/v2"

	"inventory-modular-monolith/internal/database"
	inventoryRepo "inventory-modular-monolith/internal/modules/inventory/repository"
	"inventory-modular-monolith/internal/modules/pos/handler"
	"inventory-modular-monolith/internal/modules/pos/repository"
	"inventory-modular-monolith/internal/modules/pos/service"
)

func RegisterRoutes(router fiber.Router, db *database.MongoDB) {
	repo := repository.NewSaleRepository(db.Database)
	productRepo := inventoryRepo.NewProductRepository(db.Database)
	svc := service.NewSaleService(repo, productRepo)
	h := handler.NewSaleHandler(svc)

	pos := router.Group("/pos")
	pos.Post("/sales", h.CreateSale)
	pos.Get("/stores/:storeId/sales", h.GetSalesByStore)
}