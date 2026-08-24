package service

import (
	"crypto/rand"
	"encoding/binary"
	"time"

	"github.com/miniqxt/backend/internal/model"
	"github.com/miniqxt/backend/internal/pkg/apperr"
	"github.com/miniqxt/backend/internal/repo"
	"github.com/miniqxt/backend/internal/timeutil"
	"gorm.io/gorm"
)

type ExamIn struct {
	PaperID     uint64  `json:"paper_id"`
	Title       string  `json:"title"`
	StartAt     string  `json:"start_at"`
	EndAt       string  `json:"end_at"`
	DurationSec int     `json:"duration_sec"`
	PassScore   float64 `json:"pass_score"`
	MaxAttempts int     `json:"max_attempts"`
}

func (a *App) ListExams(tid uint64) ([]model.Exam, error) {
	var list []model.Exam
	err := repo.Tenant(a.DB, tid).Preload("Paper").Order("id desc").Find(&list).Error
	return list, err
}

func (a *App) GetExam(tid, id uint64) (model.Exam, error) {
	var e model.Exam
	if err := repo.Tenant(a.DB, tid).Preload("Paper").First(&e, id).Error; err != nil {
		return e, apperr.NotFound
	}
	return e, nil
}

func (a *App) CreateExam(tid uint64, in ExamIn) (model.Exam, error) {
	if in.Title == "" || in.PaperID == 0 {
		return model.Exam{}, apperr.Validation.With("考试标题与试卷必填")
	}
	if _, err := a.GetPaper(tid, in.PaperID); err != nil {
		return model.Exam{}, err
	}
	st, err1 := timeutil.ParseRFC3339(in.StartAt)
	et, err2 := timeutil.ParseRFC3339(in.EndAt)
	if err1 != nil || err2 != nil {
		return model.Exam{}, apperr.Validation.With("时间格式错误")
	}
	if in.DurationSec < 60 {
		in.DurationSec = 3600
	}
	if in.MaxAttempts < 1 {
		in.MaxAttempts = 1
	}
	e := model.Exam{
		TenantID: tid, PaperID: in.PaperID, Title: in.Title,
		StartAt: model.Time(st), EndAt: model.Time(et), DurationSec: in.DurationSec,
		PassScore: in.PassScore, MaxAttempts: in.MaxAttempts, Status: "published",
		CreatedAt: model.Time(timeutil.Now()),
	}
	a.DB.Create(&e)
	return e, nil
}

func (a *App) AssignExam(tid, examID uint64, scope string, ids []uint64) error {
	if _, err := a.GetExam(tid, examID); err != nil {
		return err
	}
	if scope == "all" {
		return a.DB.Create(&model.ExamAssignment{TenantID: tid, ExamID: examID, Scope: "all", CreatedAt: model.Time(timeutil.Now())}).Error
	}
	for _, id := range ids {
		a.DB.Create(&model.ExamAssignment{TenantID: tid, ExamID: examID, Scope: scope, TargetID: id, CreatedAt: model.Time(timeutil.Now())})
	}
	return nil
}

func (a *App) MyExams(tid, uid uint64) ([]map[string]any, error) {
	var exams []model.Exam
	repo.Tenant(a.DB, tid).Find(&exams)
	var u model.User
	a.DB.First(&u, uid)
	out := make([]map[string]any, 0)
	for _, e := range exams {
		if !a.assigned(tid, e.ID, u) {
			continue
		}
		var sess model.ExamSession
		_ = a.DB.Where("exam_id = ? AND user_id = ?", e.ID, uid).Order("attempt desc").First(&sess).Error
		out = append(out, map[string]any{
			"id": e.ID, "title": e.Title, "start_at": e.StartAt, "end_at": e.EndAt,
			"duration_sec": e.DurationSec, "pass_score": e.PassScore,
			"session": sess, "remaining_sec": remaining(sess),
		})
	}
	return out, nil
}

func (a *App) assigned(tid, examID uint64, u model.User) bool {
	var as []model.ExamAssignment
	a.DB.Where("tenant_id = ? AND exam_id = ?", tid, examID).Find(&as)
	if len(as) == 0 {
		return true
	}
	for _, x := range as {
		if x.Scope == "all" {
			return true
		}
		if x.Scope == "user" && x.TargetID == u.ID {
			return true
		}
		if x.Scope == "dept" && u.DeptID != nil && x.TargetID == *u.DeptID {
			return true
		}
	}
	return false
}

func remaining(s model.ExamSession) int {
	if s.ID == 0 || s.Status != model.SessInProgress {
		return 0
	}
	sec := int(s.EndsAt.Time().Sub(timeutil.Now()).Seconds())
	if sec < 0 {
		return 0
	}
	return sec
}

func randSeed() int64 {
	var b [8]byte
	_, _ = rand.Read(b[:])
	n := int64(binary.LittleEndian.Uint64(b[:]))
	if n < 0 {
		n = -n
	}
	if n == 0 {
		n = time.Now().UnixNano()
	}
	return n
}

