package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/miniqxt/backend/internal/logger"
)

func RequestLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		logger.Info("http",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"ms", time.Since(start).Milliseconds(),
		)
	}
}

func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, rec any) {
		logger.Error("panic", "err", rec, "path", c.Request.URL.Path)
		c.AbortWithStatusJSON(500, gin.H{"code": "INTERNAL", "message": "内部错误", "details": nil})
	})
}
