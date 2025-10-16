package migrations

import (
	"time"

	"gorm.io/gorm"
)

// User model for migration 001
type User struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`

	// Basic user fields - keep it simple for now
	Name     string `json:"name" gorm:"not null"`
	Email    string `json:"email" gorm:"unique;not null"`
	Password string `json:"-" gorm:"not null"`        // Hidden from JSON
	Role     string `json:"role" gorm:"default:user"` // user, admin, author
	IsActive bool   `json:"is_active" gorm:"default:true"`
}

// Migration001CreateUsers creates the users table
func Migration001CreateUsers() Migration {
	return Migration{
		ID:          "001_create_users",
		Description: "Create users table with basic authentication fields",
		Up: func(db *gorm.DB) error {
			return db.AutoMigrate(&User{})
		},
		Down: func(db *gorm.DB) error {
			return db.Migrator().DropTable(&User{})
		},
	}
}
