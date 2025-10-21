package utils

import (
	"github.com/FeisalDy/nogo/internal/common/dto"
	"gorm.io/gorm"
)

// PaginationBuilder helps build cursor-based pagination queries
type PaginationBuilder struct {
	cursor    *dto.Cursor
	limit     int
	sortOrder string
	sortField string // Default: "id"
}

// NewPaginationBuilder creates a new pagination builder
func NewPaginationBuilder() *PaginationBuilder {
	return &PaginationBuilder{
		limit:     20,
		sortOrder: "desc",
		sortField: "id",
	}
}

// WithRequest sets pagination parameters from request
func (b *PaginationBuilder) WithRequest(req *dto.CursorPaginationRequest) *PaginationBuilder {
	if req == nil {
		return b
	}

	req = dto.NormalizePaginationRequest(req)

	b.limit = req.Limit
	b.sortOrder = req.SortOrder

	if req.Cursor != "" {
		cursor, err := dto.DecodeCursor(req.Cursor)
		if err == nil {
			b.cursor = cursor
		}
	}

	return b
}

// WithCursor sets the cursor
func (b *PaginationBuilder) WithCursor(cursor *dto.Cursor) *PaginationBuilder {
	b.cursor = cursor
	return b
}

// WithLimit sets the limit
func (b *PaginationBuilder) WithLimit(limit int) *PaginationBuilder {
	if limit > 0 && limit <= 100 {
		b.limit = limit
	}
	return b
}

// WithSortOrder sets the sort order ("asc" or "desc")
func (b *PaginationBuilder) WithSortOrder(order string) *PaginationBuilder {
	if order == "asc" || order == "desc" {
		b.sortOrder = order
	}
	return b
}

// WithSortField sets the field to sort by (default: "id")
func (b *PaginationBuilder) WithSortField(field string) *PaginationBuilder {
	if field != "" {
		b.sortField = field
	}
	return b
}

// ApplyToQuery applies pagination to a GORM query
func (b *PaginationBuilder) ApplyToQuery(query *gorm.DB) *gorm.DB {
	// Apply cursor filtering
	if b.cursor != nil && b.cursor.ID > 0 {
		if b.sortOrder == "desc" {
			query = query.Where(b.sortField+" < ?", b.cursor.ID)
		} else {
			query = query.Where(b.sortField+" > ?", b.cursor.ID)
		}
	}

	// Apply sorting
	orderClause := b.sortField + " " + b.sortOrder
	query = query.Order(orderClause)

	// Apply limit (+1 to check if there's a next page)
	query = query.Limit(b.limit + 1)

	return query
}

// BuildResponse builds a paginated response from results
func (b *PaginationBuilder) BuildResponse(results interface{}, getID func(int) uint) *dto.CursorPageInfo {
	// Use reflection or type assertion to get slice length
	// For now, we'll use a simpler approach with a callback

	pageInfo := &dto.CursorPageInfo{
		HasPreviousPage: b.cursor != nil && b.cursor.ID > 0,
	}

	return pageInfo
}

// BuildPageInfo builds page info from query results
// itemsCount: number of items returned from query
// limit: the requested limit
// getFirstID: function to get ID of first item
// getLastID: function to get ID of last item
func BuildPageInfo(itemsCount, limit int, getFirstID, getLastID func() uint) (dto.CursorPageInfo, error) {
	pageInfo := dto.CursorPageInfo{
		HasNextPage: itemsCount > limit,
	}

	if itemsCount == 0 {
		return pageInfo, nil
	}

	// Set start cursor (first item)
	if firstID := getFirstID(); firstID > 0 {
		startCursor, err := dto.EncodeCursorSimple(firstID)
		if err != nil {
			return pageInfo, err
		}
		pageInfo.StartCursor = startCursor
	}

	// Set end cursor (last item before the extra one)
	actualCount := itemsCount
	if pageInfo.HasNextPage {
		actualCount = limit // Don't count the extra item
	}

	if actualCount > 0 {
		lastID := getLastID()
		if lastID > 0 {
			endCursor, err := dto.EncodeCursorSimple(lastID)
			if err != nil {
				return pageInfo, err
			}
			pageInfo.EndCursor = endCursor
			pageInfo.NextCursor = endCursor // For next page, use end cursor
		}
	}

	return pageInfo, nil
}

