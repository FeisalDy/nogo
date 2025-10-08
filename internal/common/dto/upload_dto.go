package dto

import "time"

// UploadRequest represents the upload request structure
type UploadRequest struct {
	AllowedTypes []string `json:"allowed_types,omitempty"` // Optional: restrict file types
	MaxSize      int64    `json:"max_size,omitempty"`      // Optional: max file size in bytes
}

// UploadResponse represents a single file upload response
type UploadResponse struct {
	URL        string    `json:"url"`
	Filename   string    `json:"filename"`
	Size       int64     `json:"size"`
	MimeType   string    `json:"mime_type"`
	UploadedAt time.Time `json:"uploaded_at"`
}

// MultipleUploadResponse represents multiple file upload response
type MultipleUploadResponse struct {
	Files   []UploadResponse `json:"files"`
	Count   int              `json:"count"`
	Message string           `json:"message"`
}

// UploadErrorResponse represents an error response
type UploadErrorResponse struct {
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}

// DeleteFileRequest represents a file deletion request
type DeleteFileRequest struct {
	Filename string `json:"filename" binding:"required"`
}

// DeleteFileResponse represents a file deletion response
type DeleteFileResponse struct {
	Message  string `json:"message"`
	Filename string `json:"filename"`
}
