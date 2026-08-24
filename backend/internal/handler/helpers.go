package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/miniqxt/backend/internal/middleware"
	"github.com/miniqxt/backend/internal/model"
	"github.com/miniqxt/backend/internal/service"
)

type API struct {
	App *service.App
}

func uid64(c *gin.Context, name string) uint64 {
	n, _ := strconv.ParseUint(c.Param(name), 10, 64)
	return n
}

func queryU64(c *gin.Context, name string) uint64 {
	n, _ := strconv.ParseUint(c.Query(name), 10, 64)
	return n
}

func pageSize(c *gin.Context) (int, int) {
	p, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	s, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	return p, s
}

func staff(c *gin.Context) bool {
	r := middleware.Role(c)
	return r == model.RolePlatform || r == model.RoleTenant || r == model.RoleInstructor
}
