package main

import (
	"fmt"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type User struct {
	gorm.Model
	Name      string
	Email     string
	Posts     []Post
	PostCount uint `gorm:"default:0"`
}
type Post struct {
	gorm.Model
	Title         string
	Content       string
	Type          string
	UserID        uint
	User          User
	CommentStatus string `gorm:"default:comment"`
	Comments      []Comment
}

type Comment struct {
	gorm.Model
	Content string
	PostID  uint
	Post    Post
}

func newSQLiteDB(filename string) (*gorm.DB, error) {
	dbDir := "./db"

	_, err := os.Stat(dbDir)
	if os.IsNotExist(err) {
		if err := os.Mkdir(dbDir, 0755); err != nil {
			panic("创建 db 文件夹失败: " + err.Error())
		}
	}

	dbPath := filepath.Join(dbDir, filename)

	return gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		// Logger configuration
		Logger: logger.Default.LogMode(logger.Info),

		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "",    // Prefix for all table names (e.g., "app_")
			SingularTable: false, // Use singular table names (User -> user instead of users)
			NoLowerCase:   false, // Disable automatic lowercasing
			NameReplacer:  nil,   // Custom name replacer function
		},
	})
}

func getUserPostsWithComments(db *gorm.DB, userID uint) (User, error) {
	var user User
	err := db.
		Preload("Posts.Comments").
		Find(&user, userID).Error
	return user, err
}

func getMostCommentedPost(db *gorm.DB) (Post, error) {
	var post Post
	err := db.Preload("Comments").
		Table("post").
		Select("posts.*, (SELECT COUNT(*) FROM comments WHERE comments.post_id = posts.id) AS comment_count").
		Order("comment_count DESC").
		Limit(1).
		Find(&post).Error

	return post, err
}

func (p *Post) AfterCreate(tx *gorm.DB) error {
	return tx.Model(&User{}).
		Where("id = ?", p.UserID).
		UpdateColumn("post_count", gorm.Expr("post_count + 1")).Error
}

/*
func (c *Comment) AfterDelete(tx *gorm.DB) error {
	var CommetCount int64
	if err := tx.Model(&Comment{}).
		Where("post_id = ?", c.PostID).Count(&CommetCount).Error; err != nil {
		return err
	}
	if CommetCount == 0 {
		return tx.Model(&Post{}).
			UpdateColumn("comment_status", "no comments").Error
	}
	return nil

}
*/
//  AfterDelete change to BeforeDelete
func (c *Comment) BeforeDelete(tx *gorm.DB) error {
	var commentCount int64

	if err := tx.Model(&Comment{}).
		Where("post_id = ?", c.PostID).
		Count(&commentCount).Error; err != nil {
		return err
	}

	if commentCount == 1 {
		return tx.Model(&Post{}).
			Where("id = ?", c.PostID). // need
			UpdateColumn("comment_status", "no comments").Error
	}

	return nil
}

func main() {
	db, err := newSQLiteDB("test.db")
	if err != nil {
		panic("db link failed: " + err.Error())
	}
	err = db.AutoMigrate(&User{}, &Post{}, &Comment{})
	if err != nil {
		panic("auto migrate failed:" + err.Error())
	}

	// Test 1: Create user and post, check PostCount increment
	user := User{Name: "John Doe", Email: "john@example.com"}
	db.Create(&user)
	fmt.Println("\n[Test 1] User created, ID =", user.ID)

	post1 := Post{
		Title:   "Go GORM Best Practices",
		Content: "GORM hooks are very useful",
		UserID:  user.ID,
	}
	db.Create(&post1)
	fmt.Println("[Test 1] Post created")

	var checkUser User
	db.First(&checkUser, user.ID)
	fmt.Println("[Test 1] User post count:", checkUser.PostCount)

	// Test 2: Create comment, delete it, check status change
	comment := Comment{
		Content: "Nice article!",
		PostID:  post1.ID,
	}
	db.Create(&comment)
	fmt.Println("\n[Test 2] Comment created")

	db.Delete(&comment)
	fmt.Println("[Test 2] Comment deleted")

	var checkPost Post
	db.First(&checkPost, post1.ID)

	fmt.Println("[Test 2] Post comment status:", checkPost.CommentStatus)

	// Test 3: Get user posts with comments
	userWithPosts, _ := getUserPostsWithComments(db, user.ID)
	fmt.Println("\n[Test 3] User:", userWithPosts.Name, "| Posts count:", len(userWithPosts.Posts))
	for _, p := range userWithPosts.Posts {
		fmt.Println("  Post:", p.Title, "| Comments:", len(p.Comments))
	}

	// Test 4: Get most commented post
	topPost, _ := getMostCommentedPost(db)
	fmt.Println("\n[Test 4] Most commented post ID:", topPost.ID)

	fmt.Println("\n✅ All tests completed!")
}
