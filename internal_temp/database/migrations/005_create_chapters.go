package migrations

import (
	"gorm.io/gorm"
)

type Chapter struct {
	gorm.Model
	NovelId   uint   `json:"novel_id" gorm:"not null;uniqueIndex:idx_novel_chapter_unique"`
	Novel     *Novel `json:"novel" gorm:"foreignKey:NovelId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Number    int    `json:"number" gorm:"not null;uniqueIndex:idx_novel_chapter_unique"`
	WordCount *int   `json:"word_count"`
}

type ChapterTranslation struct {
	gorm.Model
	ChapterId    uint     `json:"chapter_id" gorm:"not null;uniqueIndex:idx_chapter_lang_unique"`
	Chapter      *Chapter `json:"chapter" gorm:"foreignKey:ChapterId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Language     string   `json:"language" gorm:"not null;uniqueIndex:idx_chapter_lang_unique"`
	Title        string   `json:"title" gorm:"not null"`
	Content      string   `json:"content" gorm:"type:text"`
	TranslatorId *uint    `json:"translator_id"`
	Translator   *User    `json:"translator" gorm:"foreignKey:TranslatorId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

func Migration005CreateChapters() Migration {
	return Migration{
		ID:          "005_create_chapters",
		Description: "Create chapters table with novel relationship",
		Up: func(db *gorm.DB) error {
			return db.AutoMigrate(&Chapter{}, &ChapterTranslation{})
		},
		Down: func(db *gorm.DB) error {
			return db.Migrator().DropTable(&ChapterTranslation{}, &Chapter{})
		},
	}
}
