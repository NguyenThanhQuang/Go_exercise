package service

import (
	"bus_booking_api/internal/model"
	"bus_booking_api/internal/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService interface {
	GetUserProfile(userID primitive.ObjectID) (*model.User, error)
}

type userServiceImpl struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userServiceImpl{userRepo: userRepo}
}

func (s *userServiceImpl) GetUserProfile(userID primitive.ObjectID) (*model.User, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	if user != nil {
		user.PasswordHash = "" 
	}
	return user, nil
}