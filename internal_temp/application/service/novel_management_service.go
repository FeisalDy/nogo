package service

import (
	"errors"

	"gorm.io/gorm"

	appDto "github.com/FeisalDy/nogo/internal/application/dto"
	novelDto "github.com/FeisalDy/nogo/internal/novel/dto"
	novelRepo "github.com/FeisalDy/nogo/internal/novel/repository"
	novelService "github.com/FeisalDy/nogo/internal/novel/service"
	userRepo "github.com/FeisalDy/nogo/internal/user/repository"
	// When Media domain is created:
	// mediaRepo "github.com/FeisalDy/nogo/internal/media/repository"
)

// NovelManagementService handles cross-domain operations for novels
// This service coordinates between Novel, User, and Media domains
// Following DDD principles:
// - Application layer coordinates multiple domains
// - Domain services remain pure and independent
type NovelManagementService struct {
	novelService *novelService.NovelService
	novelRepo    *novelRepo.NovelRepository
	userRepo     *userRepo.UserRepository
	// mediaRepo    *mediaRepo.MediaRepository  // Add when Media domain is created
	db *gorm.DB
}

func NewNovelManagementService(
	novelService *novelService.NovelService,
	novelRepo *novelRepo.NovelRepository,
	userRepo *userRepo.UserRepository,
	db *gorm.DB,
) *NovelManagementService {
	return &NovelManagementService{
		novelService: novelService,
		novelRepo:    novelRepo,
		userRepo:     userRepo,
		db:           db,
	}
}

// GetNovelWithDetails retrieves a novel with creator and cover media details
// This is a cross-domain operation that:
// 1. Gets novel from Novel domain
// 2. Gets creator from User domain (if CreatedBy is set)
// 3. Gets cover media from Media domain (if CoverMediaId is set - when implemented)
func (s *NovelManagementService) GetNovelWithDetails(novelID uint) (*appDto.NovelWithDetailsDTO, error) {
	// 1. Get novel from Novel domain
	novelDTO, err := s.novelService.GetNovelByID(novelID)
	if err != nil {
		return nil, err
	}

	// 2. Build response with novel data
	response := &appDto.NovelWithDetailsDTO{
		ID:               novelDTO.ID,
		OriginalLanguage: novelDTO.OriginalLanguage,
		OriginalAuthor:   novelDTO.OriginalAuthor,
		Status:           novelDTO.Status,
		Source:           novelDTO.Source,
		WordCount:        novelDTO.WordCount,
		CreatedAt:        novelDTO.CreatedAt,
		UpdatedAt:        novelDTO.UpdatedAt,
	}

	// 3. Get creator from User domain (if exists)
	if novelDTO.CreatedBy != nil {
		creator, err := s.userRepo.GetUserByID(*novelDTO.CreatedBy)
		if err == nil {
			response.Creator = &appDto.UserBasicDTO{
				ID:       creator.ID,
				Username: creator.Username,
				Email:    creator.Email,
			}
		}
		// If user not found, continue without creator info
	}

	// 4. Get cover media from Media domain (when implemented)
	// if novelDTO.CoverMediaId != nil {
	// 	media, err := s.mediaRepo.GetByID(*novelDTO.CoverMediaId)
	// 	if err == nil {
	// 		response.CoverMedia = &appDto.MediaBasicDTO{
	// 			ID:       media.ID,
	// 			URL:      media.URL,
	// 			Type:     media.Type,
	// 			FileName: media.FileName,
	// 		}
	// 	}
	// }

	return response, nil
}

// GetNovelComplete retrieves complete novel data with all translations and related entities
func (s *NovelManagementService) GetNovelComplete(novelID uint) (*appDto.NovelCompleteDTO, error) {
	// 1. Get novel with details
	novelWithDetails, err := s.GetNovelWithDetails(novelID)
	if err != nil {
		return nil, err
	}

	// 2. Get all translations
	translations, err := s.novelService.GetTranslationsByNovelID(novelID)
	if err != nil {
		return nil, err
	}

	// 3. Build translation DTOs with translator info
	translationDetails := make([]appDto.NovelTranslationWithDetailsDTO, len(translations))
	for i, trans := range translations {
		translationDetails[i] = appDto.NovelTranslationWithDetailsDTO{
			ID:        trans.ID,
			NovelId:   trans.NovelId,
			Language:  trans.Language,
			Title:     trans.Title,
			Synopsis:  trans.Synopsis,
			CreatedAt: trans.CreatedAt,
			UpdatedAt: trans.UpdatedAt,
		}

		// Get translator info from User domain (if exists)
		if trans.TranslatorId != nil {
			translator, err := s.userRepo.GetUserByID(*trans.TranslatorId)
			if err == nil {
				translationDetails[i].Translator = &appDto.UserBasicDTO{
					ID:       translator.ID,
					Username: translator.Username,
					Email:    translator.Email,
				}
			}
		}
	}

	return &appDto.NovelCompleteDTO{
		Novel:        *novelWithDetails,
		Translations: translationDetails,
	}, nil
}

// CreateNovelWithCreator creates a novel and validates creator exists
func (s *NovelManagementService) CreateNovelWithCreator(createDTO *novelDto.CreateNovelDTO, creatorID uint) (*appDto.NovelWithDetailsDTO, error) {
	// 1. Validate creator exists in User domain
	_, err := s.userRepo.GetUserByID(creatorID)
	if err != nil {
		return nil, errors.New("creator user not found")
	}

	// 2. Validate cover media exists (when Media domain is implemented)
	// if createDTO.CoverMediaId != nil {
	// 	_, err := s.mediaRepo.GetByID(*createDTO.CoverMediaId)
	// 	if err != nil {
	// 		return nil, errors.New("cover media not found")
	// 	}
	// }

	// 3. Set creator ID
	createDTO.CreatedBy = &creatorID

	// 4. Create novel in Novel domain
	novelDTO, err := s.novelService.CreateNovel(createDTO)
	if err != nil {
		return nil, err
	}

	// 5. Get novel with full details
	return s.GetNovelWithDetails(novelDTO.ID)
}

// CreateTranslationWithTranslator creates a translation and validates translator exists
func (s *NovelManagementService) CreateTranslationWithTranslator(
	createDTO *novelDto.CreateNovelTranslationDTO,
	translatorID *uint,
) (*appDto.NovelTranslationWithDetailsDTO, error) {
	// 1. Validate novel exists
	_, err := s.novelService.GetNovelByID(createDTO.NovelId)
	if err != nil {
		return nil, errors.New("novel not found")
	}

	// 2. Validate translator exists (if provided)
	var translator *appDto.UserBasicDTO
	if translatorID != nil {
		user, err := s.userRepo.GetUserByID(*translatorID)
		if err != nil {
			return nil, errors.New("translator user not found")
		}
		translator = &appDto.UserBasicDTO{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		}
		createDTO.TranslatorId = translatorID
	}

	// 3. Create translation in Novel domain
	translationDTO, err := s.novelService.CreateTranslation(createDTO)
	if err != nil {
		return nil, err
	}

	// 4. Build response with translator info
	return &appDto.NovelTranslationWithDetailsDTO{
		ID:         translationDTO.ID,
		NovelId:    translationDTO.NovelId,
		Language:   translationDTO.Language,
		Title:      translationDTO.Title,
		Synopsis:   translationDTO.Synopsis,
		CreatedAt:  translationDTO.CreatedAt,
		UpdatedAt:  translationDTO.UpdatedAt,
		Translator: translator,
	}, nil
}
