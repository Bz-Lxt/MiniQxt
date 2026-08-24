package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/miniqxt/backend/internal/model"
	"github.com/miniqxt/backend/internal/pkg/apperr"
	"github.com/miniqxt/backend/internal/pkg/response"
)

func Roles(roles ...string) gin.HandlerFunc {
	allow := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		allow[r] = struct{}{}
	}
	return func(c *gin.Context) {
		r := Role(c)
		if _, ok := allow[r]; !ok {
			response.Fail(c, apperr.Forbidden)
			c.Abort()
			return
		}
		c.Next()
	}
}

func Staff() gin.HandlerFunc {
	return Roles(model.RolePlatform, model.RoleTenant, model.RoleInstructor)
}

func Admins() gin.HandlerFunc {
	return Roles(model.RolePlatform, model.RoleTenant)
}
