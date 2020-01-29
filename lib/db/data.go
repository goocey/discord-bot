package db

import (
	"github.com/jinzhu/gorm"
)

// User user
type User struct {
	gorm.Model
	UserID   string `gorm:"not null"`
	UserName string `gorm:"not null"`
}

// UserJoin joinカウント用
type UserJoin struct {
	gorm.Model
	UserID     string `gorm: "unique;not null;"`
	UserName   string `gorm:"not null"`
	JoinCount  int    `gorm:"not null" sql:"DEFAULT:0"`
	LeaveCount int    `gorm:"not null" sql:"DEFAULT:0"`
}
