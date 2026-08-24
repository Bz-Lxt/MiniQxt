package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/miniqxt/backend/internal/middleware"
	"github.com/miniqxt/backend/internal/pkg/response"
)

func (h *API) Login(c *gin.Context) {
	var in struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, err)
		return
	}
	tok, user, err := h.App.Login(in.Username, in.Password)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{"token": tok, "user": user})
}

func (h *API) Me(c *gin.Context) {
	u, err := h.App.Me(middleware.UserID(c))
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, u.Public())
}

func (h *API) Health(c *gin.Context) {
	depth := 0
	workers := 0
	if h.App.Grader != nil {
		depth = h.App.Grader.Depth()
		workers = h.App.Grader.Workers()
	}
	response.OK(c, gin.H{"status": "ok", "queue_depth": depth, "workers": workers})
}

func (h *API) Metrics(c *gin.Context) {
	depth, workers := 0, 0
	if h.App.Grader != nil {
		depth = h.App.Grader.Depth()
		workers = h.App.Grader.Workers()
	}
	response.OK(c, h.App.Met.Snapshot(depth, workers))
}
