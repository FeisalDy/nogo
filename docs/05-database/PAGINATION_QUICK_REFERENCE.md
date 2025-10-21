# Cursor Pagination Quick Reference

## 🚀 For Developers: How to Add Pagination to Any Domain

### Step 1: Model (1 line)
```go
func (m YourModel) GetID() uint { return m.ID }
```

### Step 2: Repository (3 lines)
```go
func (r *YourRepository) GetAllWithCursor(req *dto.CursorPaginationRequest) ([]model.YourModel, dto.CursorPageInfo, error) {
    return utils.PaginateWithIDGetter[model.YourModel](r.db.Model(&model.YourModel{}), req)
}
```

### Step 3: Service (10 lines)
```go
func (s *YourService) GetAllWithCursor(req *commonDto.CursorPaginationRequest) ([]dto.YourDTO, commonDto.CursorPageInfo, error) {
    items, pageInfo, err := s.yourRepo.GetAllWithCursor(req)
    if err != nil {
        return nil, commonDto.CursorPageInfo{}, err
    }
    
    dtos := make([]dto.YourDTO, len(items))
    for i, item := range items {
        dtos[i] = *s.toDTO(&item)
    }
    
    return dtos, pageInfo, nil
}
```

### Step 4: Handler (15 lines)
```go
func (h *YourHandler) GetAll(c *gin.Context) {
    var req commonDto.CursorPaginationRequest
    if err := c.ShouldBindQuery(&req); err != nil {
        utils.RespondValidationError(c, err, errors.ErrCodeInvalidParam)
        return
    }
    
    items, pageInfo, err := h.yourService.GetAllWithCursor(&req)
    if err != nil {
        utils.HandleServiceError(c, err)
        return
    }
    
    utils.RespondSuccessWithPagination(c, http.StatusOK, items, pageInfo, commonDto.PaginationMetadata{
        Count: len(items), Limit: req.Limit, SortOrder: req.SortOrder,
    })
}
```

## 📡 API Examples

### Get First Page
```bash
GET /api/v1/novels?limit=20&sort_order=desc

Response:
{
  "data": [...],
  "page_info": {
    "has_next_page": true,
    "next_cursor": "eyJpZCI6ODF9"  ← Use this for next page
  }
}

### Get Next Page
```bash
GET /api/v1/novels?limit=20&sort_order=desc&cursor=eyJpZCI6ODF9
```

### With Filters
```bash
GET /api/v1/novels?status=ongoing&limit=20&cursor=eyJpZCI6MTAwfQ==
```

## 🎯 Response Structure

```json
{
  "success": true,
  "data": [
    { "id": 100, "title": "Novel 1", ... },
    { "id": 99, "title": "Novel 2", ... }
  ],
  "page_info": {
    "has_next_page": true,
    "has_previous_page": false,
    "next_cursor": "eyJpZCI6ODF9",
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

## 📊 Query Parameters

| Parameter | Type | Default | Max | Description |
|-----------|------|---------|-----|-------------|
| `cursor` | string | "" | - | Base64 cursor from previous response |
| `limit` | int | 20 | 100 | Items per page |
| `sort_order` | string | "desc" | - | "asc" or "desc" |

## ⚡ Performance

```
Offset Pagination:
Page 1:    0.001s ✅
Page 1000: 2.500s ❌
Page 5000: 15.00s ❌

Cursor Pagination:
Page 1:    0.001s ✅
Page 1000: 0.001s ✅
Page 5000: 0.001s ✅
```

## 🔑 Key Files

```
internal/common/
├── dto/pagination_dto.go          ← DTOs (reusable)
└── utils/pagination.go            ← Helpers (reusable)

internal/your-domain/
├── model/your_model.go            ← Add GetID()
├── repository/your_repository.go  ← Add GetAllWithCursor()
├── service/your_service.go        ← Add GetAllWithCursor()
└── handler/your_handler.go        ← Add GetAll() endpoint
```

## ✅ Checklist

- [ ] Model implements `GetID() uint`
- [ ] Repository has `GetAllWithCursor()` method
- [ ] Service has `GetAllWithCursor()` method
- [ ] Handler binds `CursorPaginationRequest`
- [ ] Handler returns `CursorPaginationResponse`
- [ ] Tested with large dataset

## 💡 Tips

1. **Always use cursor pagination** for user-facing lists
2. **Default limit: 20** (good balance)
3. **Max limit: 100** (prevent abuse)
4. **Index your sort column** (usually `id`)
5. **Cursor is opaque** - don't let clients construct it

## 🚨 Common Mistakes

```go
// ❌ Don't use offset in new code
GetAll(limit, offset int)

// ✅ Use cursor pagination
GetAllWithCursor(req *dto.CursorPaginationRequest)
```

```go
// ❌ Don't expose raw IDs
"next_id": 100

// ✅ Use encoded cursors
"next_cursor": "eyJpZCI6MTAwfQ=="
```

```go
// ❌ Don't forget GetID()
type Novel struct { ... }

// ✅ Implement GetID()
func (n Novel) GetID() uint { return n.ID }
```

## 🔗 See Full Guide

[CURSOR_PAGINATION_GUIDE.md](CURSOR_PAGINATION_GUIDE.md) - Complete documentation with examples

---

**Remember:** Copy the 4 code blocks above and adapt them to your domain! 🚀
