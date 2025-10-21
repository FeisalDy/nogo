package migrations

import "gorm.io/gorm"

// NovelGenre junction table for migration 005
type NovelGenre struct {
	NovelID uint   `gorm:"primaryKey;index"`
	GenreID uint   `gorm:"primaryKey;index"`
	Novel   *Novel `gorm:"foreignKey:NovelID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Genre   *Genre `gorm:"foreignKey:GenreID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type NovelTag struct {
	NovelID uint   `gorm:"primaryKey;index"`
	TagID   uint   `gorm:"primaryKey;index"`
	Novel   *Novel `gorm:"foreignKey:NovelID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Tag     *Tag   `gorm:"foreignKey:TagID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// Migration007AddNovelGenresAndTags creates the novel-genre and novel-tag many-to-many relationships
func Migration007AddNovelGenresAndTags() Migration {
	return Migration{
		ID:          "007_add_novel_genres_and_tags",
		Description: "Create novel-genre and novel-tag many-to-many relationship tables",
		Up: func(db *gorm.DB) error {
			return db.AutoMigrate(&NovelGenre{}, &NovelTag{})
		},
		Down: func(db *gorm.DB) error {
			return db.Migrator().DropTable(&NovelGenre{}, &NovelTag{})
		},
	}
}
