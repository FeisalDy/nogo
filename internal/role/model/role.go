package model

import "gorm.io/gorm"

// Role represents a user role in the system
type Role struct {
	gorm.Model

	Name        string  `json:"name" gorm:"unique;not null;index"`
	Description *string `json:"description" gorm:"type:text"`

	// Note: Don't define Users []User here to avoid circular dependency
	// Use the UserRole junction table in common/model instead
}
