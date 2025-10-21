package dto

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// CursorPaginationRequest represents the request parameters for cursor-based pagination
type CursorPaginationRequest struct {
	// Cursor is the base64-encoded cursor pointing to the last item of the previous page
	Cursor string `form:"cursor" json:"cursor"`

	// Limit is the maximum number of items to return (default: 20, max: 100)
	Limit int `form:"limit" json:"limit" binding:"omitempty,min=1,max=100"`

	// SortOrder is the sort direction: "asc" or "desc" (default: "desc")
	SortOrder string `form:"sort_order" json:"sort_order" binding:"omitempty,oneof=asc desc"`
}

// CursorPaginationResponse represents the response structure for cursor-based pagination
type CursorPaginationResponse[T any] struct {
	Data     []T                `json:"data"`
	PageInfo CursorPageInfo     `json:"page_info"`
	Metadata PaginationMetadata `json:"metadata,omitempty"`
}

// CursorPageInfo contains pagination navigation information
type CursorPageInfo struct {
	// HasNextPage indicates if there are more items after the current page
	HasNextPage bool `json:"has_next_page"`

	// HasPreviousPage indicates if there are items before the current page
	HasPreviousPage bool `json:"has_previous_page"`

	// NextCursor is the cursor to fetch the next page (base64-encoded)
	NextCursor string `json:"next_cursor,omitempty"`

	// PreviousCursor is the cursor to fetch the previous page (base64-encoded)
	PreviousCursor string `json:"previous_cursor,omitempty"`

	// StartCursor is the cursor of the first item in the current page
	StartCursor string `json:"start_cursor,omitempty"`

	// EndCursor is the cursor of the last item in the current page
	EndCursor string `json:"end_cursor,omitempty"`
}

// PaginationMetadata contains additional pagination information
type PaginationMetadata struct {
	// Count is the number of items in the current page
	Count int `json:"count"`

	// Limit is the requested limit
	Limit int `json:"limit"`

	// SortOrder is the sort direction used
	SortOrder string `json:"sort_order"`
}

// Cursor represents the internal cursor structure
// This is encoded/decoded to/from base64 for external use
type Cursor struct {
	// ID is the primary key of the item (used for pagination)
	ID uint `json:"id"`

	// AdditionalFields can store extra fields for complex sorting
	// For example, if sorting by created_at then ID, store created_at here
	AdditionalFields map[string]interface{} `json:"fields,omitempty"`
}

// EncodeCursor encodes a cursor to base64 string
func EncodeCursor(cursor *Cursor) (string, error) {
	if cursor == nil || cursor.ID == 0 {
		return "", nil
	}

	jsonData, err := json.Marshal(cursor)
	if err != nil {
		return "", fmt.Errorf("failed to marshal cursor: %w", err)
	}

	return base64.URLEncoding.EncodeToString(jsonData), nil
}

// DecodeCursor decodes a base64 cursor string to Cursor struct
func DecodeCursor(encoded string) (*Cursor, error) {
	if encoded == "" {
		return nil, nil
	}

	jsonData, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("failed to decode cursor: %w", err)
	}

	var cursor Cursor
	if err := json.Unmarshal(jsonData, &cursor); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cursor: %w", err)
	}

	return &cursor, nil
}

// EncodeCursorSimple creates a simple cursor from just an ID
func EncodeCursorSimple(id uint) (string, error) {
	if id == 0 {
		return "", nil
	}

	cursor := &Cursor{ID: id}
	return EncodeCursor(cursor)
}

// DecodeCursorID extracts just the ID from a cursor (simplified version)
func DecodeCursorID(encoded string) (uint, error) {
	if encoded == "" {
		return 0, nil
	}

	cursor, err := DecodeCursor(encoded)
	if err != nil {
		return 0, err
	}

	if cursor == nil {
		return 0, nil
	}

	return cursor.ID, nil
}

// GetDefaultPaginationRequest returns a pagination request with default values
func GetDefaultPaginationRequest() *CursorPaginationRequest {
	return &CursorPaginationRequest{
		Limit:     20,
		SortOrder: "desc",
	}
}

// NormalizePaginationRequest sets default values if not provided
func NormalizePaginationRequest(req *CursorPaginationRequest) *CursorPaginationRequest {
	if req == nil {
		return GetDefaultPaginationRequest()
	}

	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	if req.SortOrder != "asc" && req.SortOrder != "desc" {
		req.SortOrder = "desc"
	}

	return req
}

// OffsetPaginationRequest represents offset-based pagination (for backward compatibility)
type OffsetPaginationRequest struct {
	Page  int `form:"page" json:"page" binding:"omitempty,min=1"`
	Limit int `form:"limit" json:"limit" binding:"omitempty,min=1,max=100"`
}

// OffsetPaginationResponse represents offset-based pagination response
type OffsetPaginationResponse[T any] struct {
	Data       []T                  `json:"data"`
	Pagination OffsetPaginationInfo `json:"pagination"`
}

// OffsetPaginationInfo contains offset pagination information
type OffsetPaginationInfo struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalPages int   `json:"total_pages"`
	TotalItems int64 `json:"total_items"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

// GetDefaultOffsetPaginationRequest returns offset pagination with defaults
func GetDefaultOffsetPaginationRequest() *OffsetPaginationRequest {
	return &OffsetPaginationRequest{
		Page:  1,
		Limit: 20,
	}
}

// NormalizeOffsetPaginationRequest sets default values
func NormalizeOffsetPaginationRequest(req *OffsetPaginationRequest) *OffsetPaginationRequest {
	if req == nil {
		return GetDefaultOffsetPaginationRequest()
	}

	if req.Page <= 0 {
		req.Page = 1
	}

	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	return req
}

// CalculateOffset calculates the offset from page and limit
func (r *OffsetPaginationRequest) CalculateOffset() int {
	return (r.Page - 1) * r.Limit
}

// CalculateTotalPages calculates total pages from total items
func CalculateTotalPages(totalItems int64, limit int) int {
	if limit <= 0 {
		return 0
	}

	totalPages := int(totalItems) / limit
	if int(totalItems)%limit > 0 {
		totalPages++
	}

	return totalPages
}

// BuildOffsetPaginationInfo builds pagination info from parameters
func BuildOffsetPaginationInfo(page, limit int, totalItems int64) OffsetPaginationInfo {
	totalPages := CalculateTotalPages(totalItems, limit)

	return OffsetPaginationInfo{
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
		TotalItems: totalItems,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
}

// ConvertOffsetToCursorID converts offset pagination to approximate cursor ID
// This is a helper for migration from offset to cursor pagination
func ConvertOffsetToCursorID(page, limit int, firstID uint) uint {
	offset := (page - 1) * limit
	return firstID + uint(offset)
}

// Simple type alias for backward compatibility
type PaginatedResponse[T any] = CursorPaginationResponse[T]
type PaginationRequest = CursorPaginationRequest
