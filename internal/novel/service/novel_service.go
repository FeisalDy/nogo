package service

import (
	commonDto "github.com/FeisalDy/nogo/internal/common/dto"
	"github.com/FeisalDy/nogo/internal/novel/dto"
	"github.com/FeisalDy/nogo/internal/novel/model"
	"github.com/FeisalDy/nogo/internal/novel/repository"
)

type NovelService struct {
	novelRepo *repository.NovelRepository
}

func NewNovelService(novelRepo *repository.NovelRepository) *NovelService {
	return &NovelService{
		novelRepo: novelRepo,
	}
}

func (s *NovelService) CreateNovel(createDTO *dto.CreateNovelDTO) (*dto.NovelDTO, error) {
	novel := &model.Novel{
		OriginalLanguage: createDTO.OriginalLanguage,
		OriginalAuthor:   createDTO.OriginalAuthor,
		Status:           createDTO.Status,
		Source:           createDTO.Source,
		WordCount:        createDTO.WordCount,
		CoverMediaId:     createDTO.CoverMediaId,
		CreatedBy:        createDTO.CreatedBy,
	}

	if err := s.novelRepo.Create(novel); err != nil {
		return nil, err
	}

	return s.toNovelDTO(novel), nil
}

func (s *NovelService) GetNovelByID(id uint) (*dto.NovelDTO, error) {
	novel, err := s.novelRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return s.toNovelDTO(novel), nil
}

func (s *NovelService) GetAllNovelsWithCursor(req *commonDto.CursorPaginationRequest) ([]dto.NovelDTO, commonDto.CursorPageInfo, error) {
	novels, pageInfo, err := s.novelRepo.GetAllWithCursor(req)
	if err != nil {
		return nil, commonDto.CursorPageInfo{}, err
	}

	novelDTOs := make([]dto.NovelDTO, len(novels))
	for i, novel := range novels {
		novelDTOs[i] = *s.toNovelDTO(&novel)
	}

	return novelDTOs, pageInfo, nil
}

func (s *NovelService) UpdateNovel(id uint, updateDTO *dto.UpdateNovelDTO) (*dto.NovelDTO, error) {
	novel, err := s.novelRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if updateDTO.OriginalLanguage != nil {
		novel.OriginalLanguage = *updateDTO.OriginalLanguage
	}
	if updateDTO.OriginalAuthor != nil {
		novel.OriginalAuthor = updateDTO.OriginalAuthor
	}
	if updateDTO.Status != nil {
		novel.Status = updateDTO.Status
	}
	if updateDTO.Source != nil {
		novel.Source = updateDTO.Source
	}
	if updateDTO.WordCount != nil {
		novel.WordCount = updateDTO.WordCount
	}
	if updateDTO.CoverMediaId != nil {
		novel.CoverMediaId = updateDTO.CoverMediaId
	}

	if err := s.novelRepo.Update(novel); err != nil {
		return nil, err
	}

	return s.toNovelDTO(novel), nil
}

func (s *NovelService) DeleteNovel(id uint) error {
	return s.novelRepo.Delete(id)
}

// ==================== Translation Methods ====================

func (s *NovelService) CreateTranslation(createDTO *dto.CreateNovelTranslationDTO) (*dto.NovelTranslationDTO, error) {
	translation := &model.NovelTranslation{
		NovelId:      createDTO.NovelId,
		Language:     createDTO.Language,
		Title:        createDTO.Title,
		Synopsis:     createDTO.Synopsis,
		TranslatorId: createDTO.TranslatorId,
	}

	if err := s.novelRepo.CreateTranslation(translation); err != nil {
		return nil, err
	}

	return s.toTranslationDTO(translation), nil
}

// GetTranslationByID retrieves a translation by ID
func (s *NovelService) GetTranslationByID(id uint) (*dto.NovelTranslationDTO, error) {
	translation, err := s.novelRepo.GetTranslationByID(id)
	if err != nil {
		return nil, err
	}

	return s.toTranslationDTO(translation), nil
}

// GetTranslationsByNovelID retrieves all translations for a novel
func (s *NovelService) GetTranslationsByNovelID(novelID uint) ([]dto.NovelTranslationDTO, error) {
	translations, err := s.novelRepo.GetTranslationsByNovelID(novelID)
	if err != nil {
		return nil, err
	}

	translationDTOs := make([]dto.NovelTranslationDTO, len(translations))
	for i, translation := range translations {
		translationDTOs[i] = *s.toTranslationDTO(&translation)
	}

	return translationDTOs, nil
}

// UpdateTranslation updates a translation
func (s *NovelService) UpdateTranslation(id uint, updateDTO *dto.UpdateNovelTranslationDTO) (*dto.NovelTranslationDTO, error) {
	translation, err := s.novelRepo.GetTranslationByID(id)
	if err != nil {
		return nil, err
	}

	if updateDTO.Title != nil {
		translation.Title = *updateDTO.Title
	}
	if updateDTO.Synopsis != nil {
		translation.Synopsis = updateDTO.Synopsis
	}
	if updateDTO.TranslatorId != nil {
		translation.TranslatorId = updateDTO.TranslatorId
	}

	if err := s.novelRepo.UpdateTranslation(translation); err != nil {
		return nil, err
	}

	return s.toTranslationDTO(translation), nil
}

// DeleteTranslation deletes a translation
func (s *NovelService) DeleteTranslation(id uint) error {
	return s.novelRepo.DeleteTranslation(id)
}

// ==================== Helper Methods ====================

// toNovelDTO converts a Novel model to NovelDTO
// Note: Only includes IDs for foreign keys, not full objects
func (s *NovelService) toNovelDTO(novel *model.Novel) *dto.NovelDTO {
	return &dto.NovelDTO{
		ID:               novel.ID,
		OriginalLanguage: novel.OriginalLanguage,
		OriginalAuthor:   novel.OriginalAuthor,
		Status:           novel.Status,
		Source:           novel.Source,
		WordCount:        novel.WordCount,
		CoverMediaId:     novel.CoverMediaId, // Just the ID
		CreatedBy:        novel.CreatedBy,    // Just the ID
		CreatedAt:        novel.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:        novel.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// toTranslationDTO converts a NovelTranslation model to NovelTranslationDTO
func (s *NovelService) toTranslationDTO(translation *model.NovelTranslation) *dto.NovelTranslationDTO {
	return &dto.NovelTranslationDTO{
		ID:           translation.ID,
		NovelId:      translation.NovelId,
		Language:     translation.Language,
		Title:        translation.Title,
		Synopsis:     translation.Synopsis,
		TranslatorId: translation.TranslatorId, // Just the ID
		CreatedAt:    translation.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    translation.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
