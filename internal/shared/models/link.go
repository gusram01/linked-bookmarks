package models

import (
	"gorm.io/gorm"
)

type Link struct {
	gorm.Model
	Url      string `gorm:"index:idx_link_models_url,unique"`
	Summary  string `gorm:"text"`
	Attempts uint   `gorm:"default:0;not null"`
	Users    []User `gorm:"many2many:user_links;joinForeignKey:LinkID;joinReferences:UserID"`
	Tags     []Tag  `gorm:"many2many:tag_links;joinForeignKey:LinkID;joinReferences:TagID"`
}
