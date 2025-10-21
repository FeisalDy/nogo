package dto

type CreateChapterDTO struct {
	NovelID   uint `json:"novel_id" binding:"required"`
	Number    int  `json:"number" binding:"required"`
	WordCount *int `json:"word_count"`
}

type ChapterDTO struct {
	ID        uint   `json:"id"`
	NovelID   uint   `json:"novel_id"`
	Number    int    `json:"number"`
	WordCount *int   `json:"word_count"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type CreateChapterTranslationDTO struct {
	ChapterID    uint   `json:"chapter_id" binding:"required"`
	Language     string `json:"language" binding:"required"`
	Title        string `json:"title" binding:"required"`
	Content      string `json:"content" binding:"required"`
	TranslatorId *uint  `json:"translator_id"`
}

type ChapterTranslationDTO struct {
	ID           uint   `json:"id"`
	ChapterID    uint   `json:"chapter_id"`
	Language     string `json:"language"`
	Title        string `json:"title"`
	Content      string `json:"content"`
	TranslatorId *uint  `json:"translator_id"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type ChapterWithTranslationDTO struct {
	ChapterDTO
	Language string `json:"language"`
	Title    string `json:"title"`
	Content  string `json:"content"`
}
