package dto

type MediaDTO struct {
	ID          uint    `json:"id"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
	UploadBy    *uint   `json:"upload_by"`
	URL         string  `json:"url"`
	Type        *string `json:"type"`
	Description *string `json:"description"`
	FileSize    *int    `json:"file_size"`
	MimeType    *string `json:"mime_type"`
}

type CreateMediaDTO struct {
	UploadBy    *uint   `json:"upload_by"`
	URL         string  `json:"url" binding:"required"`
	Type        *string `json:"type"`
	Description *string `json:"description"`
	FileSize    *int    `json:"file_size"`
	MimeType    *string `json:"mime_type"`
}
