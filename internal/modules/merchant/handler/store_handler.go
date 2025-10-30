package handler

import (
	"strconv"

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
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	store, err := h.service.CreateStore(c.Context(), &dto)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(store)
}

func (h *StoreHandler) GetStore(c *fiber.Ctx) error {
	id := c.Params("id")
	store, err := h.service.GetStore(c.Context(), id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Store not found"})
	}
	return c.JSON(store)
}

func (h *StoreHandler) GetAllStores(c *fiber.Ctx) error {
	var page, page_size int64 = 1, 25
	if p, err := strconv.ParseInt(c.Query("page_number", "1"), 10, 64); err == nil {
		page = p
	}
	if ps, err := strconv.ParseInt(c.Query("page_size", "10"), 10, 64); err == nil {
		page_size = ps
	}
	keyword := c.Query("keyword", "")

	stores, err := h.service.GetAllStores(c.Context(), keyword, page, page_size)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(stores)
}