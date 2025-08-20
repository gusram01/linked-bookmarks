package models

import "gorm.io/gorm"

type UserLink struct {
	gorm.Model
	UserID uint `gorm:"uniqueIndex:idx_user_link_unique"`
	LinkID uint `gorm:"uniqueIndex:idx_user_link_unique"`
}
