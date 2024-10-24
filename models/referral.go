package models

import (
	"gorm.io/gorm"
	"time"
)


type Referral struct {
	gorm.Model
	UserID uint      `gorm:"unique;not null"`
	Code   string    `gorm:"unique;not null"`
	Expiry time.Time `gorm:"not null"`
}
