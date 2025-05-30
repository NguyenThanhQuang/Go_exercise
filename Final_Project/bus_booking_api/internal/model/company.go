package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Company struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name" validate:"required"`
	Code        string             `json:"code" bson:"code" validate:"required"`
	Address     string             `json:"address,omitempty" bson:"address,omitempty"`
	Phone       string             `json:"phone,omitempty" bson:"phone,omitempty"`
	Email       string             `json:"email,omitempty" bson:"email,omitempty"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
	LogoURL     string             `json:"logoUrl,omitempty" bson:"logoUrl,omitempty"`
	IsActive    bool               `json:"isActive" bson:"isActive" default:"true"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt"`
}
