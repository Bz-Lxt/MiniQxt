package handler

import (
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/miniqxt/backend/internal/middleware"
	"github.com/miniqxt/backend/internal/pkg/apperr"
	"github.com/miniqxt/backend/internal/pkg/response"
)

func (h *API) ListCourses(c *gin.Context) {
	list, err := h.App.ListCourses(middleware.TenantID(c))
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, list)
}

func (h *API) GetCourse(c *gin.Context) {
	co, err := h.App.GetCourse(middleware.TenantID(c), uid64(c, "id"))
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, co)
}

func (h *API) CreateCourse(c *gin.Context) {
	var in struct {
		Title    string `json:"title"`
		Summary  string `json:"summary"`
		Required bool   `json:"required"`
	}
	_ = c.ShouldBindJSON(&in)
	co, err := h.App.CreateCourse(middleware.TenantID(c), in.Title, in.Summary, in.Required)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Created(c, co)
}

func (h *API) PatchCourse(c *gin.Context) {
	var in struct {
		Title    string `json:"title"`
		Summary  string `json:"summary"`
		Required *bool  `json:"required"`
	}
	_ = c.ShouldBindJSON(&in)
	if err := h.App.PatchCourse(middleware.TenantID(c), uid64(c, "id"), in.Title, in.Summary, in.Required); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *API) DeleteCourse(c *gin.Context) {
	if err := h.App.DeleteCourse(middleware.TenantID(c), uid64(c, "id")); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *API) AddChapter(c *gin.Context) {
	var in struct {
		Title       string `json:"title"`
		VideoFile   string `json:"video_file"`
		DurationSec int    `json:"duration_sec"`
		SortNo      int    `json:"sort_no"`
	}
	_ = c.ShouldBindJSON(&in)
	ch, err := h.App.AddChapter(middleware.TenantID(c), uid64(c, "id"), in.Title, in.VideoFile, in.DurationSec, in.SortNo)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Created(c, ch)
}

func (h *API) PatchChapter(c *gin.Context) {
	var in struct {
		Title       string `json:"title"`
		DurationSec int    `json:"duration_sec"`
		SortNo      int    `json:"sort_no"`
	}
	_ = c.ShouldBindJSON(&in)
	if err := h.App.PatchChapter(middleware.TenantID(c), uid64(c, "id"), in.Title, in.DurationSec, in.SortNo); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *API) DeleteChapter(c *gin.Context) {
	if err := h.App.DeleteChapter(middleware.TenantID(c), uid64(c, "id")); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *API) AssignCourse(c *gin.Context) {
	var in struct {
		Scope    string `json:"scope"`
		TargetID uint64 `json:"target_id"`
	}
	_ = c.ShouldBindJSON(&in)
	if err := h.App.AssignCourse(middleware.TenantID(c), uid64(c, "id"), in.Scope, in.TargetID); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *API) MyCourses(c *gin.Context) {
	list, err := h.App.MyCourses(middleware.TenantID(c), middleware.UserID(c))
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, list)
}

func (h *API) ReportProgress(c *gin.Context) {
	var in struct {
		DeltaSeconds int `json:"delta_seconds"`
		PositionSec  int `json:"position_sec"`
	}
	_ = c.ShouldBindJSON(&in)
	p, err := h.App.ReportProgress(middleware.TenantID(c), middleware.UserID(c), uid64(c, "id"), in.DeltaSeconds, in.PositionSec)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, p)
}

func (h *API) Media(c *gin.Context) {
	name := filepath.Base(c.Param("name"))
	if name == "" || strings.Contains(name, "..") {
		response.Fail(c, apperr.Validation)
		return
	}
	c.File(filepath.Join(h.App.Cfg.MediaDir, name))
}
