package controllers

import (
	// "encoding/json" // Không cần nữa nếu không cache
	"errors"
	"net/http"
	"social_media_server/config"
	"social_media_server/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PostController struct{}

func NewPostController() *PostController {
	return &PostController{}
}

// Bỏ các hằng số cache
// const (
// 	allPostsCacheKey = "all_posts"
// 	postCacheKeyPrefix = "post:"
// 	cacheDuration    = 5 * time.Minute
// )

// @Summary Get all posts
// @Description Get a list of all posts with their comments
// @Tags posts
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Post "Successfully retrieved list of posts"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /posts [get]
func (pc *PostController) GetPosts(c *gin.Context) {
	var posts []models.Post
	if err := config.DB.Preload("Comments").Order("created_at DESC").Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve posts"})
		return
	}
	c.JSON(http.StatusOK, posts)
}

// @Summary Create a new post
// @Description Create a new post with title and content
// @Tags posts
// @Accept  json
// @Produce  json
// @Param post body models.Post true "Post object that needs to be created"
// @Success 201 {object} models.Post "Successfully created post"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /posts [post]
func (pc *PostController) CreatePost(c *gin.Context) {
	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	c.JSON(http.StatusCreated, post)
}

// @Summary Get a single post by ID
// @Description Get details of a specific post by its ID, including comments
// @Tags posts
// @Accept  json
// @Produce  json
// @Param id path int true "Post ID"
// @Success 200 {object} models.Post "Successfully retrieved post"
// @Failure 400 {object} map[string]string "Invalid post ID"
// @Failure 404 {object} map[string]string "Post not found"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /posts/{id} [get]
func (pc *PostController) GetPost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	var post models.Post
	// Chỉ lấy từ DB
	if err := config.DB.Preload("Comments").First(&post, uint(id)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve post"})
		return
	}
	c.JSON(http.StatusOK, post)
}

// @Summary Update an existing post
// @Description Update title and content of an existing post by its ID
// @Tags posts
// @Accept  json
// @Produce  json
// @Param id path int true "Post ID"
// @Param post body models.Post true "Post object with updated fields (only Title and Content are used)"
// @Success 200 {object} models.Post "Successfully updated post"
// @Failure 400 {object} map[string]string "Invalid post ID or Bad Request"
// @Failure 404 {object} map[string]string "Post not found"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /posts/{id} [put]
func (pc *PostController) UpdatePost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	var post models.Post
	if err := config.DB.First(&post, uint(id)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve post for update"})
		return
	}

	var postUpdates models.Post
	if err := c.ShouldBindJSON(&postUpdates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post.Title = postUpdates.Title
	post.Content = postUpdates.Content

	if err := config.DB.Save(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		return
	}
	// Không còn invalidate cache
	c.JSON(http.StatusOK, post)
}

// @Summary Delete a post
// @Description Delete a post by its ID and its associated comments
// @Tags posts
// @Accept  json
// @Produce  json
// @Param id path int true "Post ID"
// @Success 200 {object} map[string]string "Message: Post and associated comments deleted successfully"
// @Failure 400 {object} map[string]string "Invalid post ID"
// @Failure 404 {object} map[string]string "Post not found"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /posts/{id} [delete]
func (pc *PostController) DeletePost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	var post models.Post
	if err := config.DB.First(&post, uint(id)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve post for deletion"})
		return
	}

	if err := config.DB.Where("post_id = ?", uint(id)).Delete(&models.Comment{}).Error; err != nil {
		// log.Printf("Error deleting comments for post ID %s: %v\n", idStr, err)
		// Có thể log lỗi này, nhưng không cần thiết nếu chỉ tạm bỏ Redis
	}

	if err := config.DB.Delete(&models.Post{}, uint(id)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}
	// Không còn invalidate cache
	c.JSON(http.StatusOK, gin.H{"message": "Post and associated comments deleted successfully"})
}