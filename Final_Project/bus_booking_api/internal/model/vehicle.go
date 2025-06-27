package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Vehicle struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	CompanyID   primitive.ObjectID `json:"companyId" bson:"companyId" validate:"required"`
	Type        string             `json:"type" bson:"type" validate:"required"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
	SeatMap     SeatMapDefinition  `json:"seatMap" bson:"seatMap"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt"`
}
