package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"blog-api/config"
	"blog-api/internal/router"
	"blog-api/model"

	"github.com/gin-gonic/gin"
)

var (
	r           *gin.Engine
	token       string
	postID      string
	successList []string
	failList    []struct {
		Name    string
		Message string
	}
)

func init() {
	gin.SetMode(gin.TestMode)
	db := config.InitDB()
	db.AutoMigrate(&model.User{}, &model.Post{}, &model.Comment{})
	r = router.SetupRouter(db)
	successList = []string{}
	failList = []struct {
		Name    string
		Message string
	}{}
}

func pass(name string) {
	successList = append(successList, name)
}

func fail(name string, msg string) {
	failList = append(failList, struct {
		Name    string
		Message string
	}{name, msg})
}

func summary() {
	fmt.Println("==================================================")
	fmt.Println("              TEST SUMMARY")
	fmt.Println("==================================================")

	fmt.Println("\n✅ PASSED:")
	if len(successList) == 0 {
		fmt.Println("  None")
	} else {
		for _, name := range successList {
			fmt.Printf("  • %s\n", name)
		}
	}

	fmt.Println("\n❌ FAILED:")
	if len(failList) == 0 {
		fmt.Println("  None")
	} else {
		for _, item := range failList {
			fmt.Printf("  • %s → Reason: %s\n", item.Name, item.Message)
		}
	}

	fmt.Println("\n📊 STATISTICS:")
	fmt.Printf("  Total Passed: %d\n", len(successList))
	fmt.Printf("  Total Failed: %d\n", len(failList))
	fmt.Println("==================================================")
}

func TestBlogAPI(t *testing.T) {
	TestRegister(t)
	TestLogin(t)
	TestCreatePost(t)
	TestListPosts(t)
	TestGetPost(t)
	TestUpdatePost(t)
	TestCreateComment(t)
	TestListComments(t)
	TestDeleteComment(t)
	TestDeletePost(t)

	summary()
}

// 1 Register
func TestRegister(t *testing.T) {
	name := "User Register"
	jsonStr := `{"name":"testuser","password":"123456","email":"test@example.com"}`
	req := httptest.NewRequest("POST", "/api/register", bytes.NewBufferString(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		fail(name, fmt.Sprintf("HTTP status error: %d, response: %s", w.Code, w.Body.String()))
		return
	}
	pass(name)
}

// 2 Login
func TestLogin(t *testing.T) {
	name := "User Login"
	jsonStr := `{"username":"testuser","password":"123456"}`
	req := httptest.NewRequest("POST", "/api/login", bytes.NewBufferString(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		fail(name, fmt.Sprintf("HTTP status error: %d", w.Code))
		return
	}

	var res map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &res)
	token = res["token"].(string)
	if token == "" {
		fail(name, "Token not found")
		return
	}

	pass(name)
}

// 3 Create Post
func TestCreatePost(t *testing.T) {
	name := "Create Post"
	jsonStr := `{"title":"test title","content":"test content"}`
	req := httptest.NewRequest("POST", "/api/posts", bytes.NewBufferString(jsonStr))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		fail(name, fmt.Sprintf("HTTP status error: %d", w.Code))
		return
	}

	var res map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &res)
	data, ok := res["data"].(map[string]interface{})
	if !ok || data["ID"] == nil {
		fail(name, "Post ID not found")
		return
	}

	postID = fmt.Sprintf("%.0f", data["ID"].(float64))
	pass(name)
}

// 4 Post List
func TestListPosts(t *testing.T) {
	name := "Get Post List"
	req := httptest.NewRequest("GET", "/api/posts", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		fail(name, fmt.Sprintf("HTTP status error: %d", w.Code))
		return
	}
	pass(name)
}

// 5 Single Post
func TestGetPost(t *testing.T) {
	name := "Get Post Detail"
	req := httptest.NewRequest("GET", "/api/posts/"+postID, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		fail(name, fmt.Sprintf("HTTP status error: %d", w.Code))
		return
	}
	pass(name)
}

// 6 Update Post
func TestUpdatePost(t *testing.T) {
	name := "Update Post"
	jsonStr := `{"title":"updated","content":"updated content"}`
	req := httptest.NewRequest("PUT", "/api/posts/"+postID, bytes.NewBufferString(jsonStr))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		fail(name, fmt.Sprintf("HTTP status error: %d", w.Code))
		return
	}
	pass(name)
}

// 7 Create Comment
func TestCreateComment(t *testing.T) {
	name := "Create Comment"
	jsonStr := `{"post_id":` + postID + `,"content":"test comment"}`
	req := httptest.NewRequest("POST", "/api/comments", bytes.NewBufferString(jsonStr))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		fail(name, fmt.Sprintf("HTTP status error: %d", w.Code))
		return
	}
	pass(name)
}

// 8 Comment List
func TestListComments(t *testing.T) {
	name := "Get Comment List"
	req := httptest.NewRequest("GET", "/api/posts/"+postID+"/comments", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		fail(name, fmt.Sprintf("HTTP status error: %d", w.Code))
		return
	}
	pass(name)
}

// 9 Delete Comment
func TestDeleteComment(t *testing.T) {
	name := "Delete Comment"
	req := httptest.NewRequest("DELETE", "/api/comments/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		fail(name, fmt.Sprintf("HTTP status error: %d", w.Code))
		return
	}
	pass(name)
}

// 10 Delete Post
func TestDeletePost(t *testing.T) {
	name := "Delete Post"
	req := httptest.NewRequest("DELETE", "/api/posts/"+postID, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		fail(name, fmt.Sprintf("HTTP status error: %d", w.Code))
		return
	}
	pass(name)
}
