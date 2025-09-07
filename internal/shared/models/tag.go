package models

import "gorm.io/gorm"

type Tag struct {
	gorm.Model
	Name string `gorm:"index:idx_tag_name,unique"`
}
