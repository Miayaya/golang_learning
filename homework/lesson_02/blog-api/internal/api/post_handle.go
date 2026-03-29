package api

import (
	"blog-api/model"
	"blog-api/pkgs/resp"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreatePost 创建文章（必须登录）
func CreatePost(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var post model.Post
		if err := c.ShouldBindJSON(&post); err != nil {
			resp.Fail(c, http.StatusBadRequest, "invalid params")
			return
		}

		// 从JWT拿当前登录用户ID
		userId := uint(c.GetFloat64("userId"))
		post.UserID = userId

		if err := db.Create(&post).Error; err != nil {
			resp.Fail(c, http.StatusInternalServerError, "create failed")
			// c.JSON(http.StatusInternalServerError, gin.H{"error": "create failed"})
			return
		}

		resp.Success(c, post)
	}
}

// ListPosts 获取所有文章列表
func ListPosts(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var posts []model.Post
		// db.Find(&posts)
		if err := db.Find(&posts).Error; err != nil {
			resp.Fail(c, http.StatusInternalServerError, "query failed")
			return
		}
		// c.JSON(http.StatusOK, gin.H{"data": posts})
		resp.Success(c, posts)
	}
}

// GetPost 获取单篇文章详情
func GetPost(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var post model.Post

		if err := db.First(&post, id).Error; err != nil {
			resp.Fail(c, http.StatusNotFound, "post not found")
			// c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
			return
		}

		resp.Success(c, post)
	}
}

// UpdatePost 更新文章（仅作者）
func UpdatePost(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		userId := uint(c.GetFloat64("userId"))

		var post model.Post
		if err := db.First(&post, id).Error; err != nil {
			resp.Fail(c, http.StatusNotFound, "post not found")
			// c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
			return
		}

		// 校验作者
		if post.UserID != userId {
			resp.Fail(c, http.StatusForbidden, "permission denied")
			// c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
			return
		}

		// 只更新标题和内容
		var input model.Post
		if err := c.ShouldBindJSON(&input); err != nil {
			resp.Fail(c, http.StatusBadRequest, "invalid params")
			// c.JSON(http.StatusBadRequest, gin.H{"error": "invalid params"})
			return
		}

		db.Model(&post).Updates(map[string]interface{}{
			"title":   input.Title,
			"content": input.Content,
		})

		resp.Success(c, nil)
	}
}

// DeletePost 删除文章（仅作者）
func DeletePost(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		userId := uint(c.GetFloat64("userId"))

		var post model.Post
		if err := db.First(&post, id).Error; err != nil {
			resp.Fail(c, http.StatusNotFound, "post not found")
			// c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
			return
		}

		if post.UserID != userId {
			resp.Fail(c, http.StatusForbidden, "permission denied")
			// c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
			return
		}

		if err := db.Delete(&post).Error; err != nil {
			resp.Fail(c, http.StatusInternalServerError, "delete failed")
			return
		}
		resp.Success(c, nil)
	}
}
