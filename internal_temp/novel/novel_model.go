package novel

import "gorm.io/gorm"

type Novel struct {
	gorm.Model
	OriginalLanguage string  `json:"original_language" gorm:"not null;index"`
	OriginalAuthor   *string `json:"original_author" gorm:"index"`
	Status           *string `json:"status" gorm:"index"`
	Source           *string `json:"source"`
	WordCount        *int    `json:"word_count"`

	CoverMediaId *uint              `json:"cover_media_id" gorm:"index"`
	CreatedBy    *uint              `json:"created_by" gorm:"index"`
	Translations []NovelTranslation `gorm:"foreignKey:NovelId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// NovelTranslation represents translations of a novel in different languages
type NovelTranslation struct {
	gorm.Model
	NovelId      uint    `json:"novel_id" gorm:"not null;uniqueIndex:idx_novel_lang_unique"`
	Language     string  `json:"language" gorm:"not null;uniqueIndex:idx_novel_lang_unique"`
	Title        string  `json:"title" gorm:"not null"`
	Synopsis     *string `json:"synopsis" gorm:"type:text"`
	TranslatorId *uint   `json:"translator_id" gorm:"index"`
}

// TableName specifies the table name for Novel
func (Novel) TableName() string {
	return "novels"
}

// GetID implements IDGetter interface for pagination
func (n Novel) GetID() uint {
	return n.ID
}

// TableName specifies the table name for NovelTranslation
func (NovelTranslation) TableName() string {
	return "novel_translations"
}

// GetID implements IDGetter interface for pagination
func (nt NovelTranslation) GetID() uint {
	return nt.ID
}
