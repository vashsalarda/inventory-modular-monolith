package repository

import (
	"context"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

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

func (r *ProductRepository) FindAll(ctx context.Context, keyword string, page int64, page_size int64) (*domain.ProductPage, error) {
	filter := bson.M{
		"deleted": bson.M{
			"$ne": true,
		},
	}

	if keyword != "" {
		filter["$or"] = bson.A{
			bson.M{"name": bson.M{"$regex": keyword, "$options": "i"}},
			bson.M{"sku": bson.M{"$regex": keyword, "$options": "i"}},
			bson.M{"description": bson.M{"$regex": keyword, "$options": "i"}},
		}
	}

	docs := make([]domain.Product, 0, page_size)

	totalItems, err := r.collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	// Calculate total pages
	totalPages := int64(math.Ceil(float64(totalItems) / float64(page_size)))

	// Calculate skip
	skip := (page - 1) * page_size

	// Find options with limit and skip
	findOptions := options.Find()
	findOptions.SetLimit(page_size)
	findOptions.SetSkip(skip)
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}

	resp := &domain.ProductPage{
		TotalRows:  totalItems,
		TotalPages: totalPages,
		Total:      totalItems,
		PageNumber: page,
		PageSize:   page_size,
		Data:       docs,
	}

	return resp, nil
}

func (r *ProductRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*domain.Product, error) {
	var product domain.Product
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&product)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) FindByStoreID(ctx context.Context, storeID primitive.ObjectID, keyword string, page int64, page_size int64) (*domain.ProductPage, error) {
	filter := bson.M{
		"store_id": storeID,
		"deleted": bson.M{
			"$ne": true,
		},
	}

	if keyword != "" {
		filter["$or"] = bson.A{
			bson.M{"name": bson.M{"$regex": keyword, "$options": "i"}},
			bson.M{"sku": bson.M{"$regex": keyword, "$options": "i"}},
			bson.M{"description": bson.M{"$regex": keyword, "$options": "i"}},
		}
	}

	docs := make([]domain.Product, 0, page_size)

	totalItems, err := r.collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	// Calculate total pages
	totalPages := int64(math.Ceil(float64(totalItems) / float64(page_size)))

	// Calculate skip
	skip := (page - 1) * page_size

	// Find options with limit and skip
	findOptions := options.Find()
	findOptions.SetLimit(page_size)
	findOptions.SetSkip(skip)
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}

	resp := &domain.ProductPage{
		TotalRows:  totalItems,
		TotalPages: totalPages,
		Total:      totalItems,
		PageNumber: page,
		PageSize:   page_size,
		Data:       docs,
	}

	return resp, nil
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
