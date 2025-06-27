package repository

import (
	"bus_booking_api/internal/config"
	"bus_booking_api/internal/model"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const tripCollectionName = "trips"

type TripRepository interface {
	CreateTrip(ctx context.Context, trip *model.Trip) (*model.Trip, error)
	GetTripByID(ctx context.Context, tripID primitive.ObjectID) (*model.Trip, error)
	SearchTrips(ctx context.Context, from, to, date string) ([]model.Trip, error)
	UpdateTripSeats(ctx context.Context, tripID primitive.ObjectID, seatNumbers []string, newStatus string, bookingID *primitive.ObjectID) error
	// GetTripByIDForUpdate (dùng với transaction sau này nếu cần)
}

type tripRepositoryImpl struct {
	collection *mongo.Collection
	// Thêm các collection khác nếu cần populate, ví dụ vehicleCollection
	vehicleCollection *mongo.Collection
	companyCollection *mongo.Collection
}

func NewTripRepository() TripRepository {
	return &tripRepositoryImpl{
		collection:        config.GetCollection(tripCollectionName),
		vehicleCollection: config.GetCollection(vehicleCollectionName), // Giả sử vehicleCollectionName = "vehicles"
		companyCollection: config.GetCollection(companyCollectionName), // Giả sử companyCollectionName = "companies"
	}
}

func (r *tripRepositoryImpl) CreateTrip(ctx context.Context, trip *model.Trip) (*model.Trip, error) {
	trip.ID = primitive.NewObjectID()
	trip.CreatedAt = time.Now()
	trip.UpdatedAt = time.Now()

	// Lấy thông tin vehicle để khởi tạo seats
	var vehicle model.Vehicle
	err := r.vehicleCollection.FindOne(ctx, bson.M{"_id": trip.VehicleID, "companyId": trip.CompanyID}).Decode(&vehicle)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("vehicle not found or does not belong to the company")
		}
		return nil, fmt.Errorf("error finding vehicle: %w", err)
	}

	// Khởi tạo seats dựa trên seatMap của vehicle
	if vehicle.SeatMap.Layout != nil && len(vehicle.SeatMap.Layout) > 0 {
		trip.Seats = []model.Seat{} // Khởi tạo slice rỗng
		for _, row := range vehicle.SeatMap.Layout {
			for _, seatNum := range row {
				if seatNum != "" && seatNum != "null" { // Bỏ qua lối đi hoặc các giá trị null/empty
					trip.Seats = append(trip.Seats, model.Seat{
						SeatNumber: seatNum,
						Status:     "available", // Mặc định
					})
				}
			}
		}
	} else {
        // Hoặc dựa vào TotalSeats nếu không có layout chi tiết
        if vehicle.SeatMap.TotalSeats == 0 {
            return nil, errors.New("vehicle has no seat map layout or total seats defined")
        }
        log.Printf("Warning: Vehicle %s has no detailed seat map layout. Initializing %d seats with generic numbers.", vehicle.ID.Hex(), vehicle.SeatMap.TotalSeats)
        trip.Seats = make([]model.Seat, 0, vehicle.SeatMap.TotalSeats)
        for i := 1; i <= vehicle.SeatMap.TotalSeats; i++ {
             trip.Seats = append(trip.Seats, model.Seat{
                SeatNumber: fmt.Sprintf("S%d", i), // Generic seat number
                Status:     "available",
            })
        }
	}


	_, err = r.collection.InsertOne(ctx, trip)
	if err != nil {
		return nil, fmt.Errorf("error creating trip: %w", err)
	}
	return trip, nil
}

func (r *tripRepositoryImpl) GetTripByID(ctx context.Context, tripID primitive.ObjectID) (*model.Trip, error) {
	var trip model.Trip
	err := r.collection.FindOne(ctx, bson.M{"_id": tripID}).Decode(&trip)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("trip not found")
		}
		return nil, err
	}

	// (Tùy chọn) Populate company và vehicle
	// Đây là cách populate đơn giản, không dùng $lookup của MongoDB aggregation
	var company model.Company
	err = r.companyCollection.FindOne(ctx, bson.M{"_id": trip.CompanyID}).Decode(&company)
	if err != nil {
		log.Printf("Warning: could not populate company for trip %s: %v", tripID.Hex(), err)
	} else {
		// Gán thông tin cơ bản, không gán toàn bộ object để tránh vòng lặp hoặc quá nhiều data
		// Hoặc tạo một struct con trong Trip để chứa thông tin populate
		// Ví dụ: trip.PopulatedCompany = company (cần định nghĩa field này trong model.Trip)
	}

	var vehicle model.Vehicle
	err = r.vehicleCollection.FindOne(ctx, bson.M{"_id": trip.VehicleID}).Decode(&vehicle)
	if err != nil {
		log.Printf("Warning: could not populate vehicle for trip %s: %v", tripID.Hex(), err)
	} else {
		// Tương tự company
	}


	return &trip, nil
}

