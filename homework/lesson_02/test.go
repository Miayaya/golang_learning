package lesson02
package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"blog-api/config"
	"blog-api/internal/router"
	"blog-api/model"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var (
	r         *gin.Engine
	token     string
	postID    string
	commentID = "1"
)

func init() {
	gin.SetMode(gin.TestMode)
	db := config.InitDB()
	db.AutoMigrate(&model.User{}, &model.Post{}, &model.Comment{})
	r = router.SetupRouter(db)
}

func TestBlogAPI(t *testing.T) {
	t.Run("Register", TestRegister)
	t.Run("Login", TestLogin)
	t.Run("CreatePost", TestCreatePost)
	t.Run("ListPosts", TestListPosts)
	t.Run("GetPost", TestGetPost)
	t.Run("UpdatePost", TestUpdatePost)
	t.Run("CreateComment", TestCreateComment)
	t.Run("ListComments", TestListComments)
	t.Run("DeleteComment", TestDeleteComment)
	t.Run("DeletePost", TestDeletePost)
}

func TestRegister(t *testing.T) {
	data := map[string]string{
		"username": "testuser",
		"password": "123456",
		"email":    "test@example.com",
	}
	jsonStr, _ := json.Marshal(data)

	req := httptest.NewRequest("POST", "/api/register", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestLogin(t *testing.T) {
	data := map[string]string{
		"username": "testuser",
		"password": "123456",
	}
	jsonStr, _ := json.Marshal(data)

	req := httptest.NewRequest("POST", "/api/login", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var res map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &res)

	dataMap, ok := res["data"].(map[string]interface{})
	assert.True(t, ok)
	token = dataMap["token"].(string)
	assert.NotEmpty(t, token)
}

func TestCreatePost(t *testing.T) {
	data := map[string]string{
		"title":   "test title",
		"content": "test content",
	}
	jsonStr, _ := json.Marshal(data)

	req := httptest.NewRequest("POST", "/api/posts", bytes.NewBuffer(jsonStr))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var res map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &res)

	dataMap, ok := res["data"].(map[string]interface{})
	assert.True(t, ok)

	// ✅ 修复：float64 转 string
	idFloat := dataMap["ID"].(float64)
	postID = fmt.Sprintf("%.0f", idFloat)
}

func TestListPosts(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/posts", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetPost(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/posts/"+postID, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdatePost(t *testing.T) {
	data := map[string]string{
		"title":   "updated title",
		"content": "updated content",
	}
	jsonStr, _ := json.Marshal(data)

	req := httptest.NewRequest("PUT", "/api/posts/"+postID, bytes.NewBuffer(jsonStr))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateComment(t *testing.T) {
	data := map[string]interface{}{
		"post_id": postID,
		"content": "test comment",
	}
	jsonStr, _ := json.Marshal(data)

	req := httptest.NewRequest("POST", "/api/comments", bytes.NewBuffer(jsonStr))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestListComments(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/posts/"+postID+"/comments", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteComment(t *testing.T) {
	req := httptest.NewRequest("DELETE", "/api/comments/"+commentID, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeletePost(t *testing.T) {
	req := httptest.NewRequest("DELETE", "/api/posts/"+postID, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}