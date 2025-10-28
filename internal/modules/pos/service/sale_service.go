package service

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"

	inventoryRepo "inventory-modular-monolith/internal/modules/inventory/repository"
	"inventory-modular-monolith/internal/modules/pos/domain"
	"inventory-modular-monolith/internal/modules/pos/repository"
)

type SaleService struct {
	repo        *repository.SaleRepository
	productRepo *inventoryRepo.ProductRepository
}

func NewSaleService(repo *repository.SaleRepository, productRepo *inventoryRepo.ProductRepository) *SaleService {
	return &SaleService{
		repo:        repo,
		productRepo: productRepo,
	}
}

func (s *SaleService) CreateSale(ctx context.Context, dto *domain.CreateSaleDTO) (*domain.Sale, error) {
	storeID, err := primitive.ObjectIDFromHex(dto.StoreID)
	if err != nil {
		return nil, errors.New("invalid store ID")
	}

	var items []domain.SaleItem
	var total float64

	for _, item := range dto.Items {
		productID, err := primitive.ObjectIDFromHex(item.ProductID)
		if err != nil {
			return nil, errors.New("invalid product ID")
		}

		product, err := s.productRepo.FindByID(ctx, productID)
		if err != nil {
			return nil, errors.New("product not found")
		}

		if product.Quantity < item.Quantity {
			return nil, errors.New("insufficient stock for product: " + product.Name)
		}

		subtotal := float64(item.Quantity) * product.Price
		items = append(items, domain.SaleItem{
			ProductID: productID,
			SKU:       product.SKU,
			Name:      product.Name,
			Quantity:  item.Quantity,
			Price:     product.Price,
			Subtotal:  subtotal,
		})
		total += subtotal

		// Deduct stock
		if err := s.productRepo.UpdateQuantity(ctx, productID, product.Quantity-item.Quantity); err != nil {
			return nil, err
		}
	}

	tax := total * 0.12 // 12% tax
	grandTotal := total + tax

	sale := &domain.Sale{
		StoreID:     storeID,
		Items:       items,
		Total:       total,
		Tax:         tax,
		GrandTotal:  grandTotal,
		PaymentType: dto.PaymentType,
		Status:      "completed",
	}

	if err := s.repo.Create(ctx, sale); err != nil {
		return nil, err
	}

	return sale, nil
}

func (s *SaleService) GetSalesByStore(ctx context.Context, storeID string) ([]domain.Sale, error) {
	objID, err := primitive.ObjectIDFromHex(storeID)
	if err != nil {
		return nil, errors.New("invalid store ID")
	}
	return s.repo.FindByStoreID(ctx, objID)
}