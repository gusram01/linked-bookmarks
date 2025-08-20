package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	AuthID string `gorm:"index:idx_user_auth_id,unique"`
}
