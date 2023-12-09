package authentication

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       int32  `gorm:"autoIncrement;primary_key;"`
	Name     string `gorm:"type:varchar(255);not null"`
	Email    string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`

	Otp_enabled  bool `gorm:"default:false;"`
	Otp_verified bool `gorm:"default:false;"`

	Otp_secret   string
	Otp_auth_url string
}

func InitDb(db *gorm.DB) error {
	return db.AutoMigrate(&User{})
}
