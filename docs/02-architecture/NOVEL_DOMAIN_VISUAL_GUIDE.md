# Novel Domain Visual Architecture

## 🎯 Core Principle

```
╔═══════════════════════════════════════════════════════════╗
║  GOLDEN RULE: Domain Models Store IDs, Not Full Objects  ║
╚═══════════════════════════════════════════════════════════╝

Domain Layer      →  Store & return IDs (uint)
Application Layer →  Fetch & return full objects (struct)
```

## 🏗️ Layer Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      CLIENT REQUEST                         │
│                  GET /novels/1/details                      │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ↓
┌─────────────────────────────────────────────────────────────┐
│                   APPLICATION LAYER                         │
│                 (Cross-Domain Coordinator)                  │
│                                                             │
│  NovelManagementService                                     │
│  ┌─────────────────────────────────────────────────┐       │
│  │ GetNovelWithDetails(id)                         │       │
│  │                                                  │       │
│  │  1. novelService.GetNovelByID(1)    ──────┐    │       │
│  │     Returns: {id:1, created_by:5}          │    │       │
│  │                                             │    │       │
│  │  2. userRepo.GetUserByID(5)         ──────┼────┼───┐   │
│  │     Returns: {id:5, username:"john"}       │    │   │   │
│  │                                             │    │   │   │
│  │  3. mediaRepo.GetByID(10)           ──────┼────┼───┼──┐│
│  │     Returns: {id:10, url:"..."}            │    │   │  ││
│  │                                             │    │   │  ││
│  │  4. Combine all data                       │    │   │  ││
│  │     Returns: NovelWithDetailsDTO           │    │   │  ││
│  └─────────────────────────────────────────────────┘   │  ││
│                                                 │    │   │  ││
└─────────────────────────────────────────────────┼────┼───┼──┼┘
                                                  │    │   │  │
                         ┌────────────────────────┘    │   │  │
                         │                             │   │  │
                         ↓                             │   │  │
┌─────────────────────────────────────────────────────────┐│  │
│                     NOVEL DOMAIN                        ││  │
│                    (Pure & Independent)                 ││  │
│                                                         ││  │
│  NovelService                                           ││  │
│  ┌───────────────────────────────────────────┐         ││  │
│  │ GetNovelByID(1)                           │         ││  │
│  │   ↓                                       │         ││  │
│  │ NovelRepository                           │         ││  │
│  │   ↓                                       │         ││  │
│  │ SELECT * FROM novels WHERE id=1           │         ││  │
│  │   ↓                                       │         ││  │
│  │ Returns: Novel{                           │         ││  │
│  │   ID: 1,                                  │         ││  │
│  │   CreatedBy: 5,      ← Just ID            │         ││  │
│  │   CoverMediaId: 10   ← Just ID            │         ││  │
│  │ }                                         │         ││  │
│  └───────────────────────────────────────────┘         ││  │
│                                                         ││  │
│  ✅ No User import                                      ││  │
│  ✅ No Media import                                     ││  │
│  ✅ Only NovelRepository dependency                     ││  │
└─────────────────────────────────────────────────────────┘│  │
                                                           │  │
            ┌──────────────────────────────────────────────┘  │
            │                                                  │
            ↓                                                  │
┌─────────────────────────────────────────────────────────────┐│
│                      USER DOMAIN                            ││
│                                                             ││
│  UserRepository                                             ││
│  ┌───────────────────────────────────────────┐             ││
│  │ GetUserByID(5)                            │             ││
│  │   ↓                                       │             ││
│  │ SELECT * FROM users WHERE id=5            │             ││
│  │   ↓                                       │             ││
│  │ Returns: User{                            │             ││
│  │   ID: 5,                                  │             ││
│  │   Username: "john",                       │             ││
│  │   Email: "john@example.com"               │             ││
│  │ }                                         │             ││
│  └───────────────────────────────────────────┘             ││
└─────────────────────────────────────────────────────────────┘│
                                                               │
                  ┌────────────────────────────────────────────┘
                  │
                  ↓
┌─────────────────────────────────────────────────────────────┐
│                     MEDIA DOMAIN                            │
│                                                             │
│  MediaRepository                                            │
│  ┌───────────────────────────────────────────┐             │
│  │ GetByID(10)                               │             │
│  │   ↓                                       │             │
│  │ SELECT * FROM media WHERE id=10           │             │
│  │   ↓                                       │             │
│  │ Returns: Media{                           │             │
│  │   ID: 10,                                 │             │
│  │   URL: "/uploads/cover.jpg",              │             │
│  │   Type: "image/jpeg"                      │             │
│  │ }                                         │             │
│  └───────────────────────────────────────────┘             │
└─────────────────────────────────────────────────────────────┘
                         │
                         ↓