func (r *tripRepositoryImpl) SearchTrips(ctx context.Context, from, to, dateStr string) ([]model.Trip, error) {
	// Chuyển đổi date string (YYYY-MM-DD) thành time.Time cho MongoDB query
	layout := "2006-01-02"
	parsedDate, err := time.Parse(layout, dateStr)
	if err != nil {
		return nil, errors.New("invalid date format, please use YYYY-MM-DD")
	}

	// Xác định khoảng thời gian cho ngày đó (từ 00:00:00 đến 23:59:59)
	startOfDay := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 0, 0, 0, 0, parsedDate.Location())
	endOfDay := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 23, 59, 59, 999999999, parsedDate.Location())

	filter := bson.M{
		"route.from.name": bson.M{"$regex": primitive.Regex{Pattern: from, Options: "i"}}, // Tìm kiếm không phân biệt hoa thường
		"route.to.name":   bson.M{"$regex": primitive.Regex{Pattern: to, Options: "i"}},
		"departureTime": bson.M{
			"$gte": startOfDay,
			"$lte": endOfDay,
		},
		"status": model.StatusScheduled, // Chỉ tìm các chuyến 'scheduled'
	}

	cursor, err := r.collection.Find(ctx, filter, options.Find().SetSort(bson.D{{"departureTime", 1}})) // Sắp xếp theo giờ đi
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var trips []model.Trip
	if err = cursor.All(ctx, &trips); err != nil {
		return nil, err
	}

    // Tính số ghế còn trống cho mỗi chuyến (đơn giản hóa, không populate company/vehicle ở đây)
    for i := range trips {
        availableSeats := 0
        for _, seat := range trips[i].Seats {
            if seat.Status == "available" {
                availableSeats++
            }
        }
        // Bạn có thể thêm một field `AvailableSeatsCount` vào model.Trip (không lưu vào DB) để truyền cho frontend
        // Hoặc frontend tự tính dựa trên mảng seats.
        // Để đơn giản, ta không thêm field mới vào struct Trip.
        // Thông tin này có thể được tính toán và thêm vào một DTO (Data Transfer Object) nếu cần.
		log.Printf("Trip %s has %d available seats", trips[i].ID.Hex(), availableSeats)
    }


	return trips, nil
}

// UpdateTripSeats cập nhật trạng thái của các ghế được chỉ định trong một chuyến đi.
// bookingID là optional, chỉ truyền khi chuyển sang 'held' hoặc 'booked'.
func (r *tripRepositoryImpl) UpdateTripSeats(ctx context.Context, tripID primitive.ObjectID, seatNumbers []string, newStatus string, bookingID *primitive.ObjectID) error {
	if len(seatNumbers) == 0 {
		return errors.New("no seat numbers provided for update")
	}

	// Nếu chuyển sang 'held' hoặc 'booked', gán bookingId
	// Nếu chuyển về 'available', xóa bookingId (gán nil)
	if newStatus == string(model.StatusHeld) || newStatus == string(model.StatusConfirmed) { // Giả sử 'confirmed' cho booking nghĩa là ghế 'booked'
		if bookingID == nil || bookingID.IsZero() {
			return errors.New("bookingID is required when setting seat status to held or booked")
		}
	} else if newStatus == "available" {
		// No action needed here, handled below
	} else {
		return fmt.Errorf("unsupported new seat status: %s", newStatus)
	}


	// Cập nhật tất cả các ghế khớp với seatNumbers trong một lần.
    // Dùng $[] để cập nhật tất cả các phần tử trong mảng seats khớp với arrayFilters
	arrayFilters := options.UpdateOptions{
		ArrayFilters: &options.ArrayFilters{
			Filters: []interface{}{
				bson.M{"elem.seatNumber": bson.M{"$in": seatNumbers}},
			},
		},
	}

    var update bson.M
    if newStatus == string(model.StatusHeld) || newStatus == "booked" { // 'booked' tương ứng với booking 'confirmed'
        if bookingID == nil || bookingID.IsZero() {
            return errors.New("bookingID is required when setting seat status to held or booked")
        }
        update = bson.M{"$set": bson.M{
            "seats.$[elem].status":    newStatus,
            "seats.$[elem].bookingId": *bookingID,
            "updatedAt": time.Now(),
        }}
    } else if newStatus == "available" {
         update = bson.M{
            "$set": bson.M{
                "seats.$[elem].status": newStatus,
                 "updatedAt": time.Now(),
            },
            "$unset": bson.M{"seats.$[elem].bookingId": ""}, // Xóa field bookingId
        }
    } else {
        return fmt.Errorf("unsupported new seat status for bulk update: %s", newStatus)
    }


	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": tripID}, update, &arrayFilters)
	if err != nil {
		return fmt.Errorf("failed to update seats for trip %s: %w", tripID.Hex(), err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("no trip found with ID %s to update seats", tripID.Hex())
	}
	// result.ModifiedCount có thể không phản ánh đúng số ghế được cập nhật nếu trạng thái không đổi.
	// Cần kiểm tra logic kỹ hơn nếu muốn biết chính xác bao nhiêu ghế đã thay đổi.
	log.Printf("Updated seats for trip %s. Matched: %d, Modified: %d", tripID.Hex(), result.MatchedCount, result.ModifiedCount)

	return nil
}


// Đảm bảo các hằng số collection name đã được định nghĩa
const vehicleCollectionName = "vehicles"
const companyCollectionName = "companies"

// Helper function để tạo index cho collection trips
func EnsureTripIndexes() error {
	tripRepo := NewTripRepository().(*tripRepositoryImpl) // Cần cast để truy cập collection
	
	indexes := []mongo.IndexModel{
		{ // Index cho tìm kiếm
			Keys: bson.D{
				{Key: "route.from.name", Value: 1},
				{Key: "route.to.name", Value: 1},
				{Key: "departureTime", Value: 1},
				{Key: "status", Value: 1},
			},
		},
		{ // Index cho seats để cập nhật
			Keys: bson.D{
				{Key: "seats.seatNumber", Value: 1},
				{Key: "seats.status", Value: 1},
			},
		},
        { // Index cho CompanyID
            Keys: bson.D{{Key: "companyId", Value: 1}},
        },
	}

	_, err := tripRepo.collection.Indexes().CreateMany(context.Background(), indexes)
	if err != nil {
		return fmt.Errorf("failed to create indexes for trips collection: %w", err)
	}
	fmt.Println("Successfully created indexes for trips collection.")
	return nil
}