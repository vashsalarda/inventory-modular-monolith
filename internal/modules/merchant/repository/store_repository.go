package repository

import (
	"context"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

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

func (r *StoreRepository) FindAll(ctx context.Context, keyword string , page int64, page_size int64) (domain.StorePage, error) {
	filter := bson.M{
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

	docs := make([]domain.Store, 0, page_size)  
	

	totalItems, err := r.collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return domain.StorePage{}, err
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
		return domain.StorePage{}, err
	}
	defer cursor.Close(ctx)
	
	if err := cursor.All(ctx, &docs); err != nil {
		return domain.StorePage{}, err
	}

	result := domain.StorePage{
		TotalRows:  totalItems,
		TotalPages: totalPages,
		Total:      totalItems,
		PageNumber: page,
		PageSize:   page_size,
		Data:       docs,
	}

	return result, nil
}