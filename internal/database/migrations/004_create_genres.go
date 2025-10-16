package migrations

import (
	"time"

	"gorm.io/gorm"
)

// Genre model for migration 004
type Genre struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`

	// Basic genre fields
	Name        string `json:"name" gorm:"unique;not null"`
	Description string `json:"description"`
	Slug        string `json:"slug" gorm:"unique;not null"`
}

// Migration004CreateGenres creates the genres table
func Migration004CreateGenres() Migration {
	return Migration{
		ID:          "004_create_genres",
		Description: "Create genres table",
		Up: func(db *gorm.DB) error {
			return db.AutoMigrate(&Genre{})
		},
		Down: func(db *gorm.DB) error {
			return db.Migrator().DropTable(&Genre{})
		},
	}
}
