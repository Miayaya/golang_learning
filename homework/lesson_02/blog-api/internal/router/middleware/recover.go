package middleware

import (
	"blog-api/pkgs/resp"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Recovery
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, err interface{}) {
		log.Printf("panic: %v", err)
		resp.Fail(c, http.StatusInternalServerError, "server error")
	})
}
