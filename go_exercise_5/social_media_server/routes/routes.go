package routes

import (
	"social_media_server/controllers"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.RedirectTrailingSlash = false 

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour
	router.Use(cors.New(config))

	postController := controllers.NewPostController()
	commentController := controllers.NewCommentController()

	postRoutes := router.Group("/posts")
	{
		postRoutes.GET("", postController.GetPosts)        
		postRoutes.POST("", postController.CreatePost)       
		postRoutes.GET("/:id", postController.GetPost)   
		postRoutes.PUT("/:id", postController.UpdatePost)
		postRoutes.DELETE("/:id", postController.DeletePost) 
	}

	commentRoutes := router.Group("/comments")
	{
		commentRoutes.POST("", commentController.CreateComment) 
		commentRoutes.PUT("/:id", commentController.UpdateComment)
		commentRoutes.DELETE("/:id", commentController.DeleteComment)
	}

	return router
}