package model

import "gorm.io/gorm"

type Chapter struct {
	gorm.Model
	NovelId   uint `json:"novel_id" gorm:"not null;uniqueIndex:idx_novel_chapter_unique"`
	Number    int  `json:"number" gorm:"not null;uniqueIndex:idx_novel_chapter_unique"`
	WordCount *int `json:"word_count"`
}

func (c Chapter) GetID() uint {
	return c.ID
}

type ChapterTranslation struct {
	gorm.Model
	ChapterId    uint   `json:"chapter_id" gorm:"not null;uniqueIndex:idx_chapter_lang_unique"`
	Language     string `json:"language" gorm:"not null;uniqueIndex:idx_chapter_lang_unique"`
	Title        string `json:"title" gorm:"not null"`
	Content      string `json:"content" gorm:"type:text"`
	TranslatorId *uint  `json:"translator_id"`
}

func (ct ChapterTranslation) GetID() uint {
	return ct.ID
}
