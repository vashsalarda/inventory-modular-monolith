package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"inventory-modular-monolith/internal/modules/pos/domain"
)

type SaleRepository struct {
	collection *mongo.Collection
}

func NewSaleRepository(db *mongo.Database) *SaleRepository {
	return &SaleRepository{
		collection: db.Collection("sales"),
	}
}

func (r *SaleRepository) Create(ctx context.Context, sale *domain.Sale) error {
	sale.CreatedAt = time.Now()
	result, err := r.collection.InsertOne(ctx, sale)
	if err != nil {
		return err
	}
	sale.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *SaleRepository) FindByStoreID(ctx context.Context, storeID primitive.ObjectID) ([]domain.Sale, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"store_id": storeID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var sales []domain.Sale
	if err := cursor.All(ctx, &sales); err != nil {
		return nil, err
	}
	return sales, nil
}