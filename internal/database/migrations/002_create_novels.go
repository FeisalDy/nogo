package migrations

import (
	"time"

	"gorm.io/gorm"
)

// Novel model for migration 002
type Novel struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`

	// Basic novel fields
	Title       string `json:"title" gorm:"not null"`
	Description string `json:"description" gorm:"type:text"`
	Status      string `json:"status" gorm:"default:draft"` // draft, published, completed

	// Required relationship - Author (User)
	AuthorID uint `json:"author_id" gorm:"not null"`
	Author   User `json:"author" gorm:"foreignKey:AuthorID"`
}

// Migration002CreateNovels creates the novels table
func Migration002CreateNovels() Migration {
	return Migration{
		ID:          "002_create_novels",
		Description: "Create novels table with author relationship",
		Up: func(db *gorm.DB) error {
			return db.AutoMigrate(&Novel{})
		},
		Down: func(db *gorm.DB) error {
			return db.Migrator().DropTable(&Novel{})
		},
	}
}
