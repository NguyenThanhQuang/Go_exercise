package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TripStatus string

const (
	StatusScheduled TripStatus = "scheduled"
	StatusDeparted  TripStatus = "departed"
	StatusArrived   TripStatus = "arrived"
	// StatusCancelled TripStatus = "cancelled"
)

type Route struct {
	From     Location `json:"from" bson:"from"`
	To       Location `json:"to" bson:"to"`
	Stops    []Stop   `json:"stops,omitempty" bson:"stops,omitempty"`
	Polyline string   `json:"polyline,omitempty" bson:"polyline,omitempty"`
}

type Trip struct {
	ID                  primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	CompanyID           primitive.ObjectID `json:"companyId" bson:"companyId" validate:"required"`
	VehicleID           primitive.ObjectID `json:"vehicleId" bson:"vehicleId" validate:"required"`
	Route               Route              `json:"route" bson:"route"`
	DepartureTime       time.Time          `json:"departureTime" bson:"departureTime" validate:"required"`
	ExpectedArrivalTime time.Time          `json:"expectedArrivalTime" bson:"expectedArrivalTime" validate:"required"`
	Price               float64            `json:"price" bson:"price" validate:"required,gt=0"`
	Status              TripStatus         `json:"status" bson:"status" default:"scheduled"`
	Seats               []Seat             `json:"seats" bson:"seats"`
	CreatedAt           time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt           time.Time          `json:"updatedAt" bson:"updatedAt"`
}

type TripSearchQuery struct {
	From string `form:"from" binding:"required"`
	To   string `form:"to" binding:"required"`
	Date string `form:"date" binding:"required"`
}