// PaginateQuery is a generic helper for cursor-based pagination
// T: the model type
// query: the base GORM query
// req: pagination request
// returns: items (up to limit), pageInfo, error
func PaginateQuery[T any](query *gorm.DB, req *dto.CursorPaginationRequest) ([]T, dto.CursorPageInfo, error) {
	req = dto.NormalizePaginationRequest(req)

	// Decode cursor
	var cursorID uint
	if req.Cursor != "" {
		var err error
		cursorID, err = dto.DecodeCursorID(req.Cursor)
		if err != nil {
			return nil, dto.CursorPageInfo{}, err
		}
	}

	// Apply cursor filter
	if cursorID > 0 {
		if req.SortOrder == "desc" {
			query = query.Where("id < ?", cursorID)
		} else {
			query = query.Where("id > ?", cursorID)
		}
	}

	// Apply sorting and limit (+1 to check for next page)
	query = query.Order("id " + req.SortOrder).Limit(req.Limit + 1)

	// Execute query
	var items []T
	if err := query.Find(&items).Error; err != nil {
		return nil, dto.CursorPageInfo{}, err
	}

	// Build page info
	hasNextPage := len(items) > req.Limit
	if hasNextPage {
		items = items[:req.Limit] // Remove the extra item
	}

	pageInfo := dto.CursorPageInfo{
		HasNextPage:     hasNextPage,
		HasPreviousPage: cursorID > 0,
	}

	// Set cursors if we have items
	if len(items) > 0 {
		// We need a way to get ID from generic type T
		// This will be handled by the caller or we need interface
		// For now, return the pageInfo structure
	}

	return items, pageInfo, nil
}

// GetIDFromItem is a helper interface for pagination
type IDGetter interface {
	GetID() uint
}

// PaginateWithIDGetter paginates items that implement IDGetter interface
func PaginateWithIDGetter[T IDGetter](query *gorm.DB, req *dto.CursorPaginationRequest) ([]T, dto.CursorPageInfo, error) {
	items, pageInfo, err := PaginateQuery[T](query, req)
	if err != nil {
		return nil, dto.CursorPageInfo{}, err
	}

	// Set cursors
	if len(items) > 0 {
		firstID := items[0].GetID()
		lastID := items[len(items)-1].GetID()

		if startCursor, err := dto.EncodeCursorSimple(firstID); err == nil {
			pageInfo.StartCursor = startCursor
		}

		if endCursor, err := dto.EncodeCursorSimple(lastID); err == nil {
			pageInfo.EndCursor = endCursor
			if pageInfo.HasNextPage {
				pageInfo.NextCursor = endCursor
			}
		}
	}

	return items, pageInfo, nil
}

// SimplePaginateResponse builds a simple paginated response
func SimplePaginateResponse[T any](
	items []T,
	pageInfo dto.CursorPageInfo,
	req *dto.CursorPaginationRequest,
) dto.CursorPaginationResponse[T] {
	req = dto.NormalizePaginationRequest(req)

	return dto.CursorPaginationResponse[T]{
		Data:     items,
		PageInfo: pageInfo,
		Metadata: dto.PaginationMetadata{
			Count:     len(items),
			Limit:     req.Limit,
			SortOrder: req.SortOrder,
		},
	}
}

// OffsetPaginate handles offset-based pagination (for backward compatibility)
func OffsetPaginate[T any](query *gorm.DB, req *dto.OffsetPaginationRequest) ([]T, dto.OffsetPaginationInfo, error) {
	req = dto.NormalizeOffsetPaginationRequest(req)

	// Count total items
	var totalItems int64
	if err := query.Count(&totalItems).Error; err != nil {
		return nil, dto.OffsetPaginationInfo{}, err
	}

	// Apply pagination
	offset := req.CalculateOffset()
	var items []T
	if err := query.Offset(offset).Limit(req.Limit).Find(&items).Error; err != nil {
		return nil, dto.OffsetPaginationInfo{}, err
	}

	// Build pagination info
	pageInfo := dto.BuildOffsetPaginationInfo(req.Page, req.Limit, totalItems)

	return items, pageInfo, nil
}
