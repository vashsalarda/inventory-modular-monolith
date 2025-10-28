package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Store struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Email     string             `json:"email" bson:"email"`
	Phone     string             `json:"phone" bson:"phone"`
	Address   Address            `json:"address" bson:"address"`
	Status    string             `json:"status" bson:"status"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type Address struct {
	Street  string `json:"street" bson:"street"`
	City    string `json:"city" bson:"city"`
	State   string `json:"state" bson:"state"`
	ZipCode string `json:"zip_code" bson:"zip_code"`
	Country string `json:"country" bson:"country"`
}

type CreateStoreDTO struct {
	Name    string  `json:"name" validate:"required"`
	Email   string  `json:"email" validate:"required,email"`
	Phone   string  `json:"phone" validate:"required"`
	Address Address `json:"address" validate:"required"`
}