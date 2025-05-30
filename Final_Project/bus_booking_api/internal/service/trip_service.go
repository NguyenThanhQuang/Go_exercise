package service

import (
	"bus_booking_api/internal/model"
	"bus_booking_api/internal/repository"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TripService interface {
	SearchTrips(ctx context.Context, from, to, date string) ([]model.Trip, error)
	GetTripDetails(ctx context.Context, tripID string) (*model.Trip, error)
	CreateTrip(ctx context.Context, tripData *model.Trip) (*model.Trip, error) // For Admin
}

type tripServiceImpl struct {
	tripRepo    repository.TripRepository
	// companyRepo repository.CompanyRepository // (Cần sau này để validate)
	// vehicleRepo repository.VehicleRepository // (Cần sau này để validate)
}

func NewTripService(tripRepo repository.TripRepository) TripService {
	return &tripServiceImpl{
		tripRepo: tripRepo,
	}
}

func (s *tripServiceImpl) SearchTrips(ctx context.Context, from, to, date string) ([]model.Trip, error) {
	if from == "" || to == "" || date == "" {
		return nil, errors.New("from, to, and date parameters are required")
	}
	// Thêm validate cho date format nếu cần
	return s.tripRepo.SearchTrips(ctx, from, to, date)
}

func (s *tripServiceImpl) GetTripDetails(ctx context.Context, tripIDStr string) (*model.Trip, error) {
	tripID, err := primitive.ObjectIDFromHex(tripIDStr)
	if err != nil {
		return nil, errors.New("invalid trip ID format")
	}
	return s.tripRepo.GetTripByID(ctx, tripID)
}

// CreateTrip - dành cho Admin/Nhà xe
func (s *tripServiceImpl) CreateTrip(ctx context.Context, tripData *model.Trip) (*model.Trip, error) {
	// TODO: Validate companyId and vehicleId exist and are valid
	// Ví dụ:
	// company, err := s.companyRepo.GetCompanyByID(ctx, tripData.CompanyID)
	// if err != nil || company == nil { return nil, errors.New("invalid companyId")}
	// vehicle, err := s.vehicleRepo.GetVehicleByID(ctx, tripData.VehicleID)
	// if err != nil || vehicle == nil { return nil, errors.New("invalid vehicleId")}
    // if vehicle.CompanyID != tripData.CompanyID { return nil, errors.New("vehicle does not belong to the specified company")}


	// Logic khởi tạo seats đã được chuyển vào TripRepository.CreateTrip
	// Service có thể thực hiện các validation khác nếu cần trước khi gọi repo.

	if tripData.Price <= 0 {
		return nil, errors.New("price must be positive")
	}
	if tripData.DepartureTime.IsZero() || tripData.ExpectedArrivalTime.IsZero() {
		return nil, errors.New("departure and arrival times are required")
	}
	if tripData.DepartureTime.After(tripData.ExpectedArrivalTime) {
		return nil, errors.New("departure time must be before expected arrival time")
	}
    if tripData.Route.From.Name == "" || tripData.Route.To.Name == "" {
        return nil, errors.New("route from and to names are required")
    }
    // Thêm các validate khác cho route, stops...

	return s.tripRepo.CreateTrip(ctx, tripData)
}