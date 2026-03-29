package router

import (
	"blog-api/internal/api"
	"blog-api/internal/router/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORS())
	r.Use(middleware.Recovery())
	r.Use(middleware.Logger())

	// public
	r.POST("/api/register", api.Register(db))
	r.POST("/api/login", api.Login(db))
	r.GET("/api/posts", api.ListPosts(db))   // all post list
	r.GET("/api/posts/:id", api.GetPost(db)) // post detail

	//Get all comments with post
	r.GET("/api/posts/:id/comments", api.ListCommentsByPostID(db))

	// protected,need login
	auth := r.Group("/api")
	auth.Use(middleware.AuthMiddleware())
	{
		// posts
		auth.POST("/posts", api.CreatePost(db))
		auth.PUT("/posts/:id", api.UpdatePost(db))
		auth.DELETE("/posts/:id", api.DeletePost(db))

		// comments
		auth.POST("/comments", api.CreateComment(db))
		auth.DELETE("/comments/:id", api.DeleteComment(db))
	}

	return r
}
