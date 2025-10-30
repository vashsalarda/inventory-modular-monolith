package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"inventory-modular-monolith/internal/modules/inventory/domain"
	"inventory-modular-monolith/internal/modules/inventory/service"
	store_service "inventory-modular-monolith/internal/modules/merchant/service"
)

type ProductHandler struct {
	service *service.ProductService
	store_service *store_service.StoreService
}

func NewProductHandler(service *service.ProductService, store_service *store_service.StoreService) *ProductHandler {
	return &ProductHandler{service: service, store_service: store_service}
}

func (h *ProductHandler) GetAllProducts(c *fiber.Ctx) error {
	var page, page_size int64 = 1, 25
	if p, err := strconv.ParseInt(c.Query("page_number", "1"), 10, 64); err == nil {
		page = p
	}
	if ps, err := strconv.ParseInt(c.Query("page_size", "10"), 10, 64); err == nil {
		page_size = ps
	}
	keyword := c.Query("keyword", "")
	res, err := h.service.GetAllProducts(c.Context(), keyword, page, page_size)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	var dto domain.CreateProductDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	product, err := h.service.CreateProduct(c.Context(), &dto)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(product)
}

func (h *ProductHandler) GetProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	product, err := h.service.GetProduct(c.Context(), id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Product not found"})
	}
	return c.JSON(product)
}

func (h *ProductHandler) GetProductsByStore(c *fiber.Ctx) error {
	storeID := c.Params("storeId")
	_, err := h.store_service.GetStore(c.Context(), storeID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Store not found"})
	}
	var page, page_size int64 = 1, 25
	if p, err := strconv.ParseInt(c.Query("page_number", "1"), 10, 64); err == nil {
		page = p
	}
	if ps, err := strconv.ParseInt(c.Query("page_size", "10"), 10, 64); err == nil {
		page_size = ps
	}
	keyword := c.Query("keyword", "")
	res, err := h.service.GetProductsByStore(c.Context(), storeID, keyword, page, page_size)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}
