package merchant

import (
	"github.com/gofiber/fiber/v2"

	"inventory-modular-monolith/internal/database"
	"inventory-modular-monolith/internal/modules/merchant/handler"
	"inventory-modular-monolith/internal/modules/merchant/repository"
	"inventory-modular-monolith/internal/modules/merchant/service"
)

func RegisterRoutes(router fiber.Router, db *database.MongoDB) {
	repo := repository.NewStoreRepository(db.Database)
	svc := service.NewStoreService(repo)
	h := handler.NewStoreHandler(svc)

	merchant := router.Group("/merchants")
	merchant.Post("/stores", h.CreateStore)
	merchant.Get("/stores/:id", h.GetStore)
	merchant.Get("/stores", h.GetAllStores)
}