package dto

// CreateNovelDTO - Input for creating a novel (domain layer)
type CreateNovelDTO struct {
	OriginalLanguage string  `json:"original_language" binding:"required"`
	OriginalAuthor   *string `json:"original_author"`
	Status           *string `json:"status"`
	Source           *string `json:"source"`
	WordCount        *int    `json:"word_count"`
	CoverMediaId     *uint   `json:"cover_media_id"` // Just the ID, not the full Media object
	CreatedBy        *uint   `json:"created_by"`     // Just the ID, not the full User object
}

// UpdateNovelDTO - Input for updating a novel
type UpdateNovelDTO struct {
	OriginalLanguage *string `json:"original_language"`
	OriginalAuthor   *string `json:"original_author"`
	Status           *string `json:"status"`
	Source           *string `json:"source"`
	WordCount        *int    `json:"word_count"`
	CoverMediaId     *uint   `json:"cover_media_id"`
}

// NovelDTO - Simple novel response (domain layer)
// Does NOT include related entities from other domains
type NovelDTO struct {
	ID               uint    `json:"id"`
	OriginalLanguage string  `json:"original_language"`
	OriginalAuthor   *string `json:"original_author"`
	Status           *string `json:"status"`
	Source           *string `json:"source"`
	WordCount        *int    `json:"word_count"`
	CoverMediaId     *uint   `json:"cover_media_id"`
	CreatedBy        *uint   `json:"created_by"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
}

// CreateNovelTranslationDTO - Input for creating a translation
type CreateNovelTranslationDTO struct {
	NovelId      uint    `json:"novel_id" binding:"required"`
	Language     string  `json:"language" binding:"required"`
	Title        string  `json:"title" binding:"required"`
	Synopsis     *string `json:"synopsis"`
	TranslatorId *uint   `json:"translator_id"`
}

// UpdateNovelTranslationDTO - Input for updating a translation
type UpdateNovelTranslationDTO struct {
	Title        *string `json:"title"`
	Synopsis     *string `json:"synopsis"`
	TranslatorId *uint   `json:"translator_id"`
}

// NovelTranslationDTO - Translation response (domain layer)
type NovelTranslationDTO struct {
	ID           uint    `json:"id"`
	NovelId      uint    `json:"novel_id"`
	Language     string  `json:"language"`
	Title        string  `json:"title"`
	Synopsis     *string `json:"synopsis"`
	TranslatorId *uint   `json:"translator_id"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}
