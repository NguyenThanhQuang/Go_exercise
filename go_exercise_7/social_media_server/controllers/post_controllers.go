package controllers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"social_media_server/config"
	"social_media_server/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// PostController ...
type PostController struct{}

// NewPostController ...
func NewPostController() *PostController {
	return &PostController{}
}

// @Summary Get all posts
// @Description Get a list of all posts with their comments
// @Tags posts
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Post "Successfully retrieved list of posts"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /posts [get]
func (pc *PostController) GetPosts(c *gin.Context) {
	// ... (code xử lý như trước)
	// 1. Thử lấy từ Cache trước
	cachedPosts, err := config.RDB.Get(config.Ctx, allPostsCacheKey).Result()
	if err == nil { // Cache hit!
		log.Println("Cache hit for GetPosts")
		var posts []models.Post
		if err := json.Unmarshal([]byte(cachedPosts), &posts); err == nil {
			c.JSON(http.StatusOK, posts)
			return
		}
		log.Println("Error unmarshalling cached posts:", err)
	} else if err != redis.Nil {
		log.Println("Redis error on GetPosts:", err)
	} else {
		log.Println("Cache miss for GetPosts")
	}

	var posts []models.Post
	if err := config.DB.Preload("Comments").Order("created_at DESC").Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve posts"})
		return
	}

	postsJSON, err := json.Marshal(posts)
	if err != nil {
		log.Println("Error marshalling posts for cache:", err)
	} else {
		cacheDuration := 0
		err = config.RDB.Set(config.Ctx, allPostsCacheKey, postsJSON, time.Duration(cacheDuration)).Err()
		if err != nil {
			log.Println("Error setting posts to cache:", err)
		} else {
			log.Println("Posts cached successfully for GetPosts")
		}
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
	// ... (code xử lý như trước)
	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	if err := config.RDB.Del(config.Ctx, allPostsCacheKey).Err(); err != nil {
		log.Println("Error deleting all_posts cache after creating post:", err)
	} else {
		log.Println("Cache all_posts invalidated after creating post")
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
	// ... (code xử lý như trước)
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	cacheKey := postCacheKeyPrefix + idStr

	cachedPost, err := config.RDB.Get(config.Ctx, cacheKey).Result()
	if err == nil {
		log.Printf("Cache hit for GetPost ID: %s\n", idStr)
		var post models.Post
		if err := json.Unmarshal([]byte(cachedPost), &post); err == nil {
			c.JSON(http.StatusOK, post)
			return
		}
		log.Printf("Error unmarshalling cached post ID %s: %v\n", idStr, err)
	} else if err != redis.Nil {
		log.Printf("Redis error on GetPost ID %s: %v\n", idStr, err)
	} else {
		log.Printf("Cache miss for GetPost ID: %s\n", idStr)
	}

	var post models.Post
	if err := config.DB.Preload("Comments").First(&post, uint(id)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve post"})
		return
	}

	postJSON, err := json.Marshal(post)
	if err != nil {
		log.Printf("Error marshalling post ID %s for cache: %v\n", idStr, err)
	} else {
		cacheDuration := 0
		err = config.RDB.Set(config.Ctx, cacheKey, postJSON, time.Duration(cacheDuration)).Err()
		if err != nil {
			log.Printf("Error setting post ID %s to cache: %v\n", idStr, err)
		} else {
			log.Printf("Post ID %s cached successfully\n", idStr)
		}
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
	// ... (code xử lý như trước)
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

	var postUpdates models.Post // Chỉ cần Title và Content từ request body
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

	cacheKey := postCacheKeyPrefix + idStr
	if errs := config.RDB.Del(config.Ctx, allPostsCacheKey, cacheKey).Err(); errs != nil {
		log.Printf("Error deleting cache after updating post ID %s: %v\n", idStr, errs)
	} else {
		log.Printf("Cache for all_posts and post ID %s invalidated after update\n", idStr)
	}
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
	// ... (code xử lý như trước)
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
		log.Printf("Error deleting comments for post ID %s: %v\n", idStr, err)
	}

	if err := config.DB.Delete(&models.Post{}, uint(id)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}

	cacheKey := postCacheKeyPrefix + idStr
	if errs := config.RDB.Del(config.Ctx, allPostsCacheKey, cacheKey).Err(); errs != nil {
		log.Printf("Error deleting cache after deleting post ID %s: %v\n", idStr, errs)
	} else {
		log.Printf("Cache for all_posts and post ID %s invalidated after delete\n", idStr)
	}
	c.JSON(http.StatusOK, gin.H{"message": "Post and associated comments deleted successfully"})
}