package migrations

import (
	"time"

	"gorm.io/gorm"
)

// Chapter model for migration 003
type Chapter struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`

	// Basic chapter fields
	Title    string `json:"title" gorm:"not null"`
	Content  string `json:"content" gorm:"type:text"`
	Number   int    `json:"number" gorm:"not null"` // Chapter order
	IsPublic bool   `json:"is_public" gorm:"default:false"`

	// Required relationship - Novel
	NovelID uint  `json:"novel_id" gorm:"not null"`
	Novel   Novel `json:"novel" gorm:"foreignKey:NovelID"`
}

// Migration003CreateChapters creates the chapters table
func Migration003CreateChapters() Migration {
	return Migration{
		ID:          "003_create_chapters",
		Description: "Create chapters table with novel relationship",
		Up: func(db *gorm.DB) error {
			return db.AutoMigrate(&Chapter{})
		},
		Down: func(db *gorm.DB) error {
			return db.Migrator().DropTable(&Chapter{})
		},
	}
}
