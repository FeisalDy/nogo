package repository

import (
	"gorm.io/gorm"

	"github.com/FeisalDy/nogo/internal/common/dto"
	"github.com/FeisalDy/nogo/internal/common/utils"
	"github.com/FeisalDy/nogo/internal/novel/model"
)

type NovelRepository struct {
	db *gorm.DB
}

func NewNovelRepository(db *gorm.DB) *NovelRepository {
	return &NovelRepository{db: db}
}

func (r *NovelRepository) WithTx(tx *gorm.DB) *NovelRepository {
	return &NovelRepository{db: tx}
}

func (r *NovelRepository) Create(novel *model.Novel) error {
	return r.db.Create(novel).Error
}

func (r *NovelRepository) GetByID(id uint) (*model.Novel, error) {
	var novel model.Novel
	err := r.db.First(&novel, id).Error
	return &novel, err
}

func (r *NovelRepository) GetAllWithCursor(req *dto.CursorPaginationRequest) ([]model.Novel, dto.CursorPageInfo, error) {
	baseQuery := r.db.Model(&model.Novel{})
	return utils.PaginateWithIDGetter[model.Novel](baseQuery, req)
}

func (r *NovelRepository) Update(novel *model.Novel) error {
	return r.db.Save(novel).Error
}

func (r *NovelRepository) Delete(id uint) error {
	return r.db.Delete(&model.Novel{}, id).Error
}

// ==================== Translation Methods ====================

func (r *NovelRepository) CreateTranslation(translation *model.NovelTranslation) error {
	return r.db.Create(translation).Error
}

// GetTranslationByID retrieves a translation by ID
func (r *NovelRepository) GetTranslationByID(id uint) (*model.NovelTranslation, error) {
	var translation model.NovelTranslation
	err := r.db.First(&translation, id).Error
	return &translation, err
}

func (r *NovelRepository) GetTranslationsByNovelID(novelID uint) ([]model.NovelTranslation, error) {
	var translations []model.NovelTranslation
	err := r.db.Where("language = ?").Where("novel_id = ?", novelID).Find(&translations).Error
	return translations, err
}

func (r *NovelRepository) GetTranslationByNovelAndLanguage(novelID uint, language string) (*model.NovelTranslation, error) {
	var translation model.NovelTranslation
	err := r.db.Where("novel_id = ? AND language = ?", novelID, language).First(&translation).Error
	return &translation, err
}

func (r *NovelRepository) UpdateTranslation(translation *model.NovelTranslation) error {
	return r.db.Save(translation).Error
}

func (r *NovelRepository) DeleteTranslation(id uint) error {
	return r.db.Delete(&model.NovelTranslation{}, id).Error
}

func (r *NovelRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&model.Novel{}).Count(&count).Error
	return count, err
}
