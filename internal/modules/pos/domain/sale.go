package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SaleItem struct {
	ProductID primitive.ObjectID `json:"product_id" bson:"product_id"`
	SKU       string             `json:"sku" bson:"sku"`
	Name      string             `json:"name" bson:"name"`
	Quantity  int                `json:"quantity" bson:"quantity"`
	Price     float64            `json:"price" bson:"price"`
	Subtotal  float64            `json:"subtotal" bson:"subtotal"`
}

type Sale struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	StoreID     primitive.ObjectID `json:"store_id" bson:"store_id"`
	Items       []SaleItem         `json:"items" bson:"items"`
	Total       float64            `json:"total" bson:"total"`
	Tax         float64            `json:"tax" bson:"tax"`
	GrandTotal  float64            `json:"grand_total" bson:"grand_total"`
	PaymentType string             `json:"payment_type" bson:"payment_type"`
	Status      string             `json:"status" bson:"status"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
}

type CreateSaleDTO struct {
	StoreID     string        `json:"store_id" validate:"required"`
	Items       []SaleItemDTO `json:"items" validate:"required,min=1"`
	PaymentType string        `json:"payment_type" validate:"required"`
}

type SaleItemDTO struct {
	ProductID string `json:"product_id" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required,gt=0"`
}

type SalePage struct {
	PageSize   int64  `json:"page_size"`
	PageNumber int64  `json:"page_number"`
	TotalRows  int64  `json:"total_rows"`
	Total      int64  `json:"total"`
	TotalPages int64  `json:"total_pages"`
	Data       []Sale `json:"data"`
}
