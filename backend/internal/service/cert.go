package service

import (
	"fmt"

	"github.com/miniqxt/backend/internal/model"
	"github.com/miniqxt/backend/internal/pkg/apperr"
	"github.com/miniqxt/backend/internal/repo"
	"github.com/miniqxt/backend/internal/timeutil"
)

type ProgramIn struct {
	Name          string  `json:"name"`
	PositionID    uint64  `json:"position_id"`
	Level         int     `json:"level"`
	RequireCourse uint64  `json:"require_course"`
	RequireExam   uint64  `json:"require_exam"`
	MinScore      float64 `json:"min_score"`
	ValidDays     int     `json:"valid_days"`
}

func (a *App) ListPrograms(tid uint64) ([]model.CertProgram, error) {
	var list []model.CertProgram
	err := repo.Tenant(a.DB, tid).Order("id desc").Find(&list).Error
	return list, err
}

func (a *App) CreateProgram(tid uint64, in ProgramIn) (model.CertProgram, error) {
	if in.Name == "" {
		return model.CertProgram{}, apperr.Validation.With("认证名称必填")
	}
	if in.ValidDays <= 0 {
		in.ValidDays = 365
	}
	p := model.CertProgram{
		TenantID: tid, Name: in.Name, PositionID: in.PositionID, Level: in.Level,
		RequireCourse: in.RequireCourse, RequireExam: in.RequireExam,
		MinScore: in.MinScore, ValidDays: in.ValidDays, CreatedAt: model.Time(timeutil.Now()),
	}
	a.DB.Create(&p)
	return p, nil
}

func (a *App) PatchProgram(tid, id uint64, in ProgramIn) error {
	res := repo.Tenant(a.DB, tid).Model(&model.CertProgram{}).Where("id = ?", id).Updates(map[string]any{
		"name": in.Name, "position_id": in.PositionID, "level": in.Level,
		"require_course": in.RequireCourse, "require_exam": in.RequireExam,
		"min_score": in.MinScore, "valid_days": in.ValidDays,
	})
	if res.RowsAffected == 0 {
		return apperr.NotFound
	}
	return nil
}

func (a *App) ListCertificates(tid uint64, uid uint64, mine bool) ([]model.Certificate, error) {
	tx := repo.Tenant(a.DB, tid).Preload("Program")
	if mine {
		tx = tx.Where("user_id = ?", uid)
	}
	var list []model.Certificate
	err := tx.Order("id desc").Find(&list).Error
	return list, err
}

func (a *App) EvaluateProgram(tid, pid uint64) (int, error) {
	var p model.CertProgram
	if err := repo.Tenant(a.DB, tid).First(&p, pid).Error; err != nil {
		return 0, apperr.NotFound
	}
	var subs []model.Submission
	a.DB.Where("tenant_id = ? AND exam_id = ? AND status = ? AND pass = ?", tid, p.RequireExam, model.SubGraded, true).Find(&subs)
	n := 0
	for _, s := range subs {
		if s.TotalScore+0.0001 < p.MinScore {
			continue
		}
		var exist model.Certificate
		if a.DB.Where("program_id = ? AND user_id = ? AND status = ?", p.ID, s.UserID, model.CertIssued).First(&exist).Error == nil {
			continue
		}
		now := timeutil.Now()
		c := model.Certificate{
			TenantID: tid, ProgramID: p.ID, UserID: s.UserID,
			No: fmt.Sprintf("QXT-%d-%d-%d", p.ID, s.UserID, now.Unix()%100000),
			Status: model.CertIssued, IssuedAt: model.Time(now),
			ExpireAt: model.Time(now.AddDate(0, 0, p.ValidDays)), Score: s.TotalScore,
		}
		a.DB.Create(&c)
		n++
	}
	return n, nil
}
