package models

import (
	"gorm.io/gorm"
)

type Link struct {
	gorm.Model
	Url   string `gorm:"index:idx_link_models_url,unique"`
	Users []User `gorm:"many2many:user_links;joinForeignKey:LinkID;joinReferences:UserID"`
}
