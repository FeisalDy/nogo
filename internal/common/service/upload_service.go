package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// UploadService handles file upload operations
type UploadService struct {
	uploadDir string
	baseURL   string
}

// UploadResult contains the result of a file upload
type UploadResult struct {
	URL        string    `json:"url"`
	Filename   string    `json:"filename"`
	Size       int64     `json:"size"`
	MimeType   string    `json:"mime_type"`
	UploadedAt time.Time `json:"uploaded_at"`
}

// UploadConfig contains configuration for the upload service
type UploadConfig struct {
	UploadDir    string   // Directory to store uploaded files
	BaseURL      string   // Base URL for accessing files
	AllowedTypes []string // Allowed MIME types (e.g., ["image/jpeg", "image/png"])
	MaxSize      int64    // Maximum file size in bytes
}

// NewUploadService creates a new upload service
func NewUploadService(config UploadConfig) *UploadService {
	// Create upload directory if it doesn't exist
	if err := os.MkdirAll(config.UploadDir, 0755); err != nil {
		panic(fmt.Sprintf("failed to create upload directory: %v", err))
	}

	return &UploadService{
		uploadDir: config.UploadDir,
		baseURL:   config.BaseURL,
	}
}

// UploadFile uploads a file and returns its URL and metadata
func (s *UploadService) UploadFile(file *multipart.FileHeader, allowedTypes []string, maxSize int64) (*UploadResult, error) {
	// Validate file size
	if maxSize > 0 && file.Size > maxSize {
		return nil, fmt.Errorf("file size %d bytes exceeds maximum allowed size %d bytes", file.Size, maxSize)
	}

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	// Validate file type if restrictions are set
	if len(allowedTypes) > 0 {
		contentType := file.Header.Get("Content-Type")
		if !s.isAllowedType(contentType, allowedTypes) {
			return nil, fmt.Errorf("file type %s is not allowed", contentType)
		}
	}

	// Generate unique filename
	filename := s.generateUniqueFilename(file.Filename)
	filePath := filepath.Join(s.uploadDir, filename)

	// Create the destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	// Copy the uploaded file to destination
	size, err := io.Copy(dst, src)
	if err != nil {
		// Clean up the partially created file
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	// Generate the public URL
	url := s.generateFileURL(filename)

	return &UploadResult{
		URL:        url,
		Filename:   filename,
		Size:       size,
		MimeType:   file.Header.Get("Content-Type"),
		UploadedAt: time.Now(),
	}, nil
}

// UploadMultipleFiles uploads multiple files
func (s *UploadService) UploadMultipleFiles(files []*multipart.FileHeader, allowedTypes []string, maxSize int64) ([]*UploadResult, error) {
	results := make([]*UploadResult, 0, len(files))

	for _, file := range files {
		result, err := s.UploadFile(file, allowedTypes, maxSize)
		if err != nil {
			return nil, fmt.Errorf("failed to upload file %s: %w", file.Filename, err)
		}
		results = append(results, result)
	}

	return results, nil
}

// DeleteFile deletes a file by filename
func (s *UploadService) DeleteFile(filename string) error {
	filePath := filepath.Join(s.uploadDir, filename)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file %s does not exist", filename)
	}

	// Delete the file
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete file %s: %w", filename, err)
	}

	return nil
}

// generateUniqueFilename generates a unique filename using timestamp
func (s *UploadService) generateUniqueFilename(originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	nameWithoutExt := strings.TrimSuffix(originalFilename, ext)

	// Clean the filename (remove spaces and special characters)
	nameWithoutExt = strings.ReplaceAll(nameWithoutExt, " ", "_")

	timestamp := time.Now().Unix()
	return fmt.Sprintf("%s_%d%s", nameWithoutExt, timestamp, ext)
}

// generateFileURL generates the public URL for accessing the file
func (s *UploadService) generateFileURL(filename string) string {
	return fmt.Sprintf("%s/uploads/%s", strings.TrimRight(s.baseURL, "/"), filename)
}

// isAllowedType checks if the content type is allowed
func (s *UploadService) isAllowedType(contentType string, allowedTypes []string) bool {
	for _, allowed := range allowedTypes {
		if contentType == allowed {
			return true
		}
	}
	return false
}

// GetUploadDir returns the upload directory path for serving static files
func (s *UploadService) GetUploadDir() string {
	return s.uploadDir
}
