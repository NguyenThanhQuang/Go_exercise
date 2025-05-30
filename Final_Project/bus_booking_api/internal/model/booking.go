package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookingStatus string

const (
	StatusPending   BookingStatus = "pending"
	StatusHeld      BookingStatus = "held"
	StatusConfirmed BookingStatus = "confirmed"
	StatusCancelled BookingStatus = "cancelled"
	StatusExpired   BookingStatus = "expired"
)

type PaymentStatus string

const (
	PaymentPending PaymentStatus = "pending"
	PaymentPaid    PaymentStatus = "paid"
	PaymentFailed  PaymentStatus = "failed"
)

type Booking struct {
	ID                        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID                    primitive.ObjectID `json:"userId" bson:"userId" validate:"required"`
	TripID                    primitive.ObjectID `json:"tripId" bson:"tripId" validate:"required"`
	BookingTime               time.Time          `json:"bookingTime" bson:"bookingTime"`
	Status                    BookingStatus      `json:"status" bson:"status" default:"pending"`
	HeldUntil                 *time.Time         `json:"heldUntil,omitempty" bson:"heldUntil,omitempty"`
	PaymentStatus             PaymentStatus      `json:"paymentStatus" bson:"paymentStatus" default:"pending"`
	PaymentMethod             string             `json:"paymentMethod,omitempty" bson:"paymentMethod,omitempty"`
	TotalAmount               float64            `json:"totalAmount" bson:"totalAmount" validate:"required,gt=0"`
	Passengers                []Passenger        `json:"passengers" bson:"passengers" validate:"required,dive"`
	TicketCode                string             `json:"ticketCode,omitempty" bson:"ticketCode,omitempty,unique"`
	PaymentGatewayTransactionID string             `json:"paymentGatewayTransactionId,omitempty" bson:"paymentGatewayTransactionId,omitempty"`
	CreatedAt                 time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt                 time.Time          `json:"updatedAt" bson:"updatedAt"`
}

type HoldSeatsRequest struct {
	TripID      string   `json:"tripId" binding:"required"`
	SeatNumbers []string `json:"seatNumbers" binding:"required,min=1"`
}
