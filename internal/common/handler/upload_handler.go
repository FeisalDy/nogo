package handler

import (
	"net/http"
	"strconv"

	"github.com/FeisalDy/nogo/internal/common/dto"
	"github.com/FeisalDy/nogo/internal/common/service"

	"github.com/gin-gonic/gin"
)

// UploadHandler handles file upload HTTP requests
type UploadHandler struct {
	uploadService *service.UploadService
	config        UploadHandlerConfig
}

// UploadHandlerConfig contains configuration for the upload handler
type UploadHandlerConfig struct {
	DefaultAllowedTypes []string // Default allowed MIME types
	DefaultMaxSize      int64    // Default max file size in bytes
}

// NewUploadHandler creates a new upload handler
func NewUploadHandler(uploadService *service.UploadService, config UploadHandlerConfig) *UploadHandler {
	return &UploadHandler{
		uploadService: uploadService,
		config:        config,
	}
}

// UploadSingleFile handles single file upload
func (h *UploadHandler) UploadSingleFile(c *gin.Context) {
	// Get the file from form data
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.UploadErrorResponse{
			Error: "No file provided or invalid file",
			Code:  "INVALID_FILE",
		})
		return
	}

	// Get upload parameters
	allowedTypes := h.getAllowedTypes(c)
	maxSize := h.getMaxSize(c)

	// Upload the file
	result, err := h.uploadService.UploadFile(file, allowedTypes, maxSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.UploadErrorResponse{
			Error: err.Error(),
			Code:  "UPLOAD_FAILED",
		})
		return
	}

	// Convert to DTO
	response := dto.UploadResponse{
		URL:        result.URL,
		Filename:   result.Filename,
		Size:       result.Size,
		MimeType:   result.MimeType,
		UploadedAt: result.UploadedAt,
	}

	c.JSON(http.StatusOK, response)
}

// UploadMultipleFiles handles multiple file upload
func (h *UploadHandler) UploadMultipleFiles(c *gin.Context) {
	// Parse multipart form
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.UploadErrorResponse{
			Error: "Invalid multipart form",
			Code:  "INVALID_FORM",
		})
		return
	}

	// Get files from form
	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, dto.UploadErrorResponse{
			Error: "No files provided",
			Code:  "NO_FILES",
		})
		return
	}

	// Get upload parameters
	allowedTypes := h.getAllowedTypes(c)
	maxSize := h.getMaxSize(c)

	// Upload files
	results, err := h.uploadService.UploadMultipleFiles(files, allowedTypes, maxSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.UploadErrorResponse{
			Error: err.Error(),
			Code:  "UPLOAD_FAILED",
		})
		return
	}

	// Convert to DTOs
	responseFiles := make([]dto.UploadResponse, len(results))
	for i, result := range results {
		responseFiles[i] = dto.UploadResponse{
			URL:        result.URL,
			Filename:   result.Filename,
			Size:       result.Size,
			MimeType:   result.MimeType,
			UploadedAt: result.UploadedAt,
		}
	}

	response := dto.MultipleUploadResponse{
		Files:   responseFiles,
		Count:   len(responseFiles),
		Message: "Files uploaded successfully",
	}

	c.JSON(http.StatusOK, response)
}

// DeleteFile handles file deletion
func (h *UploadHandler) DeleteFile(c *gin.Context) {
	var req dto.DeleteFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.UploadErrorResponse{
			Error: "Invalid request body",
			Code:  "INVALID_REQUEST",
		})
		return
	}

	// Delete the file
	if err := h.uploadService.DeleteFile(req.Filename); err != nil {
		c.JSON(http.StatusBadRequest, dto.UploadErrorResponse{
			Error: err.Error(),
			Code:  "DELETE_FAILED",
		})
		return
	}

	response := dto.DeleteFileResponse{
		Message:  "File deleted successfully",
		Filename: req.Filename,
	}

	c.JSON(http.StatusOK, response)
}

// getAllowedTypes gets allowed types from query params or uses defaults
func (h *UploadHandler) getAllowedTypes(c *gin.Context) []string {
	allowedTypesParam := c.QueryArray("allowed_types")
	if len(allowedTypesParam) > 0 {
		return allowedTypesParam
	}
	return h.config.DefaultAllowedTypes
}

// getMaxSize gets max size from query params or uses default
func (h *UploadHandler) getMaxSize(c *gin.Context) int64 {
	maxSizeParam := c.Query("max_size")
	if maxSizeParam != "" {
		if size, err := strconv.ParseInt(maxSizeParam, 10, 64); err == nil {
			return size
		}
	}
	return h.config.DefaultMaxSize
}
