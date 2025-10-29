package handler

import (
	"strconv"

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
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	sale, err := h.service.CreateSale(c.Context(), &dto)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(sale)
}

func (h *SaleHandler) GetSalesByStore(c *fiber.Ctx) error {
	storeID := c.Params("storeId")
	var page, page_size int64 = 1, 25
	if p, err := strconv.ParseInt(c.Query("page_number", "1"), 10, 64); err == nil {
		page = p
	}
	if ps, err := strconv.ParseInt(c.Query("page_size", "10"), 10, 64); err == nil {
		page_size = ps
	}
	keyword := c.Query("keyword", "")
	sales, err := h.service.GetSalesByStore(c.Context(), storeID, keyword, page, page_size)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(sales)
}