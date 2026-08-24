package service

import (
	"github.com/miniqxt/backend/internal/model"
	"github.com/miniqxt/backend/internal/pkg/apperr"
	"github.com/miniqxt/backend/internal/pkg/shuffle"
	"github.com/miniqxt/backend/internal/repo"
	"github.com/miniqxt/backend/internal/timeutil"
)

type PaperItemIn struct {
	QuestionID uint64  `json:"question_id"`
	Score      float64 `json:"score"`
	GroupName  string  `json:"group_name"`
}

type RenderedItem struct {
	QuestionID uint64                 `json:"question_id"`
	Type       string                 `json:"type"`
	Stem       string                 `json:"stem"`
	Score      float64                `json:"score"`
	GroupName  string                 `json:"group_name"`
	Difficulty int                    `json:"difficulty"`
	DisplayNo  int                    `json:"display_no"`
	Options    []model.QuestionOption `json:"options"`
}

func (a *App) ListPapers(tid uint64) ([]model.Paper, error) {
	var list []model.Paper
	err := repo.Tenant(a.DB, tid).Preload("Items").Order("id desc").Find(&list).Error
	return list, err
}

func (a *App) GetPaper(tid, id uint64) (model.Paper, error) {
	var p model.Paper
	err := repo.Tenant(a.DB, tid).Preload("Items.Question.Options").First(&p, id).Error
	if err != nil {
		return p, apperr.NotFound
	}
	return p, nil
}

func (a *App) CreatePaper(tid uint64, title string, sq, so bool) (model.Paper, error) {
	if title == "" {
		return model.Paper{}, apperr.Validation.With("试卷标题必填")
	}
	p := model.Paper{TenantID: tid, Title: title, ShuffleQuestions: sq, ShuffleOptions: so, CreatedAt: model.Time(timeutil.Now())}
	a.DB.Create(&p)
	return p, nil
}

func (a *App) SavePaperItems(tid, id uint64, title string, sq, so bool, items []PaperItemIn) (model.Paper, error) {
	p, err := a.GetPaper(tid, id)
	if err != nil {
		return p, err
	}
	if title != "" {
		p.Title = title
	}
	p.ShuffleQuestions = sq
	p.ShuffleOptions = so
	var total float64
	a.DB.Where("paper_id = ?", id).Delete(&model.PaperItem{})
	p.Items = nil
	for i, it := range items {
		var q model.Question
		if err := repo.Tenant(a.DB, tid).First(&q, it.QuestionID).Error; err != nil {
			return p, apperr.Validation.With("题目不属于本租户")
		}
		sc := it.Score
		if sc <= 0 {
			sc = q.Score
		}
		row := model.PaperItem{PaperID: id, QuestionID: it.QuestionID, Score: sc, GroupName: it.GroupName, SortNo: i}
		a.DB.Create(&row)
		total += sc
		p.Items = append(p.Items, row)
	}
	p.TotalScore = total
	p.UpdatedAt = model.Time(timeutil.Now())
	a.DB.Save(&p)
	return a.GetPaper(tid, id)
}

func (a *App) DeletePaper(tid, id uint64) error {
	if err := repo.Tenant(a.DB, tid).First(&model.Paper{}, id).Error; err != nil {
		return apperr.NotFound
	}
	a.DB.Where("paper_id = ?", id).Delete(&model.PaperItem{})
	repo.Tenant(a.DB, tid).Delete(&model.Paper{}, id)
	return nil
}

func (a *App) RenderPaper(p model.Paper, seed int64, hideCorrect bool) []RenderedItem {
	items := append([]model.PaperItem(nil), p.Items...)
	if p.ShuffleQuestions {
		ord := shuffle.Order(seed, len(items))
		items = shuffle.Apply(items, ord)
	}
	out := make([]RenderedItem, 0, len(items))
	for i, it := range items {
		opts := append([]model.QuestionOption(nil), it.Question.Options...)
		if p.ShuffleOptions && len(opts) > 1 {
			ord := shuffle.Order(shuffle.OptionSeed(seed, it.QuestionID), len(opts))
			opts = shuffle.Apply(opts, ord)
		}
		if hideCorrect {
			opts = repo.HideCorrect(opts)
		}
		out = append(out, RenderedItem{
			QuestionID: it.QuestionID, Type: it.Question.Type, Stem: it.Question.Stem,
			Score: it.Score, GroupName: it.GroupName, Difficulty: it.Question.Difficulty,
			DisplayNo: i + 1, Options: opts,
		})
	}
	return out
}
