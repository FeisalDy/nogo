package repository

import (
	"github.com/FeisalDy/nogo/internal/media/model"
	"gorm.io/gorm"
)

type MediaRepository struct {
	db *gorm.DB
}

func NewMediaRepository(db *gorm.DB) *MediaRepository {
	return &MediaRepository{db: db}
}

func (r *MediaRepository) WithTx(tx *gorm.DB) *MediaRepository {
	return &MediaRepository{db: tx}
}

func (r *MediaRepository) Create(media *model.Media) error {
	return r.db.Create(media).Error
}

func (r *MediaRepository) GetByID(id uint) (*model.Media, error) {
	var media model.Media
	err := r.db.First(&media, id).Error
	return &media, err
}

func (r *MediaRepository) Delete(id uint) error {
	return r.db.Delete(&model.Media{}, id).Error
}