func (a *App) StartExam(tid, uid, examID uint64) (map[string]any, error) {
	e, err := a.GetExam(tid, examID)
	if err != nil {
		return nil, err
	}
	now := timeutil.Now()
	if now.Before(e.StartAt.Time()) || now.After(e.EndAt.Time()) {
		return nil, apperr.ExamNotOpen
	}
	var u model.User
	a.DB.First(&u, uid)
	if !a.assigned(tid, examID, u) {
		return nil, apperr.Forbidden
	}
	var exist model.ExamSession
	err = a.DB.Where("exam_id = ? AND user_id = ? AND status = ?", examID, uid, model.SessInProgress).First(&exist).Error
	if err == nil {
		if timeutil.Now().After(exist.EndsAt.Time()) {
			a.DB.Model(&exist).Update("status", model.SessExpired)
			return nil, apperr.SessionGone
		}
		return a.sessionPayload(tid, exist)
	}
	var cnt int64
	a.DB.Model(&model.ExamSession{}).Where("exam_id = ? AND user_id = ?", examID, uid).Count(&cnt)
	if int(cnt) >= e.MaxAttempts {
		return nil, apperr.Conflict.With("已超过最大考试次数")
	}
	sess := model.ExamSession{
		TenantID: tid, ExamID: examID, UserID: uid, Attempt: int(cnt) + 1,
		ShuffleSeed: randSeed(), StartedAt: model.Time(now),
		EndsAt: model.Time(now.Add(time.Duration(e.DurationSec) * time.Second)),
		Status: model.SessInProgress, Integrity: model.IntegrityOK,
		CreatedAt: model.Time(now),
	}
	if err := a.DB.Create(&sess).Error; err != nil {
		return nil, apperr.Conflict.With("考试会话冲突，请重试")
	}
	return a.sessionPayload(tid, sess)
}

func (a *App) GetSession(tid, uid, sid uint64, staff bool) (map[string]any, error) {
	var s model.ExamSession
	tx := repo.Tenant(a.DB, tid).Where("id = ?", sid)
	if !staff {
		tx = tx.Where("user_id = ?", uid)
	}
	if err := tx.First(&s).Error; err != nil {
		return nil, apperr.NotFound
	}
	if s.Status == model.SessInProgress && timeutil.Now().After(s.EndsAt.Time()) {
		s.Status = model.SessExpired
		a.DB.Model(&s).Update("status", model.SessExpired)
	}
	return a.sessionPayload(tid, s)
}

func (a *App) sessionPayload(tid uint64, s model.ExamSession) (map[string]any, error) {
	e, err := a.GetExam(tid, s.ExamID)
	if err != nil {
		return nil, err
	}
	p, err := a.GetPaper(tid, e.PaperID)
	if err != nil {
		return nil, err
	}
	hide := s.Status == model.SessInProgress
	paper := a.RenderPaper(p, s.ShuffleSeed, hide)
	return map[string]any{
		"session":       s,
		"remaining_sec": remaining(s),
		"exam":          e,
		"paper":         paper,
		"shuffle_seed":  s.ShuffleSeed,
	}, nil
}

func (a *App) Heartbeat(tid, uid, sid uint64) error {
	res := repo.Tenant(a.DB, tid).Model(&model.ExamSession{}).
		Where("id = ? AND user_id = ? AND status = ?", sid, uid, model.SessInProgress).
		Update("heartbeat_n", gorm.Expr("heartbeat_n + 1"))
	if res.RowsAffected == 0 {
		return apperr.NotFound
	}
	return nil
}

type AntiCheatIn struct {
	EventType  string `json:"event_type"`
	OccurredAt string `json:"occurred_at"`
	DurationMS int    `json:"duration_ms"`
}

func (a *App) ReportAntiCheat(tid, uid, sid uint64, in AntiCheatIn) (map[string]any, error) {
	var s model.ExamSession
	if err := repo.Tenant(a.DB, tid).Where("id = ? AND user_id = ?", sid, uid).First(&s).Error; err != nil {
		return nil, apperr.NotFound
	}
	if s.Status != model.SessInProgress {
		return nil, apperr.SessionGone
	}
	ot, err := timeutil.ParseRFC3339(in.OccurredAt)
	if err != nil {
		ot = timeutil.Now()
	}
	if in.EventType == "" {
		in.EventType = "blur"
	}
	ev := model.AntiCheatEvent{
		TenantID: tid, SessionID: sid, EventType: in.EventType,
		OccurredAt: model.Time(ot), DurationMS: in.DurationMS, CreatedAt: model.Time(timeutil.Now()),
	}
	a.DB.Create(&ev)
	if in.EventType == "blur" || in.EventType == "hidden" {
		s.BlurCount++
	}
	ten, _ := a.TenantOf(tid)
	force := ten.BlurForce > 0 && s.BlurCount >= ten.BlurForce
	susp := ten.BlurWarn > 0 && s.BlurCount >= ten.BlurWarn
	up := map[string]any{"blur_count": s.BlurCount, "suspicious": susp}
	if susp {
		up["integrity"] = model.IntegritySuspicious
	}
	if force {
		up["status"] = model.SessForced
		up["force_reason"] = "blur_threshold"
	}
	a.DB.Model(&s).Updates(up)
	return map[string]any{
		"blur_count": s.BlurCount, "suspicious": susp, "force_submit": force,
		"warn_at": ten.BlurWarn, "force_at": ten.BlurForce,
	}, nil
}