┌─────────────────────────────────────────────────────────────┐
│                    FINAL RESPONSE                           │
│                                                             │
│  NovelWithDetailsDTO {                                      │
│    id: 1,                                                   │
│    original_language: "en",                                 │
│    created_by: 5,                                           │
│    creator: {                  ← From User domain           │
│      id: 5,                                                 │
│      username: "john",                                      │
│      email: "john@example.com"                              │
│    },                                                       │
│    cover_media_id: 10,                                      │
│    cover_media: {              ← From Media domain          │
│      id: 10,                                                │
│      url: "/uploads/cover.jpg",                             │
│      type: "image/jpeg"                                     │
│    }                                                        │
│  }                                                          │
└─────────────────────────────────────────────────────────────┘
```

## 📦 Domain Structure Comparison

### ❌ WRONG: Domain with Cross-Domain Dependencies

```
┌─────────────────────────────────┐
│      Novel Domain (BAD)         │
│                                 │
│  import userModel  ← ❌         │
│  import mediaModel ← ❌         │
│                                 │
│  type Novel struct {            │
│    Creator   *User  ← ❌        │
│    CoverMedia *Media ← ❌       │
│  }                              │
│                                 │
│  NovelService {                 │
│    userRepo ← ❌                │
│    mediaRepo ← ❌               │
│  }                              │
│                                 │
│  ⚠️ Tight coupling              │
│  ⚠️ Hard to test                │
│  ⚠️ Can't split to microservice │
└─────────────────────────────────┘
```

### ✅ CORRECT: Pure Domain + Application Coordinator

```
┌─────────────────────────────────┐
│   Novel Domain (GOOD)           │
│                                 │
│  No cross-domain imports ✅     │
│                                 │
│  type Novel struct {            │
│    CreatedBy *uint ✅           │
│    CoverMediaId *uint ✅        │
│  }                              │
│                                 │
│  NovelService {                 │
│    novelRepo ✅                 │
│  }                              │
│                                 │
│  ✅ Pure & independent          │
│  ✅ Easy to test                │
│  ✅ Can be microservice         │
└─────────────────────────────────┘
         ↑
         │ Uses
         │
┌─────────────────────────────────┐
│  Application Layer (GOOD)       │
│                                 │
│  NovelManagementService {       │
│    novelService ✅              │
│    userRepo ✅                  │
│    mediaRepo ✅                 │
│  }                              │
│                                 │
│  Coordinates all domains ✅     │
└─────────────────────────────────┘
```

## 🔄 Create Novel Flow

```
┌─────────────────────────────────────────────────────────────┐
│  CLIENT: POST /novels/create                                │
│  {                                                          │
│    "original_language": "en",                               │
│    "created_by": 5,                                         │
│    "cover_media_id": 10                                     │
│  }                                                          │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ↓
┌────────────────────────────────────────────────────────────┐
│  STEP 1: Application Layer - Validate Creator             │
│  ┌──────────────────────────────────────────┐             │
│  │ userRepo.GetUserByID(5)                  │             │
│  │ ✅ User exists                            │             │
│  └──────────────────────────────────────────┘             │
└────────────────────┬───────────────────────────────────────┘
                     │
                     ↓
┌────────────────────────────────────────────────────────────┐
│  STEP 2: Application Layer - Validate Media               │
│  ┌──────────────────────────────────────────┐             │
│  │ mediaRepo.GetByID(10)                    │             │
│  │ ✅ Media exists                           │             │
│  └──────────────────────────────────────────┘             │
└────────────────────┬───────────────────────────────────────┘
                     │
                     ↓
┌────────────────────────────────────────────────────────────┐
│  STEP 3: Novel Domain - Create Novel                      │
│  ┌──────────────────────────────────────────┐             │
│  │ novelService.CreateNovel({               │             │
│  │   original_language: "en",               │             │
│  │   created_by: 5,      ← Just stores ID   │             │
│  │   cover_media_id: 10  ← Just stores ID   │             │
│  │ })                                       │             │
│  │                                          │             │
│  │ INSERT INTO novels (...)                 │             │
│  │ VALUES ('en', 5, 10, ...)                │             │
│  │                                          │             │
│  │ Returns: NovelDTO {id: 1, ...}           │             │
│  └──────────────────────────────────────────┘             │
└────────────────────┬───────────────────────────────────────┘
                     │
                     ↓
