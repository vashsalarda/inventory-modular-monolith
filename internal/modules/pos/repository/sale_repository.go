package repository

import (
	"context"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

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

func (r *SaleRepository) FindByStoreID(ctx context.Context, storeID primitive.ObjectID, keyword string , page int64, page_size int64) (*domain.SalePage, error) {
	filter := bson.M{
		"store_id": storeID,
		"deleted": bson.M{
			"$ne": true,
		},
	}

	if keyword != "" {
		filter["$or"] = bson.A{
			bson.M{"name": bson.M{"$regex": keyword, "$options": "i"}},
			bson.M{"email": bson.M{"$regex": keyword, "$options": "i"}},
		}
	}

	docs := make([]domain.Sale, 0, page_size)  
	

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

	cursor, err := r.collection.Find(context.TODO(), filter, findOptions,)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}

	resp := &domain.SalePage{
		TotalRows:  totalItems,
		TotalPages: totalPages,
		Total:      totalItems,
		PageNumber: page,
		PageSize:   page_size,
		Data:       docs,
	}

	return resp, nil
}