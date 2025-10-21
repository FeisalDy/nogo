package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Username  *string `json:"username"`
	Email     string  `json:"email" gorm:"unique;not null;index"`
	Password  *string `json:"-"`
	AvatarURL *string `json:"avatar_url"`
	Bio       *string `json:"bio" gorm:"type:text"`
	Status    string  `json:"status" gorm:"default:'active';index"`
}
