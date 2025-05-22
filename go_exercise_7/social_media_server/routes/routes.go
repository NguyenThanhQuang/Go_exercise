package routes

import (
	"social_media_server/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	postCtl := controllers.NewPostController()
	commentCtl := controllers.NewCommentController()

	r.POST("/posts", postCtl.CreatePost)
	r.GET("/posts", postCtl.GetPosts)
	r.GET("/posts/:id", postCtl.GetPostByID)
	r.PUT("/posts/:id", postCtl.UpdatePost)
	r.DELETE("/posts/:id", postCtl.DeletePost)

	r.GET("/comments", commentCtl.GetComments)
	r.POST("/comments", commentCtl.CreateComment)
	r.PUT("/comments/:id", commentCtl.UpdateComment)
	r.DELETE("/comments/:id", commentCtl.DeleteComment)

	return r
}