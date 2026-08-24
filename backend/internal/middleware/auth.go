package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/miniqxt/backend/internal/pkg/apperr"
	"github.com/miniqxt/backend/internal/pkg/jwtutil"
	"github.com/miniqxt/backend/internal/pkg/response"
)

const CtxClaims = "claims"

func Auth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if h == "" || !strings.HasPrefix(h, "Bearer ") {
			response.Fail(c, apperr.Unauthorized)
			c.Abort()
			return
		}
		token := strings.TrimPrefix(h, "Bearer ")
		claims, err := jwtutil.Parse(secret, token)
		if err != nil {
			response.Fail(c, apperr.Unauthorized)
			c.Abort()
			return
		}
		c.Set(CtxClaims, claims)
		c.Next()
	}
}

func Claims(c *gin.Context) *jwtutil.Claims {
	v, ok := c.Get(CtxClaims)
	if !ok {
		return nil
	}
	cl, _ := v.(*jwtutil.Claims)
	return cl
}

func TenantID(c *gin.Context) uint64 {
	cl := Claims(c)
	if cl == nil {
		return 0
	}
	return cl.TenantID
}

func UserID(c *gin.Context) uint64 {
	cl := Claims(c)
	if cl == nil {
		return 0
	}
	return cl.UserID
}

func Role(c *gin.Context) string {
	cl := Claims(c)
	if cl == nil {
		return ""
	}
	return cl.Role
}
