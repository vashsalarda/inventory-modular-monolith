package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"inventory-modular-monolith/internal/modules/inventory/domain"
)

type ProductRepository struct {
	collection *mongo.Collection
}

func NewProductRepository(db *mongo.Database) *ProductRepository {
	return &ProductRepository{
		collection: db.Collection("products"),
	}
}

func (r *ProductRepository) Create(ctx context.Context, product *domain.Product) error {
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()
	result, err := r.collection.InsertOne(ctx, product)
	if err != nil {
		return err
	}
	product.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *ProductRepository) FindAll(ctx context.Context) ([]domain.Product, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products = make([]domain.Product, 0, 100)
	if err := cursor.All(ctx, &products); err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*domain.Product, error) {
	var product domain.Product
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&product)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) FindByStoreID(ctx context.Context, storeID primitive.ObjectID) ([]domain.Product, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"store_id": storeID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []domain.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductRepository) Update(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	update["updated_at"] = time.Now()
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	return err
}

func (r *ProductRepository) UpdateQuantity(ctx context.Context, id primitive.ObjectID, quantity int) error {
	update := bson.M{
		"quantity":   quantity,
		"updated_at": time.Now(),
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	return err
}
