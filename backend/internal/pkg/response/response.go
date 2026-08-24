package response

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/miniqxt/backend/internal/pkg/apperr"
)

func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, gin.H{"data": data})
}

func Accepted(c *gin.Context, data any) {
	c.JSON(http.StatusAccepted, gin.H{"data": data})
}

func Page(c *gin.Context, items any, total int64, page, size int) {
	c.JSON(http.StatusOK, gin.H{"data": gin.H{
		"items":     items,
		"total":     total,
		"page":      page,
		"page_size": size,
	}})
}

func Fail(c *gin.Context, err error) {
	var ae *apperr.AppError
	if errors.As(err, &ae) {
		c.JSON(ae.HTTP, gin.H{"code": ae.Code, "message": ae.Message, "details": nil})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"code": "INTERNAL", "message": "内部错误", "details": nil})
}
