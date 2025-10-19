package model

import (
	"gorm.io/gorm"
)

// UserRole represents the many-to-many relationship between users and roles
type UserRole struct {
	gorm.Model
	UserID uint `gorm:"primaryKey;index:idx_user_role" json:"user_id"`
	RoleID uint `gorm:"primaryKey;index:idx_user_role" json:"role_id"`
}
