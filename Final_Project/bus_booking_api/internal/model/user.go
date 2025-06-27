package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRole string

const (
	RoleUser         UserRole = "user"
	RoleCompanyAdmin UserRole = "company_admin"
	RoleAdmin        UserRole = "admin"
)

type User struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Email        string             `json:"email" bson:"email" validate:"required,email"`
	Phone        string             `json:"phone" bson:"phone" validate:"required"`
	PasswordHash string             `json:"-" bson:"passwordHash"`
	Name         string             `json:"name" bson:"name" validate:"required"`
	Role         UserRole           `json:"role" bson:"role" default:"user"`
	CompanyID    primitive.ObjectID `json:"companyId,omitempty" bson:"companyId,omitempty"`
	CreatedAt    time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt    time.Time          `json:"updatedAt" bson:"updatedAt"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
}

type LoginRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}