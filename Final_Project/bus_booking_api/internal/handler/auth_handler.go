package handler

import (
	"bus_booking_api/internal/model"
	"bus_booking_api/internal/service"
	"bus_booking_api/internal/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register xử lý request đăng ký.
// @Summary Register a new user
// @Description Create a new user account
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param   user body model.RegisterRequest true "User Registration Info"
// @Success 201 {object} utils.APIResponse{data=model.User} "User registered successfully"
// @Failure 400 {object} utils.APIResponse "Invalid input or user already exists"
// @Failure 500 {object} utils.APIResponse "Internal server error"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	user, err := h.authService.RegisterUser(req)
	if err != nil {
		if err.Error() == "email already registered" || err.Error() == "phone number already registered" {
			utils.ErrorResponse(c, http.StatusConflict, err.Error())
		} else if strings.Contains(err.Error(), "error creating user: email or phone already exists") {
			utils.ErrorResponse(c, http.StatusConflict, "Email or phone already exists")
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to register user: "+err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "User registered successfully", user)
}

// Login xử lý request đăng nhập.
// @Summary Log in a user
// @Description Authenticate a user and return a JWT token
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param   credentials body model.LoginRequest true "User Login Credentials"
// @Success 200 {object} utils.APIResponse{data=model.LoginResponse} "Login successful"
// @Failure 400 {object} utils.APIResponse "Invalid input"
// @Failure 401 {object} utils.APIResponse "Invalid credentials"
// @Failure 500 {object} utils.APIResponse "Internal server error"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	response, err := h.authService.LoginUser(req)
	if err != nil {
		if err.Error() == "invalid credentials" {
			utils.ErrorResponse(c, http.StatusUnauthorized, err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to login: "+err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Login successful", response)
}