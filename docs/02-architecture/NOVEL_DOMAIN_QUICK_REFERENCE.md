# Novel Domain Quick Reference

## 🎯 Golden Rule

```
Domain Layer  → Store ONLY IDs
Application Layer → Fetch full objects
```

## 📁 File Structure

```
internal/
├── novel/                          # Novel Domain (Pure)
│   ├── model/
│   │   └── novel.go               # ✅ IDs only: CreatedBy, CoverMediaId
│   ├── repository/
│   │   └── novel_repository.go    # ✅ No joins with other domains
│   ├── service/
│   │   └── novel_service.go       # ✅ Only NovelRepository dependency
│   ├── dto/
│   │   └── novel_dto.go           # ✅ DTOs with IDs only
│   └── handler/
│       └── novel_handler.go       # Domain-specific endpoints
│
└── application/                    # Application Layer (Coordinator)
    ├── dto/
    │   └── novel_dto.go           # ✅ DTOs with full objects
    ├── service/
    │   └── novel_management_service.go  # ✅ Multi-domain coordination
    └── handler/
        └── novel_handler.go       # Cross-domain endpoints
```

## 🔄 Data Flow

### Single Domain Operation

```
GET /api/v1/novels/1

Handler → NovelService → NovelRepository → Database
                             ↓
                        Returns Novel with IDs

Response:
{
  "id": 1,
  "created_by": 5,       ← Just ID
  "cover_media_id": 10   ← Just ID
}
```

### Cross-Domain Operation

```
GET /api/v1/novels/1/details

Handler → NovelManagementService
              ↓
         ┌────┴────┬────────┐
         ↓         ↓        ↓
    NovelService UserRepo MediaRepo
         ↓
    Combines all data

Response:
{
  "id": 1,
  "creator": {           ← Full object from User domain
    "id": 5,
    "username": "john",
    "email": "john@example.com"
  },
  "cover_media": {       ← Full object from Media domain
    "id": 10,
    "url": "/uploads/cover.jpg"
  }
}
```

## 💻 Code Templates

### Domain Model (✅ IDs Only)

```go
type Novel struct {
    gorm.Model
    Title        string
    CoverMediaId *uint  // ✅ ID only
    CreatedBy    *uint  // ✅ ID only
}
```

### Domain Service (✅ Pure)

```go
type NovelService struct {
    novelRepo *repository.NovelRepository  // ✅ Only same domain
}

func (s *NovelService) CreateNovel(dto *dto.CreateNovelDTO) (*dto.NovelDTO, error) {
    novel := &model.Novel{
        Title:        dto.Title,
        CreatedBy:    dto.CreatedBy,    // ✅ Just stores ID
        CoverMediaId: dto.CoverMediaId, // ✅ Just stores ID
    }
    return s.novelRepo.Create(novel)
}
```

### Application Service (✅ Coordinator)

```go
type NovelManagementService struct {
    novelService *novelService.NovelService
    userRepo     *userRepo.UserRepository     // ✅ Other domain
    mediaRepo    *mediaRepo.MediaRepository   // ✅ Other domain
}

func (s *NovelManagementService) GetNovelWithDetails(id uint) (*dto.NovelWithDetailsDTO, error) {
    // 1. Get novel (IDs only)
    novel, _ := s.novelService.GetNovelByID(id)
    
    // 2. Get creator (full object)
    creator, _ := s.userRepo.GetUserByID(*novel.CreatedBy)
    
    // 3. Get media (full object)
    media, _ := s.mediaRepo.GetByID(*novel.CoverMediaId)
    
    // 4. Combine
    return &dto.NovelWithDetailsDTO{
        ID:         novel.ID,
        Title:      novel.Title,
        Creator:    mapToUserDTO(creator),
        CoverMedia: mapToMediaDTO(media),
    }
}
```

## 📋 Checklist

When creating a new domain with foreign keys:

### Domain Layer
- [ ] Model stores IDs only (`*uint`)
- [ ] No imports from other domains
- [ ] Repository doesn't join with other tables
- [ ] Service only depends on same-domain repository
- [ ] DTOs only include IDs for foreign keys

### Application Layer
- [ ] Create Application DTOs with full objects
- [ ] Create Application Service with multiple repo dependencies
- [ ] Validate foreign key references exist
- [ ] Fetch related entities from other domains
- [ ] Combine data into Application DTOs

### API Endpoints
- [ ] Domain endpoints return IDs only
- [ ] Application endpoints return full objects
- [ ] Clear endpoint naming (e.g., `/details`, `/complete`)

## ⚠️ Red Flags

```go
// ❌ Domain model importing other domains
import "server/user/model"
type Novel struct {
    Creator *userModel.User  // ❌
}

// ❌ Domain repository with joins
db.Preload("Creator").First(&novel)  // ❌

// ❌ Domain service with multiple repos
type NovelService struct {
    novelRepo *NovelRepository
    userRepo  *UserRepository  // ❌
}

// ❌ Domain DTO with full objects
type NovelDTO struct {
    Creator *UserDTO  // ❌ Should be CreatedBy uint
}
```

## ✅ Green Lights

```go
// ✅ Domain model with IDs
type Novel struct {
    CreatedBy *uint  // ✅
}

// ✅ Domain repository without joins
db.First(&novel, id)  // ✅

// ✅ Domain service with single repo
type NovelService struct {
    novelRepo *NovelRepository  // ✅
}

// ✅ Application service coordinates
type NovelManagementService struct {
    novelService *NovelService  // ✅
    userRepo     *UserRepository // ✅
}
```

## 🎓 Remember

1. **Domain = Pure & Independent**
   - Store IDs
   - No cross-domain imports
   - Single responsibility

2. **Application = Coordinator**
   - Access multiple domains
   - Fetch full objects
   - Validate relationships

3. **Benefits**
   - Clean architecture
   - Easy testing
   - Future-proof for microservices

---

**When in doubt:** If it crosses domain boundaries, it belongs in the Application layer!
