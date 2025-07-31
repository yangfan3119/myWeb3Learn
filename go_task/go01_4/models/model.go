package models

import (
	"time"

	"gorm.io/gorm"
)

var mdb *gorm.DB

type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoCreateTime" json:"update_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Username string `gorm:"size:100;unique;not null" json:"username"`
	Password string `gorm:"size:256;not null" json:"password"`
	Email    string `gorm:"size:128;unique;not null" jsom:"email"`
}

type Post struct {
	gorm.Model
	Title   string `gorm:"not null" json:"title"`
	Content string `gorm:"not null" json:"content"`
	UserID  uint   `json:"user_id"`
	User    User
}

type Comment struct {
	gorm.Model
	Content string `gorm:"not null"`
	UserID  uint
	User    User
	PostID  uint
	Post    Post `gorm:"constraint:OnDelete:CASCADE"`
}

func MigrateDB(db *gorm.DB) {
	mdb = db
	mdb.Migrator().DropTable(&User{}, &Post{}, &Comment{})
	mdb.AutoMigrate(&User{}, &Post{}, &Comment{})
}
