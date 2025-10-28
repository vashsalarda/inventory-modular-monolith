package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"inventory-modular-monolith/internal/modules/merchant/domain"
)

type StoreRepository struct {
	collection *mongo.Collection
}

func NewStoreRepository(db *mongo.Database) *StoreRepository {
	return &StoreRepository{
		collection: db.Collection("stores"),
	}
}

func (r *StoreRepository) Create(ctx context.Context, store *domain.Store) error {
	store.CreatedAt = time.Now()
	store.UpdatedAt = time.Now()
	result, err := r.collection.InsertOne(ctx, store)
	if err != nil {
		return err
	}
	store.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *StoreRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*domain.Store, error) {
	var store domain.Store
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&store)
	if err != nil {
		return nil, err
	}
	return &store, nil
}

func (r *StoreRepository) FindAll(ctx context.Context) ([]domain.Store, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var stores []domain.Store
	if err := cursor.All(ctx, &stores); err != nil {
		return nil, err
	}
	return stores, nil
}