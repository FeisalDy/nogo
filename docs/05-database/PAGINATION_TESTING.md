# Testing Cursor Pagination

## 🧪 API Testing Examples

### Using cURL

```bash
# 1. Get first page (20 novels, newest first)
curl -X GET "http://localhost:8080/api/v1/novels?limit=20&sort_order=desc"

# Expected response:
{
  "success": true,
  "data": [
    {"id": 100, "original_language": "en", ...},
    {"id": 99, "original_language": "en", ...},
    ...
  ],
  "page_info": {
    "has_next_page": true,
    "next_cursor": "eyJpZCI6ODF9"  ← Save this!
  },
  "metadata": {
    "count": 20,
    "limit": 20,
    "sort_order": "desc"
  }
}

# 2. Get next page using cursor
curl -X GET "http://localhost:8080/api/v1/novels?limit=20&sort_order=desc&cursor=eyJpZCI6ODF9"

# 3. Get novels by creator with pagination
curl -X GET "http://localhost:8080/api/v1/novels/creator/5?limit=10&sort_order=desc"

# 4. Get ongoing novels with pagination
curl -X GET "http://localhost:8080/api/v1/novels/status/ongoing?limit=30&sort_order=asc"

# 5. Small page size for testing
curl -X GET "http://localhost:8080/api/v1/novels?limit=5&sort_order=desc"
```

### Using Postman/Thunder Client

```
Method: GET
URL: http://localhost:8080/api/v1/novels

Query Params:
┌────────────┬──────────────────────┐
│ Key        │ Value                │
├────────────┼──────────────────────┤
│ limit      │ 20                   │
│ sort_order │ desc                 │
│ cursor     │ eyJpZCI6ODF9         │ ← From previous response
└────────────┴──────────────────────┘
```

### Using JavaScript (Frontend)

```javascript
// Pagination hook
async function fetchNovels(cursor = null, limit = 20) {
  const params = new URLSearchParams({
    limit: limit.toString(),
    sort_order: 'desc'
  });
  
  if (cursor) {
    params.append('cursor', cursor);
  }
  
  const response = await fetch(
    `http://localhost:8080/api/v1/novels?${params}`
  );
  
  return await response.json();
}

// Usage in React
function NovelList() {
  const [novels, setNovels] = useState([]);
  const [nextCursor, setNextCursor] = useState(null);
  const [hasMore, setHasMore] = useState(true);
  
  const loadMore = async () => {
    const result = await fetchNovels(nextCursor);
    
    // Data is now directly in result.data (not nested)
    setNovels(prev => [...prev, ...result.data]);
    setNextCursor(result.page_info.next_cursor);
    setHasMore(result.page_info.has_next_page);
  };
  
  useEffect(() => {
    loadMore();
  }, []);
  
  return (
    <div>
      {novels.map(novel => (
        <NovelCard key={novel.id} novel={novel} />
      ))}
      
      {hasMore && (
        <button onClick={loadMore}>Load More</button>
      )}
    </div>
  );
}
```

## 🧪 Manual Testing Checklist

### Basic Functionality
- [ ] GET /novels?limit=20 returns 20 items
- [ ] Response includes `page_info` with cursors
- [ ] `has_next_page` is true when more items exist
- [ ] `next_cursor` is base64-encoded string
- [ ] Using `next_cursor` returns next page
- [ ] Items don't duplicate across pages

### Edge Cases
- [ ] First page without cursor works
- [ ] Last page has `has_next_page: false`
- [ ] Empty result set works correctly
- [ ] Invalid cursor returns error
- [ ] Limit > 100 is capped at 100
- [ ] Limit < 1 defaults to 20
- [ ] Sort order "asc" works
- [ ] Sort order "desc" works (default)

### Filtered Queries
- [ ] GET /novels/creator/:id with cursor works
- [ ] GET /novels/status/:status with cursor works
- [ ] Filters + pagination work together
- [ ] No results for filter returns empty data

### Performance
- [ ] First page is fast (< 10ms)
- [ ] Page 100 is fast (< 10ms)
- [ ] Page 1000 is fast (< 10ms)
- [ ] Response time is consistent

## 🔍 Debugging

### Check Cursor Content
```bash
# Decode cursor to see what's inside
echo "eyJpZCI6ODF9" | base64 -d
# Output: {"id":81}
```

### Check Database Query
```sql
-- This is what cursor pagination generates
SELECT * FROM novels 
WHERE id < 81 
ORDER BY id DESC 
LIMIT 21;  -- +1 to check if more exist
```

### Common Issues

#### Issue: "cursor not found" error
```bash
# Bad cursor format
cursor=invalid123

