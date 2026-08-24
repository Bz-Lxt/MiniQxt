package service

import (
	"github.com/miniqxt/backend/internal/engine"
	"github.com/miniqxt/backend/internal/model"
	"github.com/miniqxt/backend/internal/pkg/apperr"
	"github.com/miniqxt/backend/internal/repo"
	"github.com/miniqxt/backend/internal/timeutil"
)

type AnswerIn struct {
	QuestionID uint64   `json:"question_id"`
	OptionIDs  []uint64 `json:"option_ids"`
	AnswerText string   `json:"answer_text"`
}

type TraceIn struct {
	QuestionID     uint64  `json:"question_id"`
	FromOptionID   *uint64 `json:"from_option_id"`
	ToOptionID     *uint64 `json:"to_option_id"`
	AnswerText     string  `json:"answer_text"`
	OccurredAt     string  `json:"occurred_at"`
	Seq            int     `json:"seq"`
}

func (a *App) IngestTraces(tid, uid, sid uint64, evs []TraceIn) error {
	if len(evs) == 0 {
		return apperr.Validation.With("events 不能为空")
	}
	if len(evs) > 500 {
		return apperr.Validation.With("单次最多 500 条流水")
	}
	var s model.ExamSession
	if err := repo.Tenant(a.DB, tid).Where("id = ? AND user_id = ?", sid, uid).First(&s).Error; err != nil {
		return apperr.NotFound
	}
	now := timeutil.Now()
	rows := make([]model.AnswerTrace, 0, len(evs))
	for _, e := range evs {
		if e.QuestionID == 0 {
			return apperr.Validation.With("question_id 必填")
		}
		ot, err := timeutil.ParseRFC3339(e.OccurredAt)
		if err != nil {
			ot = now
		}
		rows = append(rows, model.AnswerTrace{
			TenantID: tid, SessionID: sid, QuestionID: e.QuestionID,
			FromOptionID: e.FromOptionID, ToOptionID: e.ToOptionID,
			AnswerText: e.AnswerText, OccurredAt: model.Time(ot), Seq: e.Seq,
			CreatedAt: model.Time(now),
		})
	}
	a.Traces.Push(rows)
	return nil
}

func (a *App) SubmitExam(tid, uid, sid uint64, answers []AnswerIn) (map[string]any, error) {
	var s model.ExamSession
	if err := repo.Tenant(a.DB, tid).Where("id = ? AND user_id = ?", sid, uid).First(&s).Error; err != nil {
		return nil, apperr.NotFound
	}
	if s.Status == model.SessSubmitted || s.Status == model.SessForced {
		var exist model.Submission
		if err := a.DB.Where("session_id = ?", sid).First(&exist).Error; err == nil {
			return map[string]any{
				"submission_id": exist.ID, "status": exist.Status,
				"queue_position": a.queuePos(exist.ID), "integrity": exist.Integrity,
			}, nil
		}
		return nil, apperr.Conflict.With("已交卷")
	}
	now := timeutil.Now()
	expired := now.After(s.EndsAt.Time())
	if s.Status != model.SessInProgress && s.Status != model.SessExpired && !expired {
		return nil, apperr.SessionGone
	}

	var traceN, cheatN int64
	a.DB.Model(&model.AnswerTrace{}).Where("session_id = ?", sid).Count(&traceN)
	a.DB.Model(&model.AntiCheatEvent{}).Where("session_id = ?", sid).Count(&cheatN)
	integrity := engine.IntegrityFromTelemetry(s, traceN, cheatN)

	sessStatus := model.SessSubmitted
	if s.Status == model.SessForced {
		sessStatus = model.SessForced
	} else if expired {
		sessStatus = model.SessExpired
	}

	sub := model.Submission{
		TenantID: tid, SessionID: sid, ExamID: s.ExamID, UserID: uid,
		Status: model.SubQueued, Integrity: integrity,
		QueuedAt: model.Time(now), CreatedAt: model.Time(now),
	}
	if err := a.DB.Create(&sub).Error; err != nil {
		var exist model.Submission
		if e2 := a.DB.Where("session_id = ?", sid).First(&exist).Error; e2 == nil {
			return map[string]any{
				"submission_id": exist.ID, "status": exist.Status,
				"queue_position": a.queuePos(exist.ID), "integrity": exist.Integrity,
			}, nil
		}
		return nil, apperr.Conflict.With("重复交卷")
	}
	for _, ans := range answers {
		if ans.QuestionID == 0 {
			continue
		}
		row := model.SubmissionAnswer{
			SubmissionID: sub.ID, QuestionID: ans.QuestionID,
			OptionIDs: engine.JoinOptionIDs(ans.OptionIDs), AnswerText: ans.AnswerText,
		}
		a.DB.Create(&row)
	}
	a.DB.Model(&s).Updates(map[string]any{"status": sessStatus, "integrity": integrity})
	a.Grader.Enqueue(sub.ID)
	return map[string]any{
		"submission_id":  sub.ID,
		"status":         sub.Status,
		"queue_position": a.queuePos(sub.ID),
		"integrity":      integrity,
	}, nil
}

