package dao

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"size:20;unique;not null"`
	Password string `gorm:"size:255;not null"`
	Email    string `gorm:"size:50;unique;not null"`
}

type Post struct {
	gorm.Model
	Title    string `gorm:"size:50;not null"`
	Content  string `gorm:"type:text"`
	UserID   uint
	User     User
	Comments []Comment `gorm:"foreignkey:PostID;constraint:OnDelete:CASCADE;"`
}

type Comment struct {
	gorm.Model
	Content string `gorm:"type:text;not null"`
	UserID  uint
	User    User
	PostID  uint
	Post    Post
}
