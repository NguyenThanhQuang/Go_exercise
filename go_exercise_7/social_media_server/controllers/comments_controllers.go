package controllers

import (
	"errors"
	"log"
	"net/http"
	"social_media_server/config"
	"social_media_server/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CommentController struct{}

func NewCommentController() *CommentController {
	return &CommentController{}
}

// @Summary Create a new comment for a post
// @Description Create a new comment with content and associate it with a PostID
// @Tags comments
// @Accept  json
// @Produce  json
// @Param comment body models.Comment true "Comment object that needs to be created (ensure PostID is valid)"
// @Success 201 {object} models.Comment "Successfully created comment"
// @Failure 400 {object} map[string]string "Bad Request (e.g., missing content or PostID)"
// @Failure 404 {object} map[string]string "Post not found for the given PostID"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /comments [post]
func (cc *CommentController) CreateComment(c *gin.Context) {
	var comment models.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if comment.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Comment content cannot be empty"})
		return
	}
	if comment.PostID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "PostID is required to create a comment"})
		return
	}

	var post models.Post
	if err := config.DB.First(&post, comment.PostID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found, cannot create comment"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking post existence"})
		return
	}

	if err := config.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	postCacheKey := postCacheKeyPrefix + strconv.FormatUint(uint64(comment.PostID), 10)
	if errs := config.RDB.Del(config.Ctx, allPostsCacheKey, postCacheKey).Err(); errs != nil {
		log.Printf("Error deleting cache after creating comment for post ID %d: %v\n", comment.PostID, errs)
	} else {
		log.Printf("Cache for all_posts and post ID %d invalidated after creating comment\n", comment.PostID)
	}

	c.JSON(http.StatusCreated, comment)
}

// @Summary Update an existing comment
// @Description Update the content of an existing comment by its ID
// @Tags comments
// @Accept  json
// @Produce  json
// @Param id path int true "Comment ID"
// @Param comment body models.Comment true "Comment object with updated content (only Content is used)"
// @Success 200 {object} models.Comment "Successfully updated comment"
// @Failure 400 {object} map[string]string "Invalid comment ID or Bad Request (e.g., empty content)"
// @Failure 404 {object} map[string]string "Comment not found"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /comments/{id} [put]
func (cc *CommentController) UpdateComment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	var comment models.Comment
	if err := config.DB.First(&comment, uint(id)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve comment for update"})
		return
	}

	var commentUpdates models.Comment 
	if err := c.ShouldBindJSON(&commentUpdates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if commentUpdates.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Comment content cannot be empty for update"})
		return
	}

	originalPostID := comment.PostID
	comment.Content = commentUpdates.Content

	if err := config.DB.Save(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update comment"})
		return
	}

	postCacheKey := postCacheKeyPrefix + strconv.FormatUint(uint64(originalPostID), 10)
	if errs := config.RDB.Del(config.Ctx, allPostsCacheKey, postCacheKey).Err(); errs != nil {
		log.Printf("Error deleting cache after updating comment ID %d (post ID %d): %v\n", id, originalPostID, errs)
	} else {
		log.Printf("Cache for all_posts and post ID %d invalidated after updating comment ID %d\n", originalPostID, id)
	}

	c.JSON(http.StatusOK, comment)
}

// @Summary Delete a comment
// @Description Delete a comment by its ID
// @Tags comments
// @Accept  json
// @Produce  json
// @Param id path int true "Comment ID"
// @Success 200 {object} map[string]string "Message: Comment deleted successfully"
// @Failure 400 {object} map[string]string "Invalid comment ID"
// @Failure 404 {object} map[string]string "Comment not found"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /comments/{id} [delete]
func (cc *CommentController) DeleteComment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	var comment models.Comment
	if err := config.DB.First(&comment, uint(id)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve comment for deletion"})
		return
	}
	originalPostID := comment.PostID 

	if err := config.DB.Delete(&models.Comment{}, uint(id)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}

	postCacheKey := postCacheKeyPrefix + strconv.FormatUint(uint64(originalPostID), 10)
	if errs := config.RDB.Del(config.Ctx, allPostsCacheKey, postCacheKey).Err(); errs != nil {
		log.Printf("Error deleting cache after deleting comment ID %d (post ID %d): %v\n", id, originalPostID, errs)
	} else {
		log.Printf("Cache for all_posts and post ID %d invalidated after deleting comment ID %d\n", originalPostID, id)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}

const (
	allPostsCacheKey = "all_posts"
	postCacheKeyPrefix = "post:"
)