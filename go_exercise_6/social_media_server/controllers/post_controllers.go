package controllers

import (
	"net/http"
	"social_media_server/config"
	"social_media_server/models"

	"github.com/gin-gonic/gin"
)

type PostController struct{}

func NewPostController() *PostController {
	return &PostController{}
}

func (ctl *PostController) CreatePost(c *gin.Context) {
	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.DB.Create(&post)
	c.JSON(http.StatusOK, post)
}

func (ctl *PostController) GetPosts(c *gin.Context) {
	var posts []models.Post
	config.DB.Preload("Comments").Find(&posts)
	c.JSON(http.StatusOK, posts)
}

func (ctl *PostController) GetPostByID(c *gin.Context) {
	var post models.Post
	id := c.Param("id")
	if err := config.DB.Preload("Comments").First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	c.JSON(http.StatusOK, post)
}

func (ctl *PostController) UpdatePost(c *gin.Context) {
	var post models.Post
	id := c.Param("id")
	if err := config.DB.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Save(&post)
	c.JSON(http.StatusOK, post)
}

func (ctl *PostController) DeletePost(c *gin.Context) {
	var post models.Post
	id := c.Param("id")
	if err := config.DB.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	config.DB.Delete(&post)
	c.JSON(http.StatusOK, gin.H{"message": "Post deleted"})
}

