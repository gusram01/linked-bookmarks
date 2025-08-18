package models

import (
	"github.com/gusram01/linked-bookmarks/internal/onboarding/infra/models"
	"gorm.io/gorm"
)

type Link struct {
	gorm.Model
	Url   string        `gorm:"index:idx_link_models_url,unique"`
	Users []models.User `gorm:"many2many:user_links;"`
}