┌────────────────────────────────────────────────────────────┐
│  STEP 4: Application Layer - Fetch Full Details           │
│  ┌──────────────────────────────────────────┐             │
│  │ GetNovelWithDetails(1)                   │             │
│  │   → Fetches Novel + Creator + Media      │             │
│  │                                          │             │
│  │ Returns: NovelWithDetailsDTO {           │             │
│  │   id: 1,                                 │             │
│  │   creator: {id:5, username:"john"},      │             │
│  │   cover_media: {id:10, url:"..."}        │             │
│  │ }                                        │             │
│  └──────────────────────────────────────────┘             │
└────────────────────┬───────────────────────────────────────┘
                     │
                     ↓
┌────────────────────────────────────────────────────────────┐
│  RESPONSE TO CLIENT                                        │
│  {                                                         │
│    "success": true,                                        │
│    "data": {                                               │
│      "id": 1,                                              │
│      "original_language": "en",                            │
│      "creator": { "id": 5, "username": "john", ... },      │
│      "cover_media": { "id": 10, "url": "/uploads/..." }    │
│    }                                                       │
│  }                                                         │
└────────────────────────────────────────────────────────────┘
```

## 📂 File Organization

```
internal/
│
├── novel/                          ← DOMAIN LAYER
│   ├── model/
│   │   └── novel.go
│   │       ├── Novel struct        ✅ IDs: CreatedBy, CoverMediaId
│   │       └── NovelTranslation    ✅ ID: TranslatorId
│   │
│   ├── dto/
│   │   └── novel_dto.go
│   │       ├── NovelDTO            ✅ Has: created_by uint
│   │       └── TranslationDTO      ✅ Has: translator_id uint
│   │
│   ├── repository/
│   │   └── novel_repository.go
│   │       └── Methods             ✅ No joins with user/media tables
│   │
│   └── service/
│       └── novel_service.go
│           └── Dependencies        ✅ Only NovelRepository
│
└── application/                    ← APPLICATION LAYER
    ├── dto/
    │   └── novel_dto.go
    │       ├── NovelWithDetailsDTO ✅ Has: creator *UserBasicDTO
    │       └── TranslationDTO      ✅ Has: translator *UserBasicDTO
    │
    └── service/
        └── novel_management_service.go
            └── Dependencies        ✅ NovelService + UserRepo + MediaRepo
```

## 🎓 Mental Model

```
Think of it like a library system:

Novel Domain = Book Catalog
├── Stores: Book ID, Author Name, ISBN
├── Stores: BorrowerID (just a number)
└── Doesn't know: Who the borrower is, their address, etc.

Application Layer = Library Management System
├── Knows: Book Catalog (Novel Domain)
├── Knows: Member Database (User Domain)
├── Knows: Media Archive (Media Domain)
└── Can: Combine all information for complete view

When patron asks "Show me this book with borrower details":
1. Book Catalog gives: Book info + BorrowerID=5
2. Member Database gives: Who borrower 5 is
3. Library System combines: Complete information
```

## ✅ Benefits Visualization

```
┌─────────────────────────────────────────────────────────────┐
│                    BEFORE (Coupled)                         │
│                                                             │
│  ┌─────────┐                                                │
│  │  Novel  │──┐                                             │
│  │ Domain  │  │                                             │
│  │         │  │                                             │
│  │ imports │←─┼─── User Domain                             │
│  │ imports │←─┼─── Media Domain                            │
│  │ imports │←─┼─── ... more domains                        │
│  └─────────┘  │                                             │
│               └─── Can't split, can't test independently    │
│                                                             │
│  ⚠️  One big ball of mud                                    │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                  AFTER (Decoupled)                          │
│                                                             │
│  ┌─────────┐    ┌─────────┐    ┌─────────┐                │
│  │  Novel  │    │  User   │    │  Media  │                │
│  │ Domain  │    │ Domain  │    │ Domain  │                │
│  └────┬────┘    └────┬────┘    └────┬────┘                │
│       │              │              │                       │
│       └──────────────┼──────────────┘                       │
│                      ↓                                      │
│            ┌─────────────────┐                              │
│            │  Application    │                              │
│            │     Layer       │                              │
│            │  (Coordinator)  │                              │
│            └─────────────────┘                              │
│                                                             │
│  ✅ Clean separation                                        │
│  ✅ Independent testing                                     │
│  ✅ Can become microservices                                │
│  ✅ Easy to maintain                                        │
└─────────────────────────────────────────────────────────────┘
```

---

**Remember:** If it crosses domains, put it in Application layer!
