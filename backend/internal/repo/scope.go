package repo

import (
	"github.com/miniqxt/backend/internal/model"
	"gorm.io/gorm"
)

func Tenant(db *gorm.DB, tenantID uint64) *gorm.DB {
	return db.Where("tenant_id = ?", tenantID)
}

func Page(db *gorm.DB, page, size int) *gorm.DB {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 200 {
		size = 20
	}
	return db.Offset((page - 1) * size).Limit(size)
}

func HideCorrect(opts []model.QuestionOption) []model.QuestionOption {
	out := make([]model.QuestionOption, len(opts))
	for i, o := range opts {
		o.IsCorrect = false
		out[i] = o
	}
	return out
}
