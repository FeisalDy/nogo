package utils

import (
	"github.com/FeisalDy/nogo/internal/common/dto"
	"gorm.io/gorm"
)

func Paginate[T any](db *gorm.DB, p dto.PaginationRequestDTO, out *[]T) (dto.PaginationResultDTO[T], error) {
	var total int64

	if p.Page < 1 {
		p.Page = 1
	}
	if p.Limit <= 0 {
		p.Limit = 10
	}

	// Count total
	if err := db.Count(&total).Error; err != nil {
		return dto.PaginationResultDTO[T]{}, err
	}

	offset := (p.Page - 1) * p.Limit

	// Fetch page
	if err := db.Offset(offset).Limit(p.Limit).Find(out).Error; err != nil {
		return dto.PaginationResultDTO[T]{}, err
	}

	return dto.PaginationResultDTO[T]{
		Items:      *out,
		TotalCount: total,
		Page:       p.Page,
		Limit:      p.Limit,
	}, nil
}
