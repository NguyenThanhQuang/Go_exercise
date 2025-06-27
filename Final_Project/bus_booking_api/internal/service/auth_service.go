package service

import (
	"bus_booking_api/internal/model"
	"bus_booking_api/internal/repository"
	"bus_booking_api/internal/utils"
	"errors"
	"strings"
)

type AuthService interface {
	RegisterUser(req model.RegisterRequest) (*model.User, error)
	LoginUser(req model.LoginRequest) (*model.LoginResponse, error)
}

type authServiceImpl struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authServiceImpl{
		userRepo: userRepo,
	}
}

func (s *authServiceImpl) RegisterUser(req model.RegisterRequest) (*model.User, error) {
	existingUser, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, errors.New("error checking email: " + err.Error())
	}
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	existingUser, err = s.userRepo.GetUserByPhone(req.Phone)
	if err != nil {
		return nil, errors.New("error checking phone: " + err.Error())
	}
	if existingUser != nil {
		return nil, errors.New("phone number already registered")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("error hashing password: " + err.Error())
	}

	newUser := &model.User{
		Email:        strings.ToLower(req.Email),
		Phone:        req.Phone,
		PasswordHash: hashedPassword,
		Name:         req.Name,
		Role:         model.RoleUser,
	}

	createdUser, err := s.userRepo.CreateUser(newUser)
	if err != nil {
		return nil, errors.New("error creating user: " + err.Error())
	}

	createdUser.PasswordHash = ""
	return createdUser, nil
}

func (s *authServiceImpl) LoginUser(req model.LoginRequest) (*model.LoginResponse, error) {
	var user *model.User
	var err error

	if strings.Contains(req.Login, "@") {
		user, err = s.userRepo.GetUserByEmail(strings.ToLower(req.Login))
	} else {
		user, err = s.userRepo.GetUserByPhone(req.Login)
	}

	if err != nil {
		return nil, errors.New("error finding user: " + err.Error())
	}
	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	token, err := utils.GenerateToken(user)
	if err != nil {
		return nil, errors.New("error generating token: " + err.Error())
	}

	user.PasswordHash = ""

	return &model.LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}
