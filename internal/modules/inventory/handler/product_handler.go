package handler

import (
	"github.com/gofiber/fiber/v2"

	"inventory-modular-monolith/internal/modules/inventory/domain"
	"inventory-modular-monolith/internal/modules/inventory/service"
)

type ProductHandler struct {
	service *service.ProductService
}

func NewProductHandler(service *service.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}


func (h *ProductHandler) GetAllProducts(c *fiber.Ctx) error {				
	products, err := h.service.GetAllProducts(c.Context())
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(products)
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
	products, err := h.service.GetProductsByStore(c.Context(), storeID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(products)
}