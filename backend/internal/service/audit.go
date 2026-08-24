package service

import (
	"github.com/miniqxt/backend/internal/engine"
	"github.com/miniqxt/backend/internal/model"
	"github.com/miniqxt/backend/internal/pkg/apperr"
	"github.com/miniqxt/backend/internal/pkg/shuffle"
	"github.com/miniqxt/backend/internal/repo"
	"github.com/miniqxt/backend/internal/timeutil"
	"gorm.io/gorm/clause"
)

func (a *App) Timeline(tid, sid uint64) (map[string]any, error) {
	var s model.ExamSession
	if err := repo.Tenant(a.DB, tid).First(&s, sid).Error; err != nil {
		return nil, apperr.NotFound
	}
	var traces []model.AnswerTrace
	var cheats []model.AntiCheatEvent
	var flags []model.AuditFlag
	a.DB.Where("session_id = ?", sid).Order("seq asc, occurred_at asc").Find(&traces)
	a.DB.Where("session_id = ?", sid).Order("occurred_at asc").Find(&cheats)
	a.DB.Where("session_id = ?", sid).Find(&flags)

	e, err := a.GetExam(tid, s.ExamID)
	if err != nil {
		return nil, err
	}
	p, err := a.GetPaper(tid, e.PaperID)
	if err != nil {
		return nil, err
	}
	pos := map[uint64]int{}
	items := p.Items
	if p.ShuffleQuestions {
		ord := shuffle.Order(s.ShuffleSeed, len(items))
		items = shuffle.Apply(items, ord)
	}
	for i, it := range items {
		pos[it.QuestionID] = i + 1
	}
	optLabel := map[uint64]string{}
	for _, it := range p.Items {
		for _, o := range it.Question.Options {
			optLabel[o.ID] = o.Label
		}
	}
	events := make([]map[string]any, 0, len(traces))
	for _, t := range traces {
		from := ""
		to := ""
		if t.FromOptionID != nil {
			from = optLabel[*t.FromOptionID]
		}
		if t.ToOptionID != nil {
			to = optLabel[*t.ToOptionID]
		}
		events = append(events, map[string]any{
			"seq": t.Seq, "question_id": t.QuestionID, "display_no": pos[t.QuestionID],
			"from_option_id": t.FromOptionID, "to_option_id": t.ToOptionID,
			"from_label": from, "to_label": to, "occurred_at": t.OccurredAt,
			"answer_text": t.AnswerText,
		})
	}
	return map[string]any{
		"session": s, "traces": events, "anti_cheat": cheats, "flags": flags, "paper": a.RenderPaper(p, s.ShuffleSeed, false),
	}, nil
}

func (a *App) RunAudit(tid, sid uint64) ([]model.AuditFlag, error) {
	var s model.ExamSession
	if err := repo.Tenant(a.DB, tid).First(&s, sid).Error; err != nil {
		return nil, apperr.NotFound
	}
	var traces []model.AnswerTrace
	var cheats []model.AntiCheatEvent
	var sub model.Submission
	a.DB.Where("session_id = ?", sid).Order("seq asc").Find(&traces)
	a.DB.Where("session_id = ?", sid).Find(&cheats)
	_ = a.DB.Where("session_id = ?", sid).Preload("Answers").First(&sub).Error
	qIDs := make([]uint64, 0)
	for _, x := range sub.Answers {
		qIDs = append(qIDs, x.QuestionID)
	}
	var qs []model.Question
	if len(qIDs) > 0 {
		a.DB.Preload("Options").Where("id IN ?", qIDs).Find(&qs)
	}
	qmap := map[uint64]model.Question{}
	for _, q := range qs {
		qmap[q.ID] = q
	}
	hits := engine.Analyze(traces, cheats, sub.Answers, qmap)
	now := model.Time(timeutil.Now())
	a.DB.Where("session_id = ?", sid).Delete(&model.AuditFlag{})
	out := make([]model.AuditFlag, 0, len(hits))
	for _, h := range hits {
		f := model.AuditFlag{
			TenantID: tid, SessionID: sid, UserID: s.UserID, ExamID: s.ExamID,
			RuleCode: h.Rule, Title: h.Title, Evidence: engine.EvidenceJSON(h.Evidence),
			Severity: h.Severity, CreatedAt: now,
		}
		a.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&f)
		out = append(out, f)
	}
	return out, nil
}

func (a *App) ListFlags(tid uint64) ([]model.AuditFlag, error) {
	var list []model.AuditFlag
	err := repo.Tenant(a.DB, tid).Order("id desc").Limit(200).Find(&list).Error
	return list, err
}
