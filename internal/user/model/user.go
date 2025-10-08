package model

import "gorm.io/gorm"

// User represents a user in the system
// We can add more fields as needed
// and also add validation tags
type User struct {
	gorm.Model
	Name  string `json:"name"`
	Email string `json:"email" gorm:"unique"`
}
