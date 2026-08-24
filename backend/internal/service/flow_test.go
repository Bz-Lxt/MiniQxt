package service

import (
	"os"
	"testing"
	"time"

	"github.com/miniqxt/backend/internal/config"
	"github.com/miniqxt/backend/internal/db"
	"github.com/miniqxt/backend/internal/logger"
	"github.com/miniqxt/backend/internal/model"
	"github.com/miniqxt/backend/internal/queue"
	"github.com/miniqxt/backend/internal/timeutil"
	"gorm.io/gorm"
)

func liveApp(t *testing.T) (*App, func()) {
	t.Helper()
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		dsn = "miniqxt:miniqxt_pass@tcp(127.0.0.1:28391)/miniqxt?charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai"
	}
	logger.Init("test")
	conn, err := db.Open(dsn)
	if err != nil {
		t.Skip("mysql not available:", err)
	}
	met := &queue.Metrics{}
	g := queue.NewSubmitQueue(conn, 4, 256, met)
	tr := queue.NewTraceQueue(conn, 50, 50, met)
	g.Start(t.Context())
	tr.Start(t.Context())
	app := &App{DB: conn, Cfg: config.Load(), Grader: g, Traces: tr, Met: met}
	return app, func() {}
}

func TestCrossTenantHidden(t *testing.T) {
	app, _ := liveApp(t)
	var hq, xh model.User
	if err := app.DB.Where("username = ?", "emp.li@hqtech").First(&hq).Error; err != nil {
		t.Skip(err)
	}
	if err := app.DB.Where("username = ?", "emp.chen@xinghe").First(&xh).Error; err != nil {
		t.Fatal(err)
	}
	_, err := app.StartExam(xh.TenantID, xh.ID, 1)
	if err == nil {
		t.Fatal("xinghe must not start hqtech exam")
	}
}

func TestSubmitIdempotentLive(t *testing.T) {
	app, _ := liveApp(t)
	var u model.User
	if err := app.DB.Where("username = ?", "emp.wang@hqtech").First(&u).Error; err != nil {
		t.Skip(err)
	}
	// wang already submitted in smoke; create a disposable session row
	now := timeutil.Now()
	s := model.ExamSession{
		TenantID: u.TenantID, ExamID: 1, UserID: u.ID, Attempt: 99,
		ShuffleSeed: 1, StartedAt: model.Time(now), EndsAt: model.Time(now.Add(time.Hour)),
		Status: model.SessInProgress, HeartbeatN: 2, CreatedAt: model.Time(now),
	}
	if err := app.DB.Create(&s).Error; err != nil {
		t.Skip("session unique", err)
	}
	a1, err := app.SubmitExam(u.TenantID, u.ID, s.ID, nil)
	if err != nil {
		t.Fatal(err)
	}
	a2, err := app.SubmitExam(u.TenantID, u.ID, s.ID, nil)
	if err != nil {
		t.Fatal(err)
	}
	if a1["submission_id"] != a2["submission_id"] {
		t.Fatalf("%v %v", a1, a2)
	}
}

func TestRecoveryRequeues(t *testing.T) {
	app, _ := liveApp(t)
	now := timeutil.Now()
	row := model.Submission{
		TenantID: 1, SessionID: 999001, ExamID: 1, UserID: 4,
		Status: model.SubQueued, QueuedAt: model.Time(now), CreatedAt: model.Time(now),
	}
	if err := app.DB.Create(&row).Error; err != nil {
		t.Skip(err)
	}
	n := app.Grader.Recover()
	if n < 1 {
		t.Fatal("expected recover >= 1")
	}
	app.DB.Delete(&row)
}

func TestTenantScopeQuery(t *testing.T) {
	app, _ := liveApp(t)
	var n int64
	app.DB.Model(&model.Exam{}).Where("tenant_id = ?", 1).Count(&n)
	if n < 1 {
		t.Fatal("hqtech exams")
	}
	var foreign model.Exam
	err := app.DB.Where("tenant_id = ? AND id = ?", 2, 1).First(&foreign).Error
	if err == nil {
		t.Fatal("exam 1 must not belong to tenant 2")
	}
	if err != nil && err != gorm.ErrRecordNotFound && err.Error() == "" {
		t.Fatal(err)
	}
}
