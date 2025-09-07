package models

import (
	"time"

	"gorm.io/gorm"
)

type TagLink struct {
	gorm.Model
	TagID  uint `gorm:"uniqueIndex:idx_tag_link_unique"`
	LinkID uint `gorm:"uniqueIndex:idx_tag_link_unique"`
}

func (lt *TagLink) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	lt.CreatedAt = now
	lt.UpdatedAt = now
	return nil
}

func (lt *TagLink) BeforeUpdate(tx *gorm.DB) error {
	lt.UpdatedAt = time.Now()
	return nil
}
