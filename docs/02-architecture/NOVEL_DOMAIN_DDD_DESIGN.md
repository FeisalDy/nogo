# Novel Domain: DDD Design with Foreign Keys

**Date:** October 20, 2025  
**Status:** ✅ Implemented  
**Foreign Keys:** `CreatedBy` (User), `CoverMediaId` (Media), `TranslatorId` (User)

## 🎯 Key Principle

### **Domain Models Store ONLY IDs, NOT Full Objects**

```go
// ❌ WRONG - Breaks DDD
type Novel struct {
    gorm.Model
    CreatedBy uint
    Creator   *User  // ❌ Imports User model - breaks domain independence
}

// ✅ CORRECT - DDD Compliant
type Novel struct {
    gorm.Model
    CreatedBy *uint  // ✅ Just the ID - no User import needed
}
```

## 📐 Complete Architecture

### 1. Migration (Database Schema) - ✅ Can Have Full Relationships

**File:** `internal/database/migrations/004_create_novels.go`

```go
// ✅ Migrations can reference other models for schema generation
type Novel struct {
    gorm.Model
    OriginalLanguage string
    OriginalAuthor   *string
    Status           *string
    WordCount        *int
    
    // Full foreign key definitions for database
    CoverMediaId *uint
    CoverMedia   *Media `gorm:"foreignKey:CoverMediaId"` // ✅ OK here
    
    CreatedBy *uint
    Creator   *User `gorm:"foreignKey:CreatedBy"` // ✅ OK here
}
```

**Why OK:** Migrations run once to create schema, not used at runtime

### 2. Domain Model - ✅ ONLY IDs

**File:** `internal/novel/model/novel.go`

```go
// ✅ Pure domain model - NO imports from other domains
type Novel struct {
    gorm.Model
    OriginalLanguage string  `gorm:"not null;index"`
    OriginalAuthor   *string `gorm:"index"`
    Status           *string `gorm:"index"`
    WordCount        *int
    
    // ✅ Foreign keys as IDs ONLY
    CoverMediaId *uint `gorm:"index"`  // NO Media import
    CreatedBy    *uint `gorm:"index"`  // NO User import
}

type NovelTranslation struct {
    gorm.Model
    NovelId      uint    `gorm:"not null"`
    Language     string  `gorm:"not null"`
    Title        string  `gorm:"not null"`
    Synopsis     *string
    TranslatorId *uint `gorm:"index"`  // ✅ Just ID
}
```

### 3. Domain Repository - ✅ No Joins with Other Domains

**File:** `internal/novel/repository/novel_repository.go`

```go
type NovelRepository struct {
    db *gorm.DB
}

// ✅ Returns model with IDs only, no joins
func (r *NovelRepository) GetByID(id uint) (*model.Novel, error) {
    var novel model.Novel
    err := r.db.First(&novel, id).Error  // No Preload("Creator") ❌
    return &novel, err
}

// ✅ Queries by foreign key without joining
func (r *NovelRepository) GetByCreator(creatorID uint) ([]model.Novel, error) {
    var novels []model.Novel
    err := r.db.Where("created_by = ?", creatorID).Find(&novels).Error
    return novels, err
}
```

### 4. Domain Service - ✅ No Cross-Domain Dependencies

**File:** `internal/novel/service/novel_service.go`

```go
type NovelService struct {
    novelRepo *repository.NovelRepository  // ✅ ONLY Novel repo
}

// ✅ Returns domain DTO with IDs
func (s *NovelService) CreateNovel(dto *dto.CreateNovelDTO) (*dto.NovelDTO, error) {
    novel := &model.Novel{
        OriginalLanguage: dto.OriginalLanguage,
        CreatedBy:        dto.CreatedBy,     // Just stores ID
        CoverMediaId:     dto.CoverMediaId,  // Just stores ID
    }
    
    err := s.novelRepo.Create(novel)
    return s.toNovelDTO(novel), err
}
```

### 5. Application Service - ✅ Coordinates Cross-Domain

**File:** `internal/application/service/novel_management_service.go`

