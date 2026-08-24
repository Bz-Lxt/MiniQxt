package queue

import (
	"context"
	"fmt"
	"sync"

	"github.com/miniqxt/backend/internal/engine"
	"github.com/miniqxt/backend/internal/logger"
	"github.com/miniqxt/backend/internal/model"
	"github.com/miniqxt/backend/internal/timeutil"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// SubmitQueue is an in-memory scheduler. Persistence lives in submission.status.
type SubmitQueue struct {
	db      *gorm.DB
	ch      chan uint64
	workers int
	Metrics *Metrics
	once    sync.Once
}

func NewSubmitQueue(db *gorm.DB, workers, size int, m *Metrics) *SubmitQueue {
	if workers < 1 {
		workers = 8
	}
	if size < 64 {
		size = 64
	}
	return &SubmitQueue{db: db, ch: make(chan uint64, size), workers: workers, Metrics: m}
}

func (q *SubmitQueue) Start(ctx context.Context) {
	for i := 0; i < q.workers; i++ {
		go q.loop(ctx, i)
	}
	logger.Info("grader workers started", "n", q.workers)
}

func (q *SubmitQueue) Depth() int { return len(q.ch) }

func (q *SubmitQueue) Workers() int { return q.workers }

// Enqueue is best-effort. Channel overflow is safe because status remains queued.
func (q *SubmitQueue) Enqueue(id uint64) {
	if id == 0 {
		return
	}
	q.Metrics.Enqueued.Add(1)
	select {
	case q.ch <- id:
	default:
		logger.Warn("grader channel full; recovery will pick up", "submission_id", id)
	}
}

func (q *SubmitQueue) Recover() int {
	var ids []uint64
	q.db.Model(&model.Submission{}).
		Where("status IN ?", []string{model.SubQueued, model.SubGrading}).
		Pluck("id", &ids)
	for _, id := range ids {
		q.Enqueue(id)
	}
	q.Metrics.Recovered.Add(int64(len(ids)))
	logger.Info("grader recovery scan", "n", len(ids))
	return len(ids)
}

func (q *SubmitQueue) loop(ctx context.Context, worker int) {
	for {
		select {
		case <-ctx.Done():
			return
		case id := <-q.ch:
			if err := q.grade(id); err != nil {
				q.Metrics.Failed.Add(1)
				logger.Error("grade failed", "worker", worker, "id", id, "err", err)
			}
		}
	}
}

func (q *SubmitQueue) grade(id uint64) error {
	var sub model.Submission
	if err := q.db.First(&sub, id).Error; err != nil {
		return err
	}
	if sub.Status == model.SubGraded || sub.Status == model.SubPendingManual {
		return nil
	}
	res := q.db.Model(&model.Submission{}).
		Where("id = ? AND status IN ?", id, []string{model.SubQueued, model.SubGrading}).
		Update("status", model.SubGrading)
	if res.RowsAffected == 0 {
		return nil
	}

	if err := q.db.Preload("Answers").First(&sub, id).Error; err != nil {
		return err
	}

	qIDs := make([]uint64, 0, len(sub.Answers))
	for _, a := range sub.Answers {
		qIDs = append(qIDs, a.QuestionID)
	}
	var questions []model.Question
	if err := q.db.Preload("Options").Where("id IN ?", qIDs).Find(&questions).Error; err != nil {
		return err
	}
	qmap := map[uint64]model.Question{}
	for _, qq := range questions {
		qmap[qq.ID] = qq
	}

	var paperItems []model.PaperItem
	var exam model.Exam
	if err := q.db.First(&exam, sub.ExamID).Error; err != nil {
		return err
	}
	q.db.Where("paper_id = ?", exam.PaperID).Find(&paperItems)
	scoreOf := map[uint64]float64{}
	for _, it := range paperItems {
		scoreOf[it.QuestionID] = it.Score
	}

	var obj float64
	hasEssay := false
	for i := range sub.Answers {
		a := &sub.Answers[i]
		qq := qmap[a.QuestionID]
		sc := scoreOf[a.QuestionID]
		if sc <= 0 {
			sc = qq.Score
		}
		got, essay := engine.GradeObjective(qq, engine.ParseOptionIDs(a.OptionIDs), sc)
		a.AutoScore = got
		a.IsEssay = essay
		if essay {
			hasEssay = true
		} else {
			obj += got
		}
		q.db.Model(a).Updates(map[string]any{"auto_score": a.AutoScore, "is_essay": a.IsEssay})
	}

	now := model.Time(timeutil.Now())
	status := model.SubGraded
	subj := 0.0
	if hasEssay {
		status = model.SubPendingManual
		for _, a := range sub.Answers {
			if a.IsEssay && a.ManualScore != nil {
				subj += *a.ManualScore
			}
		}
		if allEssayScored(sub.Answers) {
			status = model.SubGraded
		}
	}
	total := obj + subj
	pass := total+0.0001 >= exam.PassScore
	err := q.db.Model(&model.Submission{}).Where("id = ?", id).Updates(map[string]any{
		"status":           status,
		"objective_score":  obj,
		"subjective_score": subj,
		"total_score":      total,
		"pass":             pass,
		"graded_at":        now,
		"error_msg":        "",
	}).Error
	if err != nil {
		_ = q.db.Model(&model.Submission{}).Where("id = ?", id).Updates(map[string]any{
			"status": model.SubFailed, "error_msg": err.Error(),
		})
		return err
	}
	q.Metrics.Graded.Add(1)
	runAudit(q.db, sub)
	if status == model.SubGraded && pass {
		autoIssue(q.db, sub, total)
	}
	return nil
}

func allEssayScored(ans []model.SubmissionAnswer) bool {
	for _, a := range ans {
		if a.IsEssay && a.ManualScore == nil {
			return false
		}
	}
	return true
}

func runAudit(db *gorm.DB, sub model.Submission) {
	var traces []model.AnswerTrace
	var cheats []model.AntiCheatEvent
	db.Where("session_id = ?", sub.SessionID).Order("seq asc").Find(&traces)
	db.Where("session_id = ?", sub.SessionID).Find(&cheats)
	var answers []model.SubmissionAnswer
	db.Where("submission_id = ?", sub.ID).Find(&answers)
	qIDs := make([]uint64, 0, len(answers))
	for _, a := range answers {
		qIDs = append(qIDs, a.QuestionID)
	}
	var qs []model.Question
	db.Preload("Options").Where("id IN ?", qIDs).Find(&qs)
	qmap := map[uint64]model.Question{}
	for _, q := range qs {
		qmap[q.ID] = q
	}
	hits := engine.Analyze(traces, cheats, answers, qmap)
	now := model.Time(timeutil.Now())
	for _, h := range hits {
		flag := model.AuditFlag{
			TenantID:  sub.TenantID,
			SessionID: sub.SessionID,
			UserID:    sub.UserID,
			ExamID:    sub.ExamID,
			RuleCode:  h.Rule,
			Title:     h.Title,
			Evidence:  engine.EvidenceJSON(h.Evidence),
			Severity:  h.Severity,
			CreatedAt: now,
		}
		db.Clauses(clause.OnConflict{DoNothing: true}).Create(&flag)
	}
	if len(hits) > 0 {
		db.Model(&model.ExamSession{}).Where("id = ?", sub.SessionID).Updates(map[string]any{
			"suspicious": true,
			"integrity":  model.IntegritySuspicious,
		})
		db.Model(&model.Submission{}).Where("id = ?", sub.ID).Update("integrity", model.IntegritySuspicious)
	}
}

func autoIssue(db *gorm.DB, sub model.Submission, score float64) {
	var programs []model.CertProgram
	db.Where("tenant_id = ? AND require_exam = ?", sub.TenantID, sub.ExamID).Find(&programs)
	for _, p := range programs {
		if score+0.0001 < p.MinScore {
			continue
		}
		if p.RequireCourse > 0 {
			var chapters []model.Chapter
			db.Where("course_id = ? AND tenant_id = ?", p.RequireCourse, sub.TenantID).Find(&chapters)
			ok := true
			for _, ch := range chapters {
				var pr model.LearningProgress
				err := db.Where("user_id = ? AND chapter_id = ?", sub.UserID, ch.ID).First(&pr).Error
				if err != nil || !pr.Completed {
					ok = false
					break
				}
			}
			if !ok {
				continue
			}
		}
		var exist model.Certificate
		err := db.Where("program_id = ? AND user_id = ? AND status = ?", p.ID, sub.UserID, model.CertIssued).First(&exist).Error
		if err == nil {
			continue
		}
		now := timeutil.Now()
		cert := model.Certificate{
			TenantID:  sub.TenantID,
			ProgramID: p.ID,
			UserID:    sub.UserID,
			No:        fmt.Sprintf("QXT-%d-%d-%d", p.ID, sub.UserID, now.Unix()%100000),
			Status:    model.CertIssued,
			IssuedAt:  model.Time(now),
			ExpireAt:  model.Time(now.AddDate(0, 0, p.ValidDays)),
			Score:     score,
		}
		db.Create(&cert)
	}
}
