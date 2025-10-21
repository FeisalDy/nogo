# Cursor-Based Pagination System

**Date:** October 20, 2025  
**Status:** ✅ Implemented  
**Location:** `internal/common/dto/pagination_dto.go`, `internal/common/utils/pagination.go`

## 🎯 Why Cursor-Based Pagination?

### Problem with Offset Pagination
```sql
-- Page 1: OFFSET 0 LIMIT 20 (Fast)
SELECT * FROM novels OFFSET 0 LIMIT 20;

-- Page 1000: OFFSET 19980 LIMIT 20 (VERY SLOW!)
SELECT * FROM novels OFFSET 19980 LIMIT 20;
-- Database must scan 19,980 rows to skip them
```

**With 100,000+ novels:**
- ❌ Offset pagination gets exponentially slower
- ❌ Inconsistent results if data changes between pages
- ❌ Can't efficiently jump to arbitrary positions

### Solution: Cursor-Based Pagination
```sql
-- Uses WHERE clause instead of OFFSET
SELECT * FROM novels WHERE id < 5000 ORDER BY id DESC LIMIT 20;
-- Database uses index efficiently, always fast!
```

**Benefits:**
- ✅ **Constant performance** regardless of dataset size
- ✅ **Consistent results** even with concurrent writes
- ✅ **Index-friendly** queries
- ✅ **Scalable** to millions of records

## 📐 Architecture

### Components

```
┌─────────────────────────────────────────────────────────────┐
│                 Common Pagination System                    │
│                    (Reusable Across All Domains)            │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  1. DTOs (internal/common/dto/pagination_dto.go)           │
│     ├── CursorPaginationRequest                            │
│     ├── CursorPaginationResponse                           │
│     ├── CursorPageInfo                                     │
│     └── Cursor (encoded/decoded)                           │
│                                                             │
│  2. Utilities (internal/common/utils/pagination.go)        │
│     ├── PaginateWithIDGetter[T]()                          │
│     ├── BuildPageInfo()                                    │
│     └── EncodeCursor() / DecodeCursor()                    │
│                                                             │
└─────────────────────────────────────────────────────────────┘
                          ↓
                    Used by all domains
                          ↓
┌─────────────────────────────────────────────────────────────┐
│  Novel Domain                                               │
│  ├── Model: implements GetID()                             │
│  ├── Repository: GetAllWithCursor()                        │
│  ├── Service: GetAllNovelsWithCursor()                     │
│  └── Handler: GetAllNovels() with cursor params            │
└─────────────────────────────────────────────────────────────┘
```

## 🚀 Usage Guide

### 1. Make Your Model Implement IDGetter

```go
// internal/novel/model/novel.go

// GetID implements IDGetter interface for pagination
func (n Novel) GetID() uint {
    return n.ID
}
```

**That's it!** The pagination system can now work with your model.

### 2. Add Cursor Pagination to Repository

```go
// internal/novel/repository/novel_repository.go

func (r *NovelRepository) GetAllWithCursor(
    req *dto.CursorPaginationRequest,
) ([]model.Novel, dto.CursorPageInfo, error) {
    baseQuery := r.db.Model(&model.Novel{})
    return utils.PaginateWithIDGetter[model.Novel](baseQuery, req)
}

// For filtered queries
func (r *NovelRepository) GetByStatusWithCursor(
    status string,
    req *dto.CursorPaginationRequest,
) ([]model.Novel, dto.CursorPageInfo, error) {
    baseQuery := r.db.Where("status = ?", status)
    return utils.PaginateWithIDGetter[model.Novel](baseQuery, req)
}
```

### 3. Add Methods to Service

```go
// internal/novel/service/novel_service.go

func (s *NovelService) GetAllNovelsWithCursor(
    req *commonDto.CursorPaginationRequest,
) ([]dto.NovelDTO, commonDto.CursorPageInfo, error) {
    novels, pageInfo, err := s.novelRepo.GetAllWithCursor(req)
    if err != nil {
        return nil, commonDto.CursorPageInfo{}, err
    }
    
    // Convert to DTOs
    novelDTOs := make([]dto.NovelDTO, len(novels))
    for i, novel := range novels {
        novelDTOs[i] = *s.toNovelDTO(&novel)
    }
    
    return novelDTOs, pageInfo, nil
}
```

### 4. Add Handler Endpoints

