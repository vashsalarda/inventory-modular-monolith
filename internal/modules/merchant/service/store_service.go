package service

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"inventory-modular-monolith/internal/modules/merchant/domain"
	"inventory-modular-monolith/internal/modules/merchant/repository"
)

type StoreService struct {
	repo *repository.StoreRepository
}

func NewStoreService(repo *repository.StoreRepository) *StoreService {
	return &StoreService{repo: repo}
}

func (s *StoreService) CreateStore(ctx context.Context, dto *domain.CreateStoreDTO) (*domain.Store, error) {
	store := &domain.Store{
		Name:    dto.Name,
		Email:   dto.Email,
		Phone:   dto.Phone,
		Address: dto.Address,
		Status:  "active",
	}

	if err := s.repo.Create(ctx, store); err != nil {
		return nil, err
	}
	return store, nil
}

func (s *StoreService) GetStore(ctx context.Context, id string) (*domain.Store, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid store ID")
	}
	return s.repo.FindByID(ctx, objID)
}

func (s *StoreService) GetAllStores(ctx context.Context) ([]domain.Store, error) {
	return s.repo.FindAll(ctx)
}