package novel

import (
	"gorm.io/gorm"
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

func (r *NovelRepository) Create(novel *Novel) error {
	return r.db.Create(novel).Error
}

func (r *NovelRepository) GetByID(id uint) (*Novel, error) {
	var novel Novel
	err := r.db.First(&novel, id).Error
	return &novel, err
}

func (r *NovelRepository) Update(novel *Novel) error {
	return r.db.Save(novel).Error
}

func (r *NovelRepository) Delete(id uint) error {
	return r.db.Delete(&Novel{}, id).Error
}

// ==================== Translation Methods ====================

func (r *NovelRepository) CreateTranslation(translation *NovelTranslation) error {
	return r.db.Create(translation).Error
}

// GetTranslationByID retrieves a translation by ID
func (r *NovelRepository) GetTranslationByID(id uint) (*NovelTranslation, error) {
	var translation NovelTranslation
	err := r.db.First(&translation, id).Error
	return &translation, err
}

func (r *NovelRepository) GetTranslationsByNovelID(novelID uint) ([]NovelTranslation, error) {
	var translations []NovelTranslation
	err := r.db.Where("language = ?").Where("novel_id = ?", novelID).Find(&translations).Error
	return translations, err
}

func (r *NovelRepository) GetTranslationByNovelAndLanguage(novelID uint, language string) (*NovelTranslation, error) {
	var translation NovelTranslation
	err := r.db.Where("novel_id = ? AND language = ?", novelID, language).First(&translation).Error
	return &translation, err
}

func (r *NovelRepository) UpdateTranslation(translation *NovelTranslation) error {
	return r.db.Save(translation).Error
}

func (r *NovelRepository) DeleteTranslation(id uint) error {
	return r.db.Delete(&NovelTranslation{}, id).Error
}

func (r *NovelRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&Novel{}).Count(&count).Error
	return count, err
}

func (r *NovelRepository) GetAllWithTranslationCursor(
	req *GetAllNovelRequestDTO,
	language string,
) ([]Novel,  error) {

	baseQuery := r.db.Model(&Novel{}).
		Preload("Translations", "language = ?", language)

	return utils.PaginateWithIDGetter[Novel](baseQuery, req)
}
