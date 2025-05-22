package controllers

import (
	"net/http"
	"social_media_server/config"
	"social_media_server/models"

	"github.com/gin-gonic/gin"
)

type CommentController struct{}

func NewCommentController() *CommentController {
	return &CommentController{}
}

func (ctl *CommentController) CreateComment(c *gin.Context) {
	var comment models.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.DB.Create(&comment)
	c.JSON(http.StatusOK, comment)
}

func (ctl *CommentController) GetComments(c *gin.Context) {
	var comments []models.Comment
	config.DB.Find(&comments)
	c.JSON(http.StatusOK, comments)
}

func (ctl *CommentController) UpdateComment(c *gin.Context) {
	var comment models.Comment
	id := c.Param("id")
	if err := config.DB.First(&comment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Save(&comment)
	c.JSON(http.StatusOK, comment)
}

func (ctl *CommentController) DeleteComment(c *gin.Context) {
	var comment models.Comment
	id := c.Param("id")
	if err := config.DB.First(&comment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	config.DB.Delete(&comment)
	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted"})
}