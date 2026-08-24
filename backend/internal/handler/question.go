package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/miniqxt/backend/internal/middleware"
	"github.com/miniqxt/backend/internal/pkg/response"
	"github.com/miniqxt/backend/internal/service"
)

func (h *API) ListQuestions(c *gin.Context) {
	p, s := pageSize(c)
	list, total, err := h.App.ListQuestions(middleware.TenantID(c), c.Query("type"), c.Query("tag"), c.Query("q"), p, s)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Page(c, list, total, p, s)
}

func (h *API) GetQuestion(c *gin.Context) {
	q, err := h.App.GetQuestion(middleware.TenantID(c), uid64(c, "id"))
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, q)
}

func (h *API) CreateQuestion(c *gin.Context) {
	var in service.QuestionIn
	_ = c.ShouldBindJSON(&in)
	q, err := h.App.CreateQuestion(middleware.TenantID(c), in)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Created(c, q)
}

func (h *API) UpdateQuestion(c *gin.Context) {
	var in service.QuestionIn
	_ = c.ShouldBindJSON(&in)
	q, err := h.App.UpdateQuestion(middleware.TenantID(c), uid64(c, "id"), in)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, q)
}

func (h *API) DeleteQuestion(c *gin.Context) {
	if err := h.App.DeleteQuestion(middleware.TenantID(c), uid64(c, "id")); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *API) ListPapers(c *gin.Context) {
	list, err := h.App.ListPapers(middleware.TenantID(c))
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, list)
}

func (h *API) GetPaper(c *gin.Context) {
	p, err := h.App.GetPaper(middleware.TenantID(c), uid64(c, "id"))
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, p)
}

func (h *API) CreatePaper(c *gin.Context) {
	var in struct {
		Title            string `json:"title"`
		ShuffleQuestions bool   `json:"shuffle_questions"`
		ShuffleOptions   bool   `json:"shuffle_options"`
	}
	_ = c.ShouldBindJSON(&in)
	p, err := h.App.CreatePaper(middleware.TenantID(c), in.Title, in.ShuffleQuestions, in.ShuffleOptions)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Created(c, p)
}

func (h *API) SavePaperItems(c *gin.Context) {
	var in struct {
		Title            string               `json:"title"`
		ShuffleQuestions bool                 `json:"shuffle_questions"`
		ShuffleOptions   bool                 `json:"shuffle_options"`
		Items            []service.PaperItemIn `json:"items"`
	}
	_ = c.ShouldBindJSON(&in)
	p, err := h.App.SavePaperItems(middleware.TenantID(c), uid64(c, "id"), in.Title, in.ShuffleQuestions, in.ShuffleOptions, in.Items)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, p)
}

func (h *API) DeletePaper(c *gin.Context) {
	if err := h.App.DeletePaper(middleware.TenantID(c), uid64(c, "id")); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *API) PreviewPaper(c *gin.Context) {
	p, err := h.App.GetPaper(middleware.TenantID(c), uid64(c, "id"))
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, h.App.RenderPaper(p, 42, false))
}
