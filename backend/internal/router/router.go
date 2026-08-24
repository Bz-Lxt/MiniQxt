package router

import (
	"github.com/gin-gonic/gin"
	"github.com/miniqxt/backend/internal/handler"
	"github.com/miniqxt/backend/internal/middleware"
	"github.com/miniqxt/backend/internal/model"
	"github.com/miniqxt/backend/internal/service"
	"github.com/miniqxt/backend/internal/timeutil"
)

func New(app *service.App) *gin.Engine {
	if app.Cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(middleware.Recovery(), middleware.RequestLog(), cors())
	h := &handler.API{App: app}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok", "time": timeutil.Now().Format("2006-01-02T15:04:05+08:00"),
			"queue_depth": app.Grader.Depth(), "workers": app.Grader.Workers(),
		})
	})
	r.GET("/api/v1/media/:name", h.Media)

	v1 := r.Group("/api/v1")
	v1.POST("/auth/login", h.Login)

	auth := v1.Group("", middleware.Auth(app.Cfg.JWTSecret))
	auth.GET("/auth/me", h.Me)
	auth.GET("/grader/metrics", middleware.Staff(), h.Metrics)

	plat := auth.Group("", middleware.Roles(model.RolePlatform))
	plat.GET("/tenants", h.ListTenants)
	plat.POST("/tenants", h.CreateTenant)
	plat.PATCH("/tenants/:id", h.PatchTenant)

	staff := auth.Group("", middleware.Staff())
	staff.GET("/departments", h.ListDepts)
	staff.POST("/departments", h.CreateDept)
	staff.PATCH("/departments/:id", h.PatchDept)
	staff.DELETE("/departments/:id", h.DeleteDept)
	staff.GET("/positions", h.ListPositions)
	staff.POST("/positions", h.CreatePosition)
	staff.PATCH("/positions/:id", h.PatchPosition)
	staff.DELETE("/positions/:id", h.DeletePosition)
	staff.GET("/employees", h.ListEmployees)
	staff.POST("/employees", h.CreateEmployee)
	staff.PATCH("/employees/:id", h.PatchEmployee)
	staff.DELETE("/employees/:id", h.DeleteEmployee)

	staff.GET("/courses", h.ListCourses)
	staff.POST("/courses", h.CreateCourse)
	staff.GET("/courses/:id", h.GetCourse)
	staff.PATCH("/courses/:id", h.PatchCourse)
	staff.DELETE("/courses/:id", h.DeleteCourse)
	staff.POST("/courses/:id/chapters", h.AddChapter)
	staff.POST("/courses/:id/assign", h.AssignCourse)
	staff.PATCH("/chapters/:id", h.PatchChapter)
	staff.DELETE("/chapters/:id", h.DeleteChapter)

	staff.GET("/questions", h.ListQuestions)
	staff.POST("/questions", h.CreateQuestion)
	staff.GET("/questions/:id", h.GetQuestion)
	staff.PATCH("/questions/:id", h.UpdateQuestion)
	staff.DELETE("/questions/:id", h.DeleteQuestion)
	staff.GET("/papers", h.ListPapers)
	staff.POST("/papers", h.CreatePaper)
	staff.GET("/papers/:id", h.GetPaper)
	staff.PUT("/papers/:id/items", h.SavePaperItems)
	staff.GET("/papers/:id/preview", h.PreviewPaper)
	staff.DELETE("/papers/:id", h.DeletePaper)

	staff.GET("/exams", h.ListExams)
	staff.POST("/exams", h.CreateExam)
	staff.GET("/exams/:id", h.GetExam)
	staff.POST("/exams/:id/assign", h.AssignExam)
	staff.GET("/submissions", h.ListSubmissions)
	staff.PUT("/submissions/:id/manual-scores", h.ManualScore)
	staff.GET("/exam-sessions/:id/timeline", h.Timeline)
	staff.POST("/exam-sessions/:id/run-audit", h.RunAudit)
	staff.GET("/audit-flags", h.ListFlags)
	staff.GET("/cert-programs", h.ListPrograms)
	staff.POST("/cert-programs", h.CreateProgram)
	staff.PATCH("/cert-programs/:id", h.PatchProgram)
	staff.POST("/cert-programs/:id/evaluate", h.EvaluateProgram)
	staff.GET("/certificates", h.ListCertificates)

	auth.GET("/my/courses", h.MyCourses)
	auth.POST("/my/chapters/:id/progress", h.ReportProgress)
	auth.GET("/my/exams", h.MyExams)
	auth.POST("/exams/:id/start", h.StartExam)
	auth.GET("/exam-sessions/:id", h.GetSession)
	auth.GET("/exam-sessions/:id/paper", h.SessionPaper)
	auth.POST("/exam-sessions/:id/heartbeat", h.Heartbeat)
	auth.POST("/exam-sessions/:id/traces", h.Traces)
	auth.POST("/exam-sessions/:id/anti-cheat", h.AntiCheat)
	auth.POST("/exam-sessions/:id/submit", h.Submit)
	auth.GET("/submissions/:id", h.GetSubmission)
	auth.GET("/my/certificates", h.ListCertificates)
	return r
}

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
