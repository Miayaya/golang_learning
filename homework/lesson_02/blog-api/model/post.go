package model

import "gorm.io/gorm"

// Post corresponds to posts table
type Post struct {
	gorm.Model
	Title         string `gorm:"size:128;not null"`
	Content       string `gorm:"type:text;not null"`
	CommentStatus string `gorm:"default:'has_comments'"`
	UserID        uint   `gorm:"not null"` // Foreign key to User
	User          User
	Comments      []Comment
}

// AfterCreate hook: increase user post count
func (p *Post) AfterCreate(tx *gorm.DB) error {
	return tx.Model(&User{}).
		Where("id = ?", p.UserID).
		UpdateColumn("post_count", gorm.Expr("post_count + 1")).Error
}
