package dto

import novelDto "github.com/FeisalDy/nogo/internal/novel/dto"

type NovelWithDetailsDTO struct {
	// Novel data
	ID               uint    `json:"id"`
	OriginalLanguage string  `json:"original_language"`
	OriginalAuthor   *string `json:"original_author"`
	Status           *string `json:"status"`
	Source           *string `json:"source"`
	WordCount        *int    `json:"word_count"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`

	// Cross-domain data (populated by Application layer)
	Creator    *UserBasicDTO  `json:"creator,omitempty"`     // From User domain
	CoverMedia *MediaBasicDTO `json:"cover_media,omitempty"` // From Media domain (when implemented)
}

// NovelTranslationWithDetailsDTO - Translation with cross-domain data
type NovelTranslationWithDetailsDTO struct {
	// Translation data
	ID        uint    `json:"id"`
	NovelId   uint    `json:"novel_id"`
	Language  string  `json:"language"`
	Title     string  `json:"title"`
	Synopsis  *string `json:"synopsis"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`

	// Cross-domain data
	Translator *UserBasicDTO `json:"translator,omitempty"` // From User domain
}

// NovelCompleteDTO - Complete novel data with translations and related entities
type NovelCompleteDTO struct {
	Novel        NovelWithDetailsDTO              `json:"novel"`
	Translations []NovelTranslationWithDetailsDTO `json:"translations"`
}

// UserBasicDTO - Basic user info for cross-domain responses
// This is a simplified version to avoid circular dependencies
type UserBasicDTO struct {
	ID       uint    `json:"id"`
	Username *string `json:"username,omitempty"`
	Email    string  `json:"email"`
}

// MediaBasicDTO - Basic media info for cross-domain responses
// This will be implemented when Media domain is created
type MediaBasicDTO struct {
	ID       uint   `json:"id"`
	URL      string `json:"url"`
	Type     string `json:"type"`
	FileName string `json:"file_name"`
}

// CreateNovelWithCreatorDTO - Request to create novel with creator info
type CreateNovelWithCreatorDTO struct {
	novelDto.CreateNovelDTO
	// CreatedBy will be extracted from JWT token in handler
}