```go
// internal/novel/handler/novel_handler.go

func (h *NovelHandler) GetAllNovels(c *gin.Context) {
    var req commonDto.CursorPaginationRequest
    if err := c.ShouldBindQuery(&req); err != nil {
        utils.RespondValidationError(c, err, errors.ErrCodeInvalidParam)
        return
    }
    
    novels, pageInfo, err := h.novelService.GetAllNovelsWithCursor(&req)
    if err != nil {
        utils.HandleServiceError(c, err)
        return
    }
    
    // Clean response without nested data
    utils.RespondSuccessWithPagination(
        c,
        http.StatusOK,
        novels,   // Direct array
        pageInfo, // Pagination navigation
        commonDto.PaginationMetadata{
            Count:     len(novels),
            Limit:     req.Limit,
            SortOrder: req.SortOrder,
        },
    )
}
```

## 📡 API Usage

### Request Parameters

```bash
GET /api/v1/novels?limit=20&sort_order=desc&cursor=eyJpZCI6MTAwfQ==
```

**Query Parameters:**
- `cursor` (optional): Base64-encoded cursor from previous response
- `limit` (optional): Items per page (default: 20, max: 100)
- `sort_order` (optional): "asc" or "desc" (default: "desc")

### Response Format

```json
{
  "success": true,
  "data": [
    {
      "id": 100,
      "title": "Novel Title",
      "original_language": "en",
      "created_at": "2025-10-20T10:00:00Z"
    }
    // ... 19 more items
  ],
  "page_info": {
    "has_next_page": true,
    "has_previous_page": false,
    "next_cursor": "eyJpZCI6ODF9",           // Use this for next page
    "previous_cursor": "",
    "start_cursor": "eyJpZCI6MTAwfQ==",
    "end_cursor": "eyJpZCI6ODF9"
  },
  "metadata": {
    "count": 20,
    "limit": 20,
    "sort_order": "desc"
  }
}
```

### Pagination Flow

```bash
# 1. Get first page (no cursor)
GET /api/v1/novels?limit=20&sort_order=desc

Response:
{
  "page_info": {
    "has_next_page": true,
    "next_cursor": "eyJpZCI6ODF9"  ← Save this
  }
}

# 2. Get next page (use next_cursor)
GET /api/v1/novels?limit=20&sort_order=desc&cursor=eyJpZCI6ODF9

Response:
{
  "page_info": {
    "has_next_page": true,
    "next_cursor": "eyJpZCI6NjF9"  ← Use for page 3
  }
}

# 3. Continue until has_next_page = false
```

## 🔄 How It Works Internally

### 1. Cursor Encoding

```go
// Cursor structure
type Cursor struct {
    ID uint `json:"id"`  // The primary key
}

// Encoding example
cursor := &Cursor{ID: 100}
encoded := EncodeCursor(cursor)
// Result: "eyJpZCI6MTAwfQ==" (base64 of {"id":100})
```

### 2. Query Building

```sql
-- Without cursor (first page, descending)
SELECT * FROM novels ORDER BY id DESC LIMIT 21;
-- Returns IDs: 100, 99, 98, ..., 80 (21 items to check if more exist)

-- With cursor (next page, descending)
SELECT * FROM novels 
WHERE id < 80  -- Last ID from previous page
ORDER BY id DESC 
LIMIT 21;
-- Returns IDs: 79, 78, 77, ..., 59
```

### 3. Has Next Page Detection

```go
// Query for limit + 1 items
items := QueryWithLimit(limit + 1)  // e.g., 21 items

if len(items) > limit {
    hasNextPage = true
    items = items[:limit]  // Keep only 20 items
} else {
    hasNextPage = false
}
```

## 🎨 Complete Example: Add Pagination to New Domain

Let's say you want to add pagination to a **Chapter** domain:

### Step 1: Model
```go
// internal/chapter/model/chapter.go
type Chapter struct {
    gorm.Model
    NovelID uint
    Title   string
    Content string
}

// ✅ Implement GetID
func (c Chapter) GetID() uint {
    return c.ID
}
```

### Step 2: Repository
```go
// internal/chapter/repository/chapter_repository.go
import (
    "github.com/FeisalDy/nogo/internal/common/dto"
    "github.com/FeisalDy/nogo/internal/common/utils"
)

func (r *ChapterRepository) GetByNovelWithCursor(
    novelID uint,
    req *dto.CursorPaginationRequest,
) ([]model.Chapter, dto.CursorPageInfo, error) {
    baseQuery := r.db.Where("novel_id = ?", novelID)
    return utils.PaginateWithIDGetter[model.Chapter](baseQuery, req)
}
```

