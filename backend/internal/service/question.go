package service

import (
	"strings"

	"github.com/miniqxt/backend/internal/model"
	"github.com/miniqxt/backend/internal/pkg/apperr"
	"github.com/miniqxt/backend/internal/repo"
	"github.com/miniqxt/backend/internal/timeutil"
)

type OptionIn struct {
	ID        uint64 `json:"id"`
	Label     string `json:"label"`
	IsCorrect bool   `json:"is_correct"`
}

type QuestionIn struct {
	Type       string     `json:"type"`
	Stem       string     `json:"stem"`
	Difficulty int        `json:"difficulty"`
	Tags       string     `json:"tags"`
	Score      float64    `json:"score"`
	Options    []OptionIn `json:"options"`
}

func validType(t string) bool {
	switch t {
	case model.QSingle, model.QMulti, model.QJudge, model.QEssay:
		return true
	}
	return false
}

func (a *App) ListQuestions(tid uint64, typ, tag, q string, page, size int) ([]model.Question, int64, error) {
	tx := repo.Tenant(a.DB.Model(&model.Question{}), tid)
	if typ != "" {
		tx = tx.Where("type = ?", typ)
	}
	if tag != "" {
		tx = tx.Where("tags LIKE ?", "%"+tag+"%")
	}
	if q != "" {
		tx = tx.Where("stem LIKE ?", "%"+q+"%")
	}
	var total int64
	tx.Count(&total)
	var list []model.Question
	err := repo.Page(tx.Preload("Options"), page, size).Order("id desc").Find(&list).Error
	return list, total, err
}

func (a *App) GetQuestion(tid, id uint64) (model.Question, error) {
	var qq model.Question
	if err := repo.Tenant(a.DB, tid).Preload("Options").First(&qq, id).Error; err != nil {
		return qq, apperr.NotFound
	}
	return qq, nil
}

func (a *App) CreateQuestion(tid uint64, in QuestionIn) (model.Question, error) {
	if !validType(in.Type) || strings.TrimSpace(in.Stem) == "" {
		return model.Question{}, apperr.Validation.With("题型或题干不合法")
	}
	if in.Type != model.QEssay && len(in.Options) < 2 {
		return model.Question{}, apperr.Validation.With("客观题至少两个选项")
	}
	if in.Score <= 0 {
		in.Score = 5
	}
	if in.Difficulty < 1 {
		in.Difficulty = 1
	}
	qq := model.Question{
		TenantID: tid, Type: in.Type, Stem: in.Stem, Difficulty: in.Difficulty,
		Tags: in.Tags, Score: in.Score, CreatedAt: model.Time(timeutil.Now()),
	}
	if err := a.DB.Create(&qq).Error; err != nil {
		return qq, err
	}
	for i, o := range in.Options {
		op := model.QuestionOption{QuestionID: qq.ID, Label: o.Label, IsCorrect: o.IsCorrect, SortNo: i}
		a.DB.Create(&op)
		qq.Options = append(qq.Options, op)
	}
	return qq, nil
}

func (a *App) UpdateQuestion(tid, id uint64, in QuestionIn) (model.Question, error) {
	qq, err := a.GetQuestion(tid, id)
	if err != nil {
		return qq, err
	}
	if in.Stem != "" {
		qq.Stem = in.Stem
	}
	if validType(in.Type) {
		qq.Type = in.Type
	}
	if in.Difficulty > 0 {
		qq.Difficulty = in.Difficulty
	}
	if in.Tags != "" {
		qq.Tags = in.Tags
	}
	if in.Score > 0 {
		qq.Score = in.Score
	}
	qq.UpdatedAt = model.Time(timeutil.Now())
	a.DB.Save(&qq)
	if in.Options != nil {
		a.DB.Where("question_id = ?", qq.ID).Delete(&model.QuestionOption{})
		qq.Options = nil
		for i, o := range in.Options {
			op := model.QuestionOption{QuestionID: qq.ID, Label: o.Label, IsCorrect: o.IsCorrect, SortNo: i}
			a.DB.Create(&op)
			qq.Options = append(qq.Options, op)
		}
	}
	return qq, nil
}

func (a *App) DeleteQuestion(tid, id uint64) error {
	if err := repo.Tenant(a.DB, tid).First(&model.Question{}, id).Error; err != nil {
		return apperr.NotFound
	}
	a.DB.Where("question_id = ?", id).Delete(&model.QuestionOption{})
	repo.Tenant(a.DB, tid).Delete(&model.Question{}, id)
	return nil
}
