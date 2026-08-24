package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/miniqxt/backend/internal/middleware"
	"github.com/miniqxt/backend/internal/pkg/apperr"
	"github.com/miniqxt/backend/internal/pkg/response"
	"github.com/miniqxt/backend/internal/service"
)

func (h *API) ListTenants(c *gin.Context) {
	list, err := h.App.ListTenants()
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, list)
}

func (h *API) CreateTenant(c *gin.Context) {
	var in struct{ Name, Code string }
	_ = c.ShouldBindJSON(&in)
	t, err := h.App.CreateTenant(in.Name, in.Code)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Created(c, t)
}

func (h *API) PatchTenant(c *gin.Context) {
	var in service.TenantPatch
	_ = c.ShouldBindJSON(&in)
	t, err := h.App.PatchTenant(uid64(c, "id"), in)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, t)
}

func (h *API) ListDepts(c *gin.Context) {
	list, err := h.App.ListDepts(middleware.TenantID(c))
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, list)
}

func (h *API) CreateDept(c *gin.Context) {
	var in struct {
		Name     string `json:"name"`
		ParentID uint64 `json:"parent_id"`
	}
	_ = c.ShouldBindJSON(&in)
	d, err := h.App.CreateDept(middleware.TenantID(c), in.Name, in.ParentID)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Created(c, d)
}

func (h *API) PatchDept(c *gin.Context) {
	var in struct{ Name string `json:"name"` }
	_ = c.ShouldBindJSON(&in)
	if err := h.App.PatchDept(middleware.TenantID(c), uid64(c, "id"), in.Name); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *API) DeleteDept(c *gin.Context) {
	if err := h.App.DeleteDept(middleware.TenantID(c), uid64(c, "id")); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *API) ListPositions(c *gin.Context) {
	list, err := h.App.ListPositions(middleware.TenantID(c))
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, list)
}

func (h *API) CreatePosition(c *gin.Context) {
	var in struct {
		Name  string `json:"name"`
		Level int    `json:"level"`
	}
	_ = c.ShouldBindJSON(&in)
	p, err := h.App.CreatePosition(middleware.TenantID(c), in.Name, in.Level)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Created(c, p)
}

func (h *API) PatchPosition(c *gin.Context) {
	var in struct {
		Name  string `json:"name"`
		Level int    `json:"level"`
	}
	_ = c.ShouldBindJSON(&in)
	if err := h.App.PatchPosition(middleware.TenantID(c), uid64(c, "id"), in.Name, in.Level); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *API) DeletePosition(c *gin.Context) {
	if err := h.App.DeletePosition(middleware.TenantID(c), uid64(c, "id")); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *API) ListEmployees(c *gin.Context) {
	p, s := pageSize(c)
	list, total, err := h.App.ListEmployees(middleware.TenantID(c), p, s, c.Query("q"))
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Page(c, list, total, p, s)
}

func (h *API) CreateEmployee(c *gin.Context) {
	var in service.EmployeeIn
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, apperr.Validation)
		return
	}
	u, err := h.App.CreateEmployee(middleware.TenantID(c), in)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Created(c, u.Public())
}

func (h *API) PatchEmployee(c *gin.Context) {
	var in service.EmployeeIn
	_ = c.ShouldBindJSON(&in)
	if err := h.App.PatchEmployee(middleware.TenantID(c), uid64(c, "id"), in); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *API) DeleteEmployee(c *gin.Context) {
	if err := h.App.DeleteEmployee(middleware.TenantID(c), uid64(c, "id")); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}