### Step 3: Service
```go
// internal/chapter/service/chapter_service.go
func (s *ChapterService) GetChaptersByNovelWithCursor(
    novelID uint,
    req *commonDto.CursorPaginationRequest,
) ([]dto.ChapterDTO, commonDto.CursorPageInfo, error) {
    chapters, pageInfo, err := s.chapterRepo.GetByNovelWithCursor(novelID, req)
    if err != nil {
        return nil, commonDto.CursorPageInfo{}, err
    }
    
    chapterDTOs := make([]dto.ChapterDTO, len(chapters))
    for i, chapter := range chapters {
        chapterDTOs[i] = *s.toChapterDTO(&chapter)
    }
    
    return chapterDTOs, pageInfo, nil
}
```

### Step 4: Handler
```go
// internal/chapter/handler/chapter_handler.go
func (h *ChapterHandler) GetChaptersByNovel(c *gin.Context) {
    novelID, _ := strconv.ParseUint(c.Param("novel_id"), 10, 32)
    
    var req commonDto.CursorPaginationRequest
    c.ShouldBindQuery(&req)
    
    chapters, pageInfo, err := h.chapterService.GetChaptersByNovelWithCursor(uint(novelID), &req)
    if err != nil {
        utils.HandleServiceError(c, err)
        return
    }
    
    utils.RespondSuccessWithPagination(
        c,
        http.StatusOK,
        chapters,
        pageInfo,
        commonDto.PaginationMetadata{
            Count:     len(chapters),
            Limit:     req.Limit,
            SortOrder: req.SortOrder,
        },
    )
}
```

### Step 5: Routes
```go
chapterRoutes.GET("/novels/:novel_id/chapters", chapterHandler.GetChaptersByNovel)
```

## 📊 Performance Comparison

### Offset Pagination (Old Way)
```
Page 1:    OFFSET 0     → 0.001s  ✅ Fast
Page 100:  OFFSET 1980  → 0.050s  ⚠️  Slower
Page 1000: OFFSET 19980 → 2.500s  ❌ Very slow
Page 5000: OFFSET 99980 → 15.00s  ❌ Extremely slow
```

### Cursor Pagination (New Way)
```
Page 1:    WHERE id < MAX      → 0.001s  ✅ Fast
Page 100:  WHERE id < 80000    → 0.001s  ✅ Fast
Page 1000: WHERE id < 1000     → 0.001s  ✅ Fast
Page 5000: WHERE id < 100      → 0.001s  ✅ Fast
```

**All pages have consistent performance!** 🚀

## ⚠️ Limitations & Considerations

### 1. Can't Jump to Arbitrary Page
```
❌ Can't do: "Go to page 500"
✅ Can do: "Next page", "Previous page"
```
**Solution:** Use offset pagination for admin interfaces where page jumping is needed.

### 2. Cursor Changes if Data Changes
If items are deleted/inserted, cursors remain valid but position may shift.
**This is expected behavior and prevents phantom reads.**

### 3. Simple Sorting Only (for now)
Current implementation sorts by ID only.
**For complex sorting:** Extend `Cursor` struct to include sort fields.

## 🔮 Future Enhancements

### Multi-Field Sorting
```go
type Cursor struct {
    ID               uint                       `json:"id"`
    AdditionalFields map[string]interface{}     `json:"fields"`
}

// Sort by created_at DESC, then id DESC
cursor := &Cursor{
    ID: 100,
    AdditionalFields: map[string]interface{}{
        "created_at": "2025-10-20T10:00:00Z",
    },
}
```

### Bidirectional Pagination
```go
// Get previous page
req := &CursorPaginationRequest{
    Cursor:    previousCursor,
    Limit:     20,
    SortOrder: "asc",  // Reverse direction
    Direction: "before",  // New field
}
```

## 📋 Summary

### ✅ What You Get

1. **Reusable pagination system** for all domains
2. **Consistent API** across all endpoints
3. **Scalable performance** for large datasets
4. **Easy to implement** (just 4 steps per domain)
5. **Type-safe** with generics

### 🎯 Best Practices

1. **Always use cursor pagination** for user-facing lists
2. **Use offset pagination** only for admin interfaces
3. **Set reasonable defaults** (limit: 20, max: 100)
4. **Add indexes** on sorted columns
5. **Document cursor format** in API docs

### 📚 Related Files

- `internal/common/dto/pagination_dto.go` - DTO definitions
- `internal/common/utils/pagination.go` - Helper functions
- `internal/novel/model/novel.go` - Example implementation
- `internal/novel/repository/novel_repository.go` - Repository example
- `internal/novel/service/novel_service.go` - Service example
- `internal/novel/handler/novel_handler.go` - Handler example

---

**Remember:** Cursor pagination scales to millions of records with constant performance! 🚀
