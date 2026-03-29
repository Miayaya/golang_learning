package api

import (
	"blog-api/model"
	"blog-api/pkgs/resp"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateComment
func CreateComment(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var comment model.Comment
		if err := c.ShouldBindJSON(&comment); err != nil {
			resp.Fail(c, http.StatusBadRequest, "invalid params")
			return
		}

		// Get login user id
		uid := c.GetFloat64("userId")
		comment.UserID = uint(uid)

		// Save comments
		if err := db.Create(&comment).Error; err != nil {
			resp.Fail(c, http.StatusInternalServerError, "create failed")
			return
		}

		resp.Success(c, nil)
	}
}

// ListCommentsByPostID
func ListCommentsByPostID(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		postID := c.Param("post_id")

		var comments []model.Comment
		// Find all comments in post_id
		if err := db.Where("post_id = ?", postID).Find(&comments).Error; err != nil {
			resp.Fail(c, http.StatusInternalServerError, "query failed")
			return
		}
		resp.Success(c, comments)
	}
}

// DeleteComment
func DeleteComment(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := db.Delete(&model.Comment{}, id).Error; err != nil {
			resp.Fail(c, http.StatusInternalServerError, "delete failed")
			return
		}
		resp.Success(c, nil)
	}
}
