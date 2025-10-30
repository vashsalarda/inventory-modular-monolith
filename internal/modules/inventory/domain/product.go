package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	StoreID     primitive.ObjectID `json:"store_id" bson:"store_id"`
	SKU         string             `json:"sku" bson:"sku"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Price       float64            `json:"price" bson:"price"`
	Cost        float64            `json:"cost" bson:"cost"`
	Quantity    int                `json:"quantity" bson:"quantity"`
	Category    string             `json:"category" bson:"category"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

type CreateProductDTO struct {
	StoreID     string  `json:"store_id" validate:"required"`
	SKU         string  `json:"sku" validate:"required"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Cost        float64 `json:"cost" validate:"required,gte=0"`
	Quantity    int     `json:"quantity" validate:"required,gte=0"`
	Category    string  `json:"category"`
}

type ProductPage struct {
	PageSize   int64     `json:"page_size"`
	PageNumber int64     `json:"page_number"`
	TotalRows  int64     `json:"total_rows"`
	Total      int64     `json:"total"`
	TotalPages int64     `json:"total_pages"`
	Data       []Product `json:"data"`
}
