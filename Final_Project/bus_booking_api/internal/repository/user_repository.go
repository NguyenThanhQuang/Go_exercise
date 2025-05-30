package repository

import (
	"bus_booking_api/internal/config"
	"bus_booking_api/internal/model"
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const userCollectionName = "users"

type UserRepository interface {
	CreateUser(user *model.User) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	GetUserByPhone(phone string) (*model.User, error)
	GetUserByID(id primitive.ObjectID) (*model.User, error)
}

type userRepositoryImpl struct {
	collection *mongo.Collection
}

func NewUserRepository() UserRepository {
	return &userRepositoryImpl{
		collection: config.GetCollection(userCollectionName),
	}
}

func (r *userRepositoryImpl) CreateUser(user *model.User) (*model.User, error) {
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(context.Background(), user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, errors.New("email or phone already exists")
		}
		return nil, err
	}
	return user, nil
}

func (r *userRepositoryImpl) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepositoryImpl) GetUserByPhone(phone string) (*model.User, error) {
	var user model.User
	err := r.collection.FindOne(context.Background(), bson.M{"phone": phone}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepositoryImpl) GetUserByID(id primitive.ObjectID) (*model.User, error) {
	var user model.User
	err := r.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func EnsureUserIndexes() error {
	userRepo := NewUserRepository().(*userRepositoryImpl)

	emailIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	phoneIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "phone", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err := userRepo.collection.Indexes().CreateMany(context.Background(), []mongo.IndexModel{emailIndex, phoneIndex})
	if err != nil {
		return fmt.Errorf("failed to create indexes for users collection: %w", err)
	}
	fmt.Println("Successfully created indexes for users collection.")
	return nil
}