# Solution: Use cursor from API response
cursor=eyJpZCI6ODF9
```

#### Issue: Duplicate items across pages
```bash
# This means you're using offset pagination
GET /novels?page=1&limit=20  ❌

# Use cursor pagination instead
GET /novels?cursor=...&limit=20  ✅
```

#### Issue: Empty page_info
```bash
# Check if service returns pageInfo
novels, pageInfo, err := service.GetAllNovelsWithCursor(req)

# Make sure to pass pageInfo to response
response := dto.CursorPaginationResponse{
    Data: novels,
    PageInfo: pageInfo,  ← Don't forget this!
}
```

## 📊 Performance Testing

### Load Testing with Apache Bench

```bash
# Test first page
ab -n 1000 -c 10 http://localhost:8080/api/v1/novels?limit=20

# Test with cursor (simulate pagination)
ab -n 1000 -c 10 "http://localhost:8080/api/v1/novels?limit=20&cursor=eyJpZCI6ODF9"

# Compare results
# Look for: Time per request (should be consistent)
```

### Simulate Heavy Load

```bash
# Create test data
for i in {1..100000}; do
  curl -X POST http://localhost:8080/api/v1/novels \
    -H "Content-Type: application/json" \
    -d "{\"original_language\":\"en\",\"title\":\"Novel $i\"}"
done

# Test pagination performance at different positions
time curl "http://localhost:8080/api/v1/novels?limit=20"  # Page 1
time curl "http://localhost:8080/api/v1/novels?limit=20&cursor=eyJpZCI6NTAwMDB9"  # Middle
time curl "http://localhost:8080/api/v1/novels?limit=20&cursor=eyJpZCI6MTAwfQ=="  # End

# All should have similar response times!
```

## 🎯 Expected Results

### Successful Response
```json
{
  "success": true,
  "data": {
    "data": [...],  // Array of novels
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
}
```

### Error Response (Invalid Cursor)
```json
{
  "success": false,
  "error": {
    "code": "INVALID_PARAM",
    "message": "Invalid cursor format"
  }
}
```

### Empty Result
```json
{
  "success": true,
  "data": {
    "data": [],  // Empty array
    "page_info": {
      "has_next_page": false,
      "has_previous_page": false
    },
    "metadata": {
      "count": 0,
      "limit": 20,
      "sort_order": "desc"
    }
  }
}
```

## 🐛 Troubleshooting

### Cursor doesn't work across restarts
**Cause:** Cursor contains IDs that may have changed  
**Solution:** This is expected. Cursors are session-independent but data-dependent.

### Performance degrades over time
**Check:**
```sql
-- Ensure index exists on id column
SHOW INDEX FROM novels;

-- If not, create it
CREATE INDEX idx_novels_id ON novels(id);
```

### Items appear out of order
**Check sort_order parameter:**
```bash
# Wrong
?sort_order=descending  ❌

# Correct
?sort_order=desc  ✅
```

## ✅ Test Checklist

- [ ] Basic pagination works (first page)
- [ ] Navigation works (next page with cursor)
- [ ] Last page detected correctly
- [ ] Empty results handled
- [ ] Invalid cursors rejected
- [ ] Limits enforced (max 100)
- [ ] Sort order works (asc/desc)
- [ ] Filters work with pagination
- [ ] Performance is consistent
- [ ] No duplicate items
- [ ] Cursors are base64-encoded
- [ ] Metadata is correct

---

**Pro Tip:** Test with a small limit (e.g., 5) to easily verify pagination works correctly!
