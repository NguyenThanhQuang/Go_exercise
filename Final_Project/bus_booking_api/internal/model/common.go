package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Location struct {
	Name    string  `json:"name" bson:"name"`
	GeoJSON GeoJSON `json:"location" bson:"location"`
}

type GeoJSON struct {
	Type        string    `json:"type" bson:"type" default:"Point"`
	Coordinates []float64 `json:"coordinates" bson:"coordinates"`
}

type Stop struct {
	Name                  string    `json:"name" bson:"name"`
	Location              GeoJSON   `json:"location" bson:"location"`
	ExpectedArrivalTime   time.Time `json:"expectedArrivalTime,omitempty" bson:"expectedArrivalTime,omitempty"`
	ExpectedDepartureTime time.Time `json:"expectedDepartureTime,omitempty" bson:"expectedDepartureTime,omitempty"`
}

type Seat struct {
	SeatNumber string             `json:"seatNumber" bson:"seatNumber"`
	Status     string             `json:"status" bson:"status" default:"available"`
	BookingID  primitive.ObjectID `json:"bookingId,omitempty" bson:"bookingId,omitempty"`
}

type Passenger struct {
	Name       string `json:"name" bson:"name" validate:"required"`
	Phone      string `json:"phone" bson:"phone" validate:"required"`
	SeatNumber string `json:"seatNumber" bson:"seatNumber" validate:"required"`
}

type SeatMapDefinition struct {
	Rows       int        `json:"rows" bson:"rows"`
	Cols       int        `json:"cols" bson:"cols"`
	Layout     [][]string `json:"layout" bson:"layout"`
	TotalSeats int        `json:"totalSeats" bson:"totalSeats"`
}
