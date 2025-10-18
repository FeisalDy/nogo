package migrations

import (
	"gorm.io/gorm"
)

// User model for migration 001
type User struct {
	gorm.Model

	UserName  *string `json:"user_name"`
	Email     string  `json:"email" gorm:"unique;not null"`
	Password  *string `json:"-"`
	AvatarURL *string `json:"avatar_url"`
	Bio       *string `json:"bio" gorm:"type:text"`
	Status    string  `json:"status" gorm:"default:'active';index"`
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
