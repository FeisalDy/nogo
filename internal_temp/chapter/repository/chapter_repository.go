package repository

import (
	"github.com/FeisalDy/nogo/internal/chapter/model"
	"github.com/FeisalDy/nogo/internal/common/dto"
	"github.com/FeisalDy/nogo/internal/common/utils"
	"gorm.io/gorm"
)

type ChapterRepository struct {
	db *gorm.DB
}

func NewChapterRepository(db *gorm.DB) *ChapterRepository {
	return &ChapterRepository{db: db}
}

func (r *ChapterRepository) WithTx(tx *gorm.DB) *ChapterRepository {
	return &ChapterRepository{db: tx}
}

func (r *ChapterRepository) Create(chapter *model.Chapter) error {
	return r.db.Create(chapter).Error
}

func (r *ChapterRepository) GetByID(id uint) (*model.Chapter, error) {
	var chapter model.Chapter
	err := r.db.First(&chapter, id).Error
	return &chapter, err
}

func (r *ChapterRepository) GetAll(req *dto.CursorPaginationRequest) ([]model.Chapter, dto.CursorPageInfo, error) {
	baseQuery := r.db.Model(&model.Chapter{})
	return utils.PaginateWithIDGetter[model.Chapter](baseQuery, req)
}
