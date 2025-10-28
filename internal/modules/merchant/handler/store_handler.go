package handler

import (
	"github.com/gofiber/fiber/v2"

	"inventory-modular-monolith/internal/modules/merchant/domain"
	"inventory-modular-monolith/internal/modules/merchant/service"
)

type StoreHandler struct {
	service *service.StoreService
}

func NewStoreHandler(service *service.StoreService) *StoreHandler {
	return &StoreHandler{service: service}
}

func (h *StoreHandler) CreateStore(c *fiber.Ctx) error {
	var dto domain.CreateStoreDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	store, err := h.service.CreateStore(c.Context(), &dto)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(store)
}

func (h *StoreHandler) GetStore(c *fiber.Ctx) error {
	id := c.Params("id")
	store, err := h.service.GetStore(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Store not found"})
	}
	return c.JSON(store)
}

func (h *StoreHandler) GetAllStores(c *fiber.Ctx) error {
	stores, err := h.service.GetAllStores(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(stores)
}