```go
type NovelManagementService struct {
    novelService *novelService.NovelService
    userRepo     *userRepo.UserRepository     // ✅ Access User domain
    mediaRepo    *mediaRepo.MediaRepository   // ✅ Access Media domain
    db           *gorm.DB
}

// ✅ Fetches related data from multiple domains
func (s *NovelManagementService) GetNovelWithDetails(id uint) (*dto.NovelWithDetailsDTO, error) {
    // 1. Get novel from Novel domain (returns IDs only)
    novel, _ := s.novelService.GetNovelByID(id)
    
    // 2. Get creator from User domain
    var creator *dto.UserBasicDTO
    if novel.CreatedBy != nil {
        user, _ := s.userRepo.GetUserByID(*novel.CreatedBy)
        creator = &dto.UserBasicDTO{
            ID:       user.ID,
            Username: user.Username,
            Email:    user.Email,
        }
    }
    
    // 3. Get media from Media domain (when implemented)
    var media *dto.MediaBasicDTO
    if novel.CoverMediaId != nil {
        m, _ := s.mediaRepo.GetByID(*novel.CoverMediaId)
        media = &dto.MediaBasicDTO{
            ID:  m.ID,
            URL: m.URL,
        }
    }
    
    // 4. Combine everything
    return &dto.NovelWithDetailsDTO{
        ID:         novel.ID,
        Title:      novel.Title,
        Creator:    creator,     // ✅ Full object
        CoverMedia: media,       // ✅ Full object
    }, nil
}

// ✅ Validates cross-domain references before creating
func (s *NovelManagementService) CreateNovelWithCreator(
    dto *novelDto.CreateNovelDTO,
    creatorID uint,
) (*dto.NovelWithDetailsDTO, error) {
    // 1. Validate creator exists
    _, err := s.userRepo.GetUserByID(creatorID)
    if err != nil {
        return nil, errors.New("creator not found")
    }
    
    // 2. Validate media exists (if provided)
    if dto.CoverMediaId != nil {
        _, err := s.mediaRepo.GetByID(*dto.CoverMediaId)
        if err != nil {
            return nil, errors.New("media not found")
        }
    }
    
    // 3. Create through domain service
    dto.CreatedBy = &creatorID
    novel, err := s.novelService.CreateNovel(dto)
    if err != nil {
        return nil, err
    }
    
    // 4. Return with full details
    return s.GetNovelWithDetails(novel.ID)
}
```

## 📊 API Design

### Domain Endpoints (Novel Only)

```bash
# Returns IDs only for foreign keys
GET /api/v1/novels/:id
Response:
{
  "id": 1,
  "original_language": "en",
  "created_by": 5,        # ✅ Just ID
  "cover_media_id": 10    # ✅ Just ID
}
```

### Application Endpoints (With Related Data)

```bash
# Returns full related objects
GET /api/v1/novels/:id/details
Response:
{
  "id": 1,
  "original_language": "en",
  "creator": {            # ✅ Full object
    "id": 5,
    "username": "john",
    "email": "john@example.com"
  },
  "cover_media": {        # ✅ Full object
    "id": 10,
    "url": "/uploads/cover.jpg",
    "type": "image/jpeg"
  }
}
```

## ✅ Benefits

1. **Domain Independence**
   - Novel domain doesn't import User or Media
   - Can test Novel without User/Media
   - Can change User/Media without affecting Novel

2. **Clear Boundaries**
   - Domain layer = Pure business logic with IDs
   - Application layer = Cross-domain coordination

3. **Flexibility**
   - Easy to split into microservices later
   - Easy to replace databases or ORMs
   - Easy to add new domains

4. **Maintainability**
   - Clear responsibilities
   - Easy to understand
   - Easy to test

## 🚫 Common Mistakes

### ❌ Mistake 1: Preloading in Repository
```go
// ❌ WRONG
func (r *NovelRepository) GetByID(id uint) {
    db.Preload("Creator").First(&novel, id)  // ❌ Breaks DDD
}
```

### ❌ Mistake 2: Domain Service with Multiple Repos
```go
// ❌ WRONG
type NovelService struct {
    novelRepo *NovelRepository
    userRepo  *UserRepository  // ❌ Cross-domain dependency
}
```

### ❌ Mistake 3: Domain Model Importing Other Models
```go
// ❌ WRONG
import "server/user/model"

type Novel struct {
    Creator *userModel.User  // ❌ Imports User
}
```

## 📝 Summary

| Layer | Foreign Keys | Related Objects | Import Other Domains |
|-------|--------------|-----------------|---------------------|
| **Migration** | ✅ Full FK definitions | ✅ Can reference models | ✅ OK |
| **Domain Model** | ✅ IDs only (`*uint`) | ❌ No | ❌ No |
| **Domain Repository** | ✅ Query by ID | ❌ No joins | ❌ No |
| **Domain Service** | ✅ Store/return IDs | ❌ No | ❌ No |
| **Application Service** | ✅ Uses IDs | ✅ Fetches full objects | ✅ Yes |

**Bottom Line:**
- **Domain = IDs only**
- **Application = Full objects**
- **Clear separation = Clean architecture**

---

**See Also:**
- [Application Layer](APPLICATION_LAYER.md)
- [DDD Repository Compliance](DDD_COMPLIANCE_FIX.md)
- [DDD Authorization Fix](DDD_AUTHORIZATION_FIX.md)
