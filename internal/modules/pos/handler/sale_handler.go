package handler

import (
	"github.com/gofiber/fiber/v2"

	"inventory-modular-monolith/internal/modules/pos/domain"
	"inventory-modular-monolith/internal/modules/pos/service"
)

type SaleHandler struct {
	service *service.SaleService
}

func NewSaleHandler(service *service.SaleService) *SaleHandler {
	return &SaleHandler{service: service}
}

func (h *SaleHandler) CreateSale(c *fiber.Ctx) error {
	var dto domain.CreateSaleDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	sale, err := h.service.CreateSale(c.Context(), &dto)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(sale)
}

func (h *SaleHandler) GetSalesByStore(c *fiber.Ctx) error {
	storeID := c.Params("storeId")
	sales, err := h.service.GetSalesByStore(c.Context(), storeID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(sales)
}