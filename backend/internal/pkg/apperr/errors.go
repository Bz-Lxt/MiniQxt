package apperr

import "fmt"

type AppError struct {
	HTTP    int
	Code    string
	Message string
}

func (e *AppError) Error() string { return e.Message }

func New(http int, code, msg string) *AppError {
	return &AppError{HTTP: http, Code: code, Message: msg}
}

func (e *AppError) With(msg string) *AppError {
	return &AppError{HTTP: e.HTTP, Code: e.Code, Message: msg}
}

func (e *AppError) Fmt(format string, args ...any) *AppError {
	return &AppError{HTTP: e.HTTP, Code: e.Code, Message: fmt.Sprintf(format, args...)}
}

var (
	Validation   = New(400, "VALIDATION", "参数校验失败")
	Unauthorized = New(401, "UNAUTHORIZED", "未登录或登录已失效")
	Forbidden    = New(403, "FORBIDDEN", "权限不足")
	NotFound     = New(404, "NOT_FOUND", "资源不存在")
	Conflict     = New(409, "CONFLICT", "资源冲突")
	ExamNotOpen  = New(422, "EXAM_NOT_OPEN", "考试未开始或已结束")
	SessionGone  = New(422, "SESSION_EXPIRED", "考试时间已到")
	ForceSubmit  = New(422, "FORCE_SUBMITTED", "切屏次数超限，已强制交卷")
	Internal     = New(500, "INTERNAL", "内部错误")
)
