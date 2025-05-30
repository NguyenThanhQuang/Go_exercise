package handler

import (
	"bus_booking_api/internal/model"
	"bus_booking_api/internal/service"
	"bus_booking_api/internal/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TripHandler struct {
	tripService service.TripService
}

func NewTripHandler(tripService service.TripService) *TripHandler {
	return &TripHandler{
		tripService: tripService,
	}
}

// SearchTrips
// @Summary Search for available trips
// @Description Search for trips based on origin, destination, and date
// @Tags Trips
// @Accept  json
// @Produce  json
// @Param from query string true "Origin location name"
// @Param to query string true "Destination location name"
// @Param date query string true "Date of travel (YYYY-MM-DD)"
// @Success 200 {object} utils.APIResponse{data=[]model.Trip} "List of found trips"
// @Failure 400 {object} utils.APIResponse "Invalid query parameters"
// @Failure 500 {object} utils.APIResponse "Internal server error"
// @Router /trips [get]
func (h *TripHandler) SearchTrips(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")
	date := c.Query("date")

	if from == "" || to == "" || date == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Query parameters 'from', 'to', and 'date' are required")
		return
	}

	trips, err := h.tripService.SearchTrips(c.Request.Context(), from, to, date)
	if err != nil {
		if err.Error() == "invalid date format, please use YYYY-MM-DD" {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to search trips: "+err.Error())
		}
		return
	}
    if trips == nil { // Đảm bảo trả về mảng rỗng thay vì null nếu không có trips
        trips = []model.Trip{}
    }
	utils.SuccessResponse(c, http.StatusOK, "Trips retrieved successfully", trips)
}

// GetTripDetails
// @Summary Get details of a specific trip
// @Description Get all details of a trip including route, stops, and seat map
// @Tags Trips
// @Accept  json
// @Produce  json
// @Param tripId path string true "Trip ID"
// @Success 200 {object} utils.APIResponse{data=model.Trip} "Trip details"
// @Failure 400 {object} utils.APIResponse "Invalid Trip ID"
// @Failure 404 {object} utils.APIResponse "Trip not found"
// @Failure 500 {object} utils.APIResponse "Internal server error"
// @Router /trips/{tripId} [get]
func (h *TripHandler) GetTripDetails(c *gin.Context) {
	tripID := c.Param("tripId")

	trip, err := h.tripService.GetTripDetails(c.Request.Context(), tripID)
	if err != nil {
		if err.Error() == "invalid trip ID format" {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		} else if err.Error() == "trip not found" {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get trip details: "+err.Error())
		}
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Trip details retrieved successfully", trip)
}

// CreateTrip (For Admin/Company - sẽ cần Auth Middleware sau)
// @Summary Create a new trip (Admin/Company)
// @Description Allows authenticated admins or company managers to create a new trip
// @Tags Trips_Admin
// @Accept  json
// @Produce  json
// @Param   trip body model.Trip true "Trip data"
// @Security ApiKeyAuth
// @Success 201 {object} utils.APIResponse{data=model.Trip} "Trip created successfully"
// @Failure 400 {object} utils.APIResponse "Invalid input"
// @Failure 401 {object} utils.APIResponse "Unauthorized"
// @Failure 403 {object} utils.APIResponse "Forbidden (User does not have permission)"
// @Failure 500 {object} utils.APIResponse "Internal server error"
// @Router /trips [post] // Sẽ được bảo vệ bởi Auth Middleware và Role Check
func (h *TripHandler) CreateTrip(c *gin.Context) {
	// TODO: Add AuthMiddleware and Role Check (Admin or CompanyAdmin)
	// claims, exists := middleware.GetAuthPayload(c)
	// if !exists || (claims.Role != model.RoleAdmin && claims.Role != model.RoleCompanyAdmin) {
	// 	utils.ErrorResponse(c, http.StatusForbidden, "You do not have permission to create a trip")
	// 	return
	// }

	var tripData model.Trip
	if err := c.ShouldBindJSON(&tripData); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	// Nếu là CompanyAdmin, companyId trong tripData phải khớp với companyId của admin đó
	// if claims.Role == model.RoleCompanyAdmin {
	//   companyIdFromToken, _ := primitive.ObjectIDFromHex(claims.CompanyID)
	// 	 if tripData.CompanyID != companyIdFromToken {
	// 		 utils.ErrorResponse(c, http.StatusForbidden, "You can only create trips for your own company")
	// 		 return
	// 	 }
	// }


	createdTrip, err := h.tripService.CreateTrip(c.Request.Context(), &tripData)
	if err != nil {
        // Cung cấp thông báo lỗi chi tiết hơn từ service
		utils.ErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Failed to create trip: %v", err))
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Trip created successfully", createdTrip)
}