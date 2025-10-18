package migrations

import (
	"gorm.io/gorm"
)

// Genre model for migration 004
type Genre struct {
	gorm.Model

	Name        string  `json:"name" gorm:"unique;not null"`
	Slug        string  `json:"slug" gorm:"unique;not null"`
	Description *string `json:"description"`
}

type Tag struct {
	gorm.Model

	Name        string  `json:"name" gorm:"unique;not null"`
	Slug        string  `json:"slug" gorm:"unique;not null"`
	Description *string `json:"description"`
}

func Migration006CreateGenresAndTags() Migration {
	return Migration{
		ID:          "006_create_genres_and_tags",
		Description: "Create genres and tags table",
		Up: func(db *gorm.DB) error {
			return db.AutoMigrate(&Genre{}, &Tag{})
		},
		Down: func(db *gorm.DB) error {
			return db.Migrator().DropTable(&Genre{}, &Tag{})
		},
	}
}
