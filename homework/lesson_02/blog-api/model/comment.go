package model

import "gorm.io/gorm"

// Comment corresponds to comments table
type Comment struct {
	gorm.Model
	Content string `gorm:"type:text;not null"`
	UserID  uint   `gorm:"not null"`
	PostID  uint   `gorm:"not null"`
	User    User
	Post    Post
}

// AfterDelete hook: update post status when no comments left
func (c *Comment) BeforeDelete(tx *gorm.DB) error {
	var count int64
	if err := tx.Model(&Comment{}).
		Where("post_id = ?", c.PostID).
		Count(&count).Error; err != nil {
		return err
	}

	if count == 1 {
		return tx.Model(&Post{}).
			Where("id = ?", c.PostID).
			UpdateColumn("comment_status", "no_comments").Error
	}
	return nil
}
