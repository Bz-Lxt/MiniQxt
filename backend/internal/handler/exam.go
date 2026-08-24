package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/miniqxt/backend/internal/middleware"
	"github.com/miniqxt/backend/internal/pkg/response"
	"github.com/miniqxt/backend/internal/service"
)

func (h *API) ListExams(c *gin.Context) {
	list, err := h.App.ListExams(middleware.TenantID(c))
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, list)
}

func (h *API) GetExam(c *gin.Context) {
	e, err := h.App.GetExam(middleware.TenantID(c), uid64(c, "id"))
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, e)
}

func (h *API) CreateExam(c *gin.Context) {
	var in service.ExamIn
	_ = c.ShouldBindJSON(&in)
	e, err := h.App.CreateExam(middleware.TenantID(c), in)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Created(c, e)
}

func (h *API) AssignExam(c *gin.Context) {
	var in struct {
		Scope string   `json:"scope"`
		IDs   []uint64 `json:"ids"`
	}
	_ = c.ShouldBindJSON(&in)
	if err := h.App.AssignExam(middleware.TenantID(c), uid64(c, "id"), in.Scope, in.IDs); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *API) MyExams(c *gin.Context) {
	list, err := h.App.MyExams(middleware.TenantID(c), middleware.UserID(c))
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, list)
}

func (h *API) StartExam(c *gin.Context) {
	out, err := h.App.StartExam(middleware.TenantID(c), middleware.UserID(c), uid64(c, "id"))
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, out)
}

func (h *API) GetSession(c *gin.Context) {
	out, err := h.App.GetSession(middleware.TenantID(c), middleware.UserID(c), uid64(c, "id"), staff(c))
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, out)
}

func (h *API) SessionPaper(c *gin.Context) {
	out, err := h.App.GetSession(middleware.TenantID(c), middleware.UserID(c), uid64(c, "id"), staff(c))
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{"paper": out["paper"], "shuffle_seed": out["shuffle_seed"], "remaining_sec": out["remaining_sec"]})
}

func (h *API) Heartbeat(c *gin.Context) {
	if err := h.App.Heartbeat(middleware.TenantID(c), middleware.UserID(c), uid64(c, "id")); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *API) AntiCheat(c *gin.Context) {
	var in service.AntiCheatIn
	_ = c.ShouldBindJSON(&in)
	out, err := h.App.ReportAntiCheat(middleware.TenantID(c), middleware.UserID(c), uid64(c, "id"), in)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, out)
}

func (h *API) Traces(c *gin.Context) {
	var in struct {
		Events []service.TraceIn `json:"events"`
	}
	_ = c.ShouldBindJSON(&in)
	if err := h.App.IngestTraces(middleware.TenantID(c), middleware.UserID(c), uid64(c, "id"), in.Events); err != nil {
		response.Fail(c, err)
		return
	}
	response.Accepted(c, gin.H{"accepted": len(in.Events)})
}

func (h *API) Submit(c *gin.Context) {
	var in struct {
		Answers []service.AnswerIn `json:"answers"`
	}
	_ = c.ShouldBindJSON(&in)
	out, err := h.App.SubmitExam(middleware.TenantID(c), middleware.UserID(c), uid64(c, "id"), in.Answers)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Accepted(c, out)
}

func (h *API) GetSubmission(c *gin.Context) {
	out, err := h.App.GetSubmission(middleware.TenantID(c), middleware.UserID(c), uid64(c, "id"), staff(c))
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, out)
}

func (h *API) ListSubmissions(c *gin.Context) {
	p, s := pageSize(c)
	list, total, err := h.App.ListSubmissions(middleware.TenantID(c), c.Query("status"), p, s)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Page(c, list, total, p, s)
}

func (h *API) ManualScore(c *gin.Context) {
	var in struct {
		Items []service.ManualItem `json:"items"`
	}
	_ = c.ShouldBindJSON(&in)
	sub, err := h.App.ManualScore(middleware.TenantID(c), uid64(c, "id"), in.Items)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, sub)
}

func (h *API) Timeline(c *gin.Context) {
	out, err := h.App.Timeline(middleware.TenantID(c), uid64(c, "id"))
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, out)
}

func (h *API) RunAudit(c *gin.Context) {
	flags, err := h.App.RunAudit(middleware.TenantID(c), uid64(c, "id"))
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, flags)
}

func (h *API) ListFlags(c *gin.Context) {
	list, err := h.App.ListFlags(middleware.TenantID(c))
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, list)
}

func (h *API) ListPrograms(c *gin.Context) {
	list, err := h.App.ListPrograms(middleware.TenantID(c))
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, list)
}

func (h *API) CreateProgram(c *gin.Context) {
	var in service.ProgramIn
	_ = c.ShouldBindJSON(&in)
	p, err := h.App.CreateProgram(middleware.TenantID(c), in)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Created(c, p)
}

func (h *API) PatchProgram(c *gin.Context) {
	var in service.ProgramIn
	_ = c.ShouldBindJSON(&in)
	if err := h.App.PatchProgram(middleware.TenantID(c), uid64(c, "id"), in); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *API) EvaluateProgram(c *gin.Context) {
	n, err := h.App.EvaluateProgram(middleware.TenantID(c), uid64(c, "id"))
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{"issued": n})
}

func (h *API) ListCertificates(c *gin.Context) {
	mine := c.FullPath() == "/api/v1/my/certificates"
	list, err := h.App.ListCertificates(middleware.TenantID(c), middleware.UserID(c), mine)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, list)
}
