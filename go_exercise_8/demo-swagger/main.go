// main.go
package main

import (
	"github.com/gin-gonic/gin"
	_ "demo-swagger/docs" // Import docs để swagger nhận biết
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files"
)

// @title Demo Swagger API
// @version 1.0
// @description Đây là ví dụ Swagger dùng Gin
// @host localhost:8080
// @BasePath /api/v1

func main() {
	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		v1.GET("/hello", HelloHandler)
	}

	// Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":8080")
}

// HelloHandler godoc
// @Summary Chào mừng đến với Swagger
// @Description Trả về câu chào
// @Tags example
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]string
// @Router /hello [get]
func HelloHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello from Swagger!",
	})
}