package routes

import (
	"social_media_server/controllers" 

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	postController := controllers.NewPostController()
	commentController := controllers.NewCommentController()

	postRoutes := router.Group("/posts")
	{
		postRoutes.GET("/", postController.GetPosts)
		postRoutes.POST("/", postController.CreatePost)
		postRoutes.GET("/:id", postController.GetPost)
		postRoutes.PUT("/:id", postController.UpdatePost)
		postRoutes.DELETE("/:id", postController.DeletePost)
	}

	commentRoutes := router.Group("/comments")
	{
		commentRoutes.POST("/", commentController.CreateComment)
		commentRoutes.PUT("/:id", commentController.UpdateComment)
		commentRoutes.DELETE("/:id", commentController.DeleteComment)
	}

	return router
}