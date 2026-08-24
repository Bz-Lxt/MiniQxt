package service

import (
	"github.com/miniqxt/backend/internal/model"
	"github.com/miniqxt/backend/internal/pkg/apperr"
	"github.com/miniqxt/backend/internal/repo"
	"github.com/miniqxt/backend/internal/timeutil"
)

func (a *App) ListCourses(tid uint64) ([]model.Course, error) {
	var list []model.Course
	err := repo.Tenant(a.DB, tid).Preload("Chapters").Order("id desc").Find(&list).Error
	return list, err
}

func (a *App) GetCourse(tid, id uint64) (model.Course, error) {
	var c model.Course
	err := repo.Tenant(a.DB, tid).Preload("Chapters").First(&c, id).Error
	if err != nil {
		return c, apperr.NotFound
	}
	return c, nil
}

func (a *App) CreateCourse(tid uint64, title, summary string, required bool) (model.Course, error) {
	if title == "" {
		return model.Course{}, apperr.Validation.With("课程标题必填")
	}
	c := model.Course{TenantID: tid, Title: title, Summary: summary, Required: required, CreatedAt: model.Time(timeutil.Now())}
	a.DB.Create(&c)
	return c, nil
}

func (a *App) PatchCourse(tid, id uint64, title, summary string, required *bool) error {
	up := map[string]any{}
	if title != "" {
		up["title"] = title
	}
	if summary != "" {
		up["summary"] = summary
	}
	if required != nil {
		up["required"] = *required
	}
	res := repo.Tenant(a.DB, tid).Model(&model.Course{}).Where("id = ?", id).Updates(up)
	if res.RowsAffected == 0 {
		return apperr.NotFound
	}
	return nil
}

func (a *App) DeleteCourse(tid, id uint64) error {
	res := repo.Tenant(a.DB, tid).Delete(&model.Course{}, id)
	if res.RowsAffected == 0 {
		return apperr.NotFound
	}
	return nil
}

func (a *App) AddChapter(tid, courseID uint64, title, video string, duration, sort int) (model.Chapter, error) {
	var c model.Course
	if err := repo.Tenant(a.DB, tid).First(&c, courseID).Error; err != nil {
		return model.Chapter{}, apperr.NotFound
	}
	ch := model.Chapter{
		TenantID: tid, CourseID: courseID, Title: title, VideoFile: video,
		DurationSec: duration, SortNo: sort, CreatedAt: model.Time(timeutil.Now()),
	}
	a.DB.Create(&ch)
	return ch, nil
}

func (a *App) PatchChapter(tid, id uint64, title string, duration, sort int) error {
	res := repo.Tenant(a.DB, tid).Model(&model.Chapter{}).Where("id = ?", id).Updates(map[string]any{
		"title": title, "duration_sec": duration, "sort_no": sort,
	})
	if res.RowsAffected == 0 {
		return apperr.NotFound
	}
	return nil
}

func (a *App) DeleteChapter(tid, id uint64) error {
	res := repo.Tenant(a.DB, tid).Delete(&model.Chapter{}, id)
	if res.RowsAffected == 0 {
		return apperr.NotFound
	}
	return nil
}

func (a *App) AssignCourse(tid, courseID uint64, scope string, target uint64) error {
	var c model.Course
	if err := repo.Tenant(a.DB, tid).First(&c, courseID).Error; err != nil {
		return apperr.NotFound
	}
	as := model.CourseAssignment{TenantID: tid, CourseID: courseID, Scope: scope, TargetID: target, CreatedAt: model.Time(timeutil.Now())}
	return a.DB.Create(&as).Error
}

func (a *App) MyCourses(tid, uid uint64) ([]map[string]any, error) {
	var courses []model.Course
	repo.Tenant(a.DB, tid).Preload("Chapters").Find(&courses)
	out := make([]map[string]any, 0, len(courses))
	for _, c := range courses {
		var progresses []model.LearningProgress
		a.DB.Where("user_id = ? AND course_id = ?", uid, c.ID).Find(&progresses)
		pm := map[uint64]model.LearningProgress{}
		sum := 0.0
		for _, p := range progresses {
			pm[p.ChapterID] = p
			sum += p.Percent
		}
		pct := 0.0
		if n := len(c.Chapters); n > 0 {
			pct = sum / float64(n)
		}
		chs := make([]map[string]any, 0, len(c.Chapters))
		for _, ch := range c.Chapters {
			p := pm[ch.ID]
			chs = append(chs, map[string]any{
				"id": ch.ID, "title": ch.Title, "video_file": ch.VideoFile,
				"duration_sec": ch.DurationSec, "percent": p.Percent, "completed": p.Completed,
				"position_sec": p.PositionSec,
			})
		}
		out = append(out, map[string]any{
			"id": c.ID, "title": c.Title, "summary": c.Summary,
			"required": c.Required, "percent": pct, "chapters": chs,
		})
	}
	return out, nil
}

func (a *App) ReportProgress(tid, uid, chapterID uint64, delta, position int) (model.LearningProgress, error) {
	var ch model.Chapter
	if err := repo.Tenant(a.DB, tid).First(&ch, chapterID).Error; err != nil {
		return model.LearningProgress{}, apperr.NotFound
	}
	if delta < 0 {
		delta = 0
	}
	if delta > 15 {
		delta = 15 // anti-seek: one report at most 15s of valid watch
	}
	now := timeutil.Now()
	var p model.LearningProgress
	err := a.DB.Where("user_id = ? AND chapter_id = ?", uid, chapterID).First(&p).Error
	if err != nil {
		p = model.LearningProgress{
			TenantID: tid, UserID: uid, ChapterID: chapterID, CourseID: ch.CourseID,
			LastReportAt: model.Time(now),
		}
	} else {
		elapsed := now.Sub(p.LastReportAt.Time()).Seconds()
		if elapsed < float64(delta) {
			delta = int(elapsed)
			if delta < 0 {
				delta = 0
			}
		}
	}
	p.WatchedSec += delta
	p.PositionSec = position
	if ch.DurationSec > 0 {
		p.Percent = float64(p.WatchedSec) * 100 / float64(ch.DurationSec)
		if p.Percent > 100 {
			p.Percent = 100
		}
	}
	p.Completed = p.Percent >= 90
	p.LastReportAt = model.Time(now)
	if p.ID == 0 {
		a.DB.Create(&p)
	} else {
		a.DB.Save(&p)
	}
	return p, nil
}
