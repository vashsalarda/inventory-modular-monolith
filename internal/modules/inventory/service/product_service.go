package service

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"inventory-modular-monolith/internal/modules/inventory/domain"
	"inventory-modular-monolith/internal/modules/inventory/repository"
)

type ProductService struct {
	repo *repository.ProductRepository
}

func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) CreateProduct(ctx context.Context, dto *domain.CreateProductDTO) (*domain.Product, error) {
	storeID, err := primitive.ObjectIDFromHex(dto.StoreID)
	if err != nil {
		return nil, errors.New("invalid store ID")
	}

	product := &domain.Product{
		StoreID:     storeID,
		SKU:         dto.SKU,
		Name:        dto.Name,
		Description: dto.Description,
		Price:       dto.Price,
		Cost:        dto.Cost,
		Quantity:    dto.Quantity,
		Category:    dto.Category,
	}

	if err := s.repo.Create(ctx, product); err != nil {
		return nil, err
	}
	return product, nil
}

func (s *ProductService) GetAllProducts(ctx context.Context) ([]domain.Product, error) {
	return s.repo.FindAll(ctx)
}

func (s *ProductService) GetProduct(ctx context.Context, id string) (*domain.Product, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid product ID")
	}
	return s.repo.FindByID(ctx, objID)
}

func (s *ProductService) GetProductsByStore(ctx context.Context, storeID string) ([]domain.Product, error) {
	objID, err := primitive.ObjectIDFromHex(storeID)
	if err != nil {
		return nil, errors.New("invalid store ID")
	}
	return s.repo.FindByStoreID(ctx, objID)
}

func (s *ProductService) DeductStock(ctx context.Context, productID primitive.ObjectID, quantity int) error {
	product, err := s.repo.FindByID(ctx, productID)
	if err != nil {
		return err
	}
	
	if product.Quantity < quantity {
		return errors.New("insufficient stock")
	}
	
	return s.repo.UpdateQuantity(ctx, productID, product.Quantity-quantity)
}