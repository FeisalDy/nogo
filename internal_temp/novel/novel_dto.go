package novel

type CreateNovelDTO struct {
	OriginalLanguage string  `json:"original_language" binding:"required"`
	OriginalAuthor   *string `json:"original_author"`
	Status           *string `json:"status"`
	Source           *string `json:"source"`
	WordCount        *int    `json:"word_count"`
	CoverMediaId     *uint   `json:"cover_media_id"`
	CreatedBy        *uint   `json:"created_by"`
}

type UpdateNovelDTO struct {
	OriginalLanguage *string `json:"original_language"`
	OriginalAuthor   *string `json:"original_author"`
	Status           *string `json:"status"`
	Source           *string `json:"source"`
	WordCount        *int    `json:"word_count"`
	CoverMediaId     *uint   `json:"cover_media_id"`
}

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

type CreateNovelTranslationDTO struct {
	NovelId      uint    `json:"novel_id" binding:"required"`
	Language     string  `json:"language" binding:"required"`
	Title        string  `json:"title" binding:"required"`
	Synopsis     *string `json:"synopsis"`
	TranslatorId *uint   `json:"translator_id"`
}

type UpdateNovelTranslationDTO struct {
	Title        *string `json:"title"`
	Synopsis     *string `json:"synopsis"`
	TranslatorId *uint   `json:"translator_id"`
}

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

type NovelWithTranslationDTO struct {
	NovelDTO
	Language string  `json:"language"`
	Title    string  `json:"title"`
	Synopsis *string `json:"synopsis"`
}

type GetAllNovelRequestDTO struct {
	Title    *string `form:"title"`
	language *string `form:"language"`
}