func (a *App) queuePos(id uint64) int {
	var n int64
	a.DB.Model(&model.Submission{}).
		Where("status = ? AND id <= ?", model.SubQueued, id).
		Count(&n)
	return int(n)
}

func (a *App) GetSubmission(tid, uid, id uint64, staff bool) (map[string]any, error) {
	var sub model.Submission
	tx := repo.Tenant(a.DB, tid).Preload("Answers").Where("id = ?", id)
	if !staff {
		tx = tx.Where("user_id = ?", uid)
	}
	if err := tx.First(&sub).Error; err != nil {
		return nil, apperr.NotFound
	}
	pos := 0
	if sub.Status == model.SubQueued || sub.Status == model.SubGrading {
		pos = a.queuePos(sub.ID)
	}
	return map[string]any{
		"submission":     sub,
		"queue_position": pos,
		"queue_depth":    a.Grader.Depth(),
	}, nil
}

type ManualItem struct {
	QuestionID uint64  `json:"question_id"`
	Score      float64 `json:"score"`
	Comment    string  `json:"comment"`
}

func (a *App) ManualScore(tid, sid uint64, items []ManualItem) (model.Submission, error) {
	var sub model.Submission
	if err := repo.Tenant(a.DB, tid).Preload("Answers").First(&sub, sid).Error; err != nil {
		return sub, apperr.NotFound
	}
	var exam model.Exam
	a.DB.First(&exam, sub.ExamID)
	var subj float64
	for i := range sub.Answers {
		ans := &sub.Answers[i]
		for _, it := range items {
			if it.QuestionID == ans.QuestionID {
				sc := it.Score
				ans.ManualScore = &sc
				ans.Comment = it.Comment
				a.DB.Model(ans).Updates(map[string]any{"manual_score": sc, "comment": it.Comment})
			}
		}
		if ans.IsEssay && ans.ManualScore != nil {
			subj += *ans.ManualScore
		}
	}
	total := sub.ObjectiveScore + subj
	status := model.SubPendingManual
	if allManualDone(sub.Answers) {
		status = model.SubGraded
	}
	now := model.Time(timeutil.Now())
	a.DB.Model(&sub).Updates(map[string]any{
		"subjective_score": subj, "total_score": total, "status": status,
		"pass": total+0.0001 >= exam.PassScore, "graded_at": now,
	})
	sub.SubjectiveScore = subj
	sub.TotalScore = total
	sub.Status = status
	return sub, nil
}

func allManualDone(ans []model.SubmissionAnswer) bool {
	for _, x := range ans {
		if x.IsEssay && x.ManualScore == nil {
			return false
		}
	}
	return true
}

func (a *App) ListSubmissions(tid uint64, status string, page, size int) ([]model.Submission, int64, error) {
	tx := repo.Tenant(a.DB.Model(&model.Submission{}), tid)
	if status != "" {
		tx = tx.Where("status = ?", status)
	}
	var total int64
	tx.Count(&total)
	var list []model.Submission
	err := repo.Page(tx, page, size).Order("id desc").Find(&list).Error
	return list, total, err
}
