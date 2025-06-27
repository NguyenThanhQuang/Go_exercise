package main

import (
	"bus_booking_api/internal/config"
	"bus_booking_api/internal/handler"
	"bus_booking_api/internal/middleware"
	"bus_booking_api/internal/repository"
	"bus_booking_api/internal/service"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Load Cấu hình
	config.LoadConfig()

	// 2. Kết nối Database
	config.ConnectDB()

	// 3. Tạo Index
	if err := repository.EnsureUserIndexes(); err != nil {
		log.Printf("Warning: Could not ensure user indexes: %v.", err)
	}
	if err := repository.EnsureTripIndexes(); err != nil { // THÊM DÒNG NÀY
		log.Printf("Warning: Could not ensure trip indexes: %v.", err)
	}
    // Gọi EnsureCompanyIndexes, EnsureVehicleIndexes khi tạo các repo đó

	// 4. Khởi tạo Gin Engine
	router := gin.Default()
	// ... (CORS Middleware) ...

	// 5. Dependency Injection
	userRepo := repository.NewUserRepository()
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// === THÊM CHO TRIPS ===
	tripRepo := repository.NewTripRepository()
	tripService := service.NewTripService(tripRepo) // Truyền các repo khác nếu TripService cần
	tripHandler := handler.NewTripHandler(tripService)
	// ======================

	// 6. Định nghĩa Routes
	apiV1 := router.Group("/api")
	{
		authRoutes := apiV1.Group("/auth")
		{
			authRoutes.POST("/register", authHandler.Register)
			authRoutes.POST("/login", authHandler.Login)
		}

		userRoutes := apiV1.Group("/users")
		{
			userRoutes.GET("/me", middleware.AuthMiddleware(), userHandler.GetMyProfile)
		}

		// === THÊM ROUTES CHO TRIPS ===
		tripPublicRoutes := apiV1.Group("/trips")
		{
			tripPublicRoutes.GET("", tripHandler.SearchTrips)             // GET /api/trips?from=A&to=B&date=YYYY-MM-DD
			tripPublicRoutes.GET("/:tripId", tripHandler.GetTripDetails) // GET /api/trips/some-trip-id
		}

		// Ví dụ route tạo trip cho admin (sẽ thêm AuthMiddleware sau)
		tripAdminRoutes := apiV1.Group("/trips") // Có thể cùng group nếu method khác
		// Hoặc tạo group riêng /api/admin/trips
		// tripAdminRoutes.Use(middleware.AuthMiddleware(), middleware.RequireRole(model.RoleAdmin)) // Hoặc kiểm tra role trong handler
		{
			// Hiện tại để public để test, sau sẽ thêm middleware
			tripAdminRoutes.POST("", tripHandler.CreateTrip) // POST /api/trips
		}
		// ============================
	}

	// Health check
	router.GET("/health", func(c *gin.Context) { // Dòng này bạn có thể đã comment ra
		c.JSON(http.StatusOK, gin.H{"status": "UP"})
	})

	// 7. Chạy Server
	serverAddr := fmt.Sprintf(":%s", config.AppConfig.Port)
	log.Printf("Server is listening on port %s", config.AppConfig.Port) // Thay đổi log một chút để chắc chắn nó chạy tới đây
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to run server: %v", err) // Đảm bảo log lỗi này
	}
	log.Println("Server exiting...") // Dòng này không bao giờ nên được thấy nếu server chạy bình thường
}