package models

import "gorm.io/gorm"

type TagLink struct {
	gorm.Model
	TagID  uint `gorm:"uniqueIndex:idx_tag_link_unique"`
	LinkID uint `gorm:"uniqueIndex:idx_tag_link_unique"`
}
