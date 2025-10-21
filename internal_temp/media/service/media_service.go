package service

import (
	"github.com/FeisalDy/nogo/internal/media/dto"
	"github.com/FeisalDy/nogo/internal/media/model"
	"github.com/FeisalDy/nogo/internal/media/repository"
)

type MediaService struct {
	mediaRepo *repository.MediaRepository
}

func NewMediaService(mediaRepo *repository.MediaRepository) *MediaService {
	return &MediaService{
		mediaRepo: mediaRepo,
	}
}

func (s *MediaService) CreateMedia(createDTO *dto.CreateMediaDTO) (*dto.MediaDTO, error) {
	media := &model.Media{
		UploadBy:    createDTO.UploadBy,
		URL:         createDTO.URL,
		Type:        createDTO.Type,
		Description: createDTO.Description,
		FileSize:    createDTO.FileSize,
		MimeType:    createDTO.MimeType,
	}

	if err := s.mediaRepo.Create(media); err != nil {
		return nil, err
	}

	return s.toMediaDTO(media), nil
}

func (s *MediaService) GetMediaByID(id uint) (*dto.MediaDTO, error) {
	media, err := s.mediaRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return s.toMediaDTO(media), nil
}

func (s *MediaService) DeleteMedia(id uint) error {
	return s.mediaRepo.Delete(id)
}

func (s *MediaService) toMediaDTO(media *model.Media) *dto.MediaDTO {
	return &dto.MediaDTO{
		ID:          media.ID,
		CreatedAt:   media.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   media.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UploadBy:    media.UploadBy,
		URL:         media.URL,
		Type:        media.Type,
		Description: media.Description,
		FileSize:    media.FileSize,
		MimeType:    media.MimeType,
	}
}
