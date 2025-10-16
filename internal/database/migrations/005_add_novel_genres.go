package migrations

import "gorm.io/gorm"

// NovelGenre junction table for migration 005
type NovelGenre struct {
	NovelID uint  `gorm:"primaryKey"`
	GenreID uint  `gorm:"primaryKey"`
	Novel   Novel `gorm:"foreignKey:NovelID"`
	Genre   Genre `gorm:"foreignKey:GenreID"`
}

// Migration005AddNovelGenres creates the novel-genre many-to-many relationship
func Migration005AddNovelGenres() Migration {
	return Migration{
		ID:          "005_add_novel_genres",
		Description: "Create novel-genre many-to-many relationship table",
		Up: func(db *gorm.DB) error {
			return db.AutoMigrate(&NovelGenre{})
		},
		Down: func(db *gorm.DB) error {
			return db.Migrator().DropTable(&NovelGenre{})
		},
	}
}
