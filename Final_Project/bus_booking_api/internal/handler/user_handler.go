package handler

import (
	"bus_booking_api/internal/middleware"
	"bus_booking_api/internal/service"
	"bus_booking_api/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetMyProfile lấy thông tin của người dùng đang đăng nhập.
// @Summary Get current user's profile
// @Description Retrieves the profile information of the authenticated user
// @Tags Users
// @Security ApiKeyAuth
// @Produce  json
// @Success 200 {object} utils.APIResponse{data=model.User} "Successfully retrieved profile"
// @Failure 401 {object} utils.APIResponse "Unauthorized"
// @Failure 404 {object} utils.APIResponse "User not found"
// @Failure 500 {object} utils.APIResponse "Internal server error"
// @Router /users/me [get]
func (h *UserHandler) GetMyProfile(c *gin.Context) {
	authPayload, exists := middleware.GetAuthPayload(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized: Missing auth payload")
		return
	}

	userID, err := primitive.ObjectIDFromHex(authPayload.UserID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized: Invalid user ID in token")
		return
	}

	user, err := h.userService.GetUserProfile(userID)
	if err != nil {
		if err.Error() == "user not found" {
			utils.ErrorResponse(c, http.StatusNotFound, "User not found")
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get user profile: "+err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User profile retrieved successfully", user)
}