package migrations

import (
	"gorm.io/gorm"
)

type Novel struct {
	gorm.Model
	OriginalLanguage string  `json:"original_language" gorm:"not null;index"`
	Status           *string `json:"status" gorm:"index"`
	Source           *string `json:"source"`
	WordCount        *int    `json:"word_count"`

	CoverMediaId *uint  `json:"cover_media_id"`
	CoverMedia   *Media `json:"cover_media" gorm:"foreignKey:CoverMediaId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

	AuthorID *uint `json:"author_id" gorm:"index"`
	Author   *User `json:"author" gorm:"foreignKey:AuthorID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

type NovelTranslation struct {
	gorm.Model
	Language     string  `json:"language" gorm:"not null;uniqueIndex:idx_novel_lang_unique"`
	Title        string  `json:"title" gorm:"not null"`
	Synopsis     *string `json:"synopsis" gorm:"type:text"`
	TranslatorId *uint   `json:"translator_id"`
	Translator   *User   `json:"translator" gorm:"foreignKey:TranslatorId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	NovelId      uint    `json:"novel_id" gorm:"not null;uniqueIndex:idx_novel_lang_unique"`
	Novel        *Novel  `json:"novel" gorm:"foreignKey:NovelId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func Migration004CreateNovels() Migration {
	return Migration{
		ID:          "004_create_novels",
		Description: "Create novels table with author relationship",
		Up: func(db *gorm.DB) error {
			return db.AutoMigrate(&Novel{}, &NovelTranslation{})
		},
		Down: func(db *gorm.DB) error {
			return db.Migrator().DropTable(&NovelTranslation{}, &Novel{})
		},
	}
}
