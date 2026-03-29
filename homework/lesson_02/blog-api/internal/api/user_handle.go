package api

import (
	"blog-api/model"
	"blog-api/pkgs/resp"
	"errors"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

/*
func Register(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user model.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

		if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
            return
        }
		user.Password = string(hash)

		if err := db.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "register failed"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "success"})
	}
}

*/

func Register(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Name     string `json:"name" binding:"required"`
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required,min=6"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			resp.Fail(c, http.StatusBadRequest, "invalid params")
			// c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 检查邮箱是否已存在
		var existingUser model.User
		if err := db.Where("Username = ?", input.Name).First(&existingUser).Error; err == nil {
			resp.Fail(c, http.StatusConflict, "email already registered")
			// c.JSON(http.StatusConflict, gin.H{"error": "name already registered"})
			return
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			resp.Fail(c, http.StatusInternalServerError, "database error")
			// c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
			return
		}

		// 密码哈希
		hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			resp.Fail(c, http.StatusInternalServerError, "failed to hash password")
			// c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
			return
		}

		user := model.User{
			Username: input.Name,
			Email:    input.Email,
			Password: string(hash),
		}
		if err := db.Create(&user).Error; err != nil {
			resp.Fail(c, http.StatusInternalServerError, "failed to create user")
			// c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
			return
		}

		resp.Success(c, gin.H{
			"user": gin.H{
				"id":    user.ID,
				"name":  user.Username,
				"email": user.Email,
			},
		})
	}
}

func Login(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input model.User
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var user model.User
		if err := db.Where("username = ?", input.Username).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "password error"})
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":  user.ID,
			"exp": time.Now().Add(time.Hour * 24).Unix(),
		})

		tokenStr, _ := token.SignedString([]byte("my_secret_key"))
		c.JSON(http.StatusOK, gin.H{"token": tokenStr})
	}
}
