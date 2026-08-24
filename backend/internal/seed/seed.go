package seed

import (
	"time"

	"github.com/miniqxt/backend/internal/logger"
	"github.com/miniqxt/backend/internal/model"
	"github.com/miniqxt/backend/internal/pkg/passwd"
	"github.com/miniqxt/backend/internal/timeutil"
	"gorm.io/gorm"
)

func Run(db *gorm.DB) error {
	var n int64
	db.Model(&model.User{}).Count(&n)
	if n > 0 {
		return nil
	}
	now := timeutil.Now()
	hash := func(p string) string {
		h, _ := passwd.Hash(p)
		return h
	}

	hq := model.Tenant{Name: "华清科技", Code: "hqtech", Status: model.StatusActive, BlurWarn: 3, BlurForce: 8, PassScore: 60, CreatedAt: model.Time(now)}
	xh := model.Tenant{Name: "星河制造", Code: "xinghe", Status: model.StatusActive, BlurWarn: 3, BlurForce: 8, PassScore: 70, CreatedAt: model.Time(now)}
	db.Create(&hq)
	db.Create(&xh)

	db.Create(&model.User{Username: "platform@qxt", PasswordHash: hash("Admin@123"), DisplayName: "平台管理员", Role: model.RolePlatform, Status: model.StatusActive, CreatedAt: model.Time(now)})

	dSafe := model.Department{TenantID: hq.ID, Name: "安全合规部", CreatedAt: model.Time(now)}
	dEng := model.Department{TenantID: hq.ID, Name: "工程研发部", CreatedAt: model.Time(now)}
	db.Create(&dSafe)
	db.Create(&dEng)
	p1 := model.Position{TenantID: hq.ID, Name: "安全员", Level: 1, CreatedAt: model.Time(now)}
	p2 := model.Position{TenantID: hq.ID, Name: "后端工程师", Level: 2, CreatedAt: model.Time(now)}
	db.Create(&p1)
	db.Create(&p2)

	db.Create(&model.User{TenantID: hq.ID, Username: "admin.hq@hqtech", PasswordHash: hash("Tenant@123"), DisplayName: "华清管理员", Role: model.RoleTenant, Status: model.StatusActive, CreatedAt: model.Time(now)})
	db.Create(&model.User{TenantID: hq.ID, Username: "teach.zhou@hqtech", PasswordHash: hash("Teach@123"), DisplayName: "周讲师", Role: model.RoleInstructor, DeptID: &dSafe.ID, Status: model.StatusActive, CreatedAt: model.Time(now)})
	li := model.User{TenantID: hq.ID, Username: "emp.li@hqtech", PasswordHash: hash("Emp@123"), DisplayName: "李明", Role: model.RoleEmployee, DeptID: &dSafe.ID, PositionID: &p1.ID, Status: model.StatusActive, CreatedAt: model.Time(now)}
	wang := model.User{TenantID: hq.ID, Username: "emp.wang@hqtech", PasswordHash: hash("Emp@123"), DisplayName: "王芳", Role: model.RoleEmployee, DeptID: &dEng.ID, PositionID: &p2.ID, Status: model.StatusActive, CreatedAt: model.Time(now)}
	db.Create(&li)
	db.Create(&wang)

	db.Create(&model.Department{TenantID: xh.ID, Name: "生产一部", CreatedAt: model.Time(now)})
	db.Create(&model.User{TenantID: xh.ID, Username: "admin.xh@xinghe", PasswordHash: hash("Tenant@123"), DisplayName: "星河管理员", Role: model.RoleTenant, Status: model.StatusActive, CreatedAt: model.Time(now)})
	db.Create(&model.User{TenantID: xh.ID, Username: "emp.chen@xinghe", PasswordHash: hash("Emp@123"), DisplayName: "陈强", Role: model.RoleEmployee, Status: model.StatusActive, CreatedAt: model.Time(now)})

	course := model.Course{TenantID: hq.ID, Title: "新员工入职与工地安全", Summary: "覆盖工地防护、应急疏散与信息安全基线，完课后方可参加合规考试。", Required: true, CreatedAt: model.Time(now)}
	db.Create(&course)
	ch := model.Chapter{TenantID: hq.ID, CourseID: course.ID, Title: "安全红线三十讲", VideoFile: "onboarding.mp4", DurationSec: 30, SortNo: 1, CreatedAt: model.Time(now)}
	db.Create(&ch)
	db.Create(&model.CourseAssignment{TenantID: hq.ID, CourseID: course.ID, Scope: "all", CreatedAt: model.Time(now)})

	type opt struct {
		l string
		c bool
	}
	addQ := func(typ, stem, tags string, diff int, score float64, opts []opt) model.Question {
		q := model.Question{TenantID: hq.ID, Type: typ, Stem: stem, Tags: tags, Difficulty: diff, Score: score, CreatedAt: model.Time(now)}
		db.Create(&q)
		for i, o := range opts {
			db.Create(&model.QuestionOption{QuestionID: q.ID, Label: o.l, IsCorrect: o.c, SortNo: i})
		}
		db.Preload("Options").First(&q, q.ID)
		return q
	}

	qs := []model.Question{
		addQ(model.QSingle, "进入施工现场必须佩戴的防护用品是？", "安全", 1, 5, []opt{{"安全帽", true}, {"棒球帽", false}, {"围巾", false}, {"墨镜即可", false}}),
		addQ(model.QSingle, "发现电气设备冒烟，首先应？", "安全", 2, 5, []opt{{"切断电源并上报", true}, {"浇水降温", false}, {"用衣服扑打", false}, {"继续作业观察", false}}),
		addQ(model.QMulti, "下列哪些属于信息安全红线？", "合规", 2, 8, []opt{{"不得转发客户数据到私人网盘", true}, {"不得共享个人账号", true}, {"可以将密码写在显示器上", false}, {"离职须交还访问权限", true}}),
		addQ(model.QJudge, "高处作业可以不系安全带，只要地面有人看护。", "安全", 1, 4, []opt{{"正确", false}, {"错误", true}}),
		addQ(model.QJudge, "考试过程中切屏超过阈值将被系统强制交卷。", "合规", 1, 4, []opt{{"正确", true}, {"错误", false}}),
		addQ(model.QSingle, "Go 中用于限制并发写入数据库的典型手段是？", "技术", 3, 6, []opt{{"Channel Worker + 批量落库", true}, {"在循环里无限制 goroutine", false}, {"关闭连接池", false}, {"只用全局锁锁住整个进程启动", false}}),
		addQ(model.QSingle, "JWT Claims 中的 tenant_id 应该来自？", "技术", 3, 6, []opt{{"登录签发的令牌，禁止请求参数传入", true}, {"前端任意 POST 字段", false}, {"URL 路径第一段", false}, {"Cookie 明文", false}}),
		addQ(model.QEssay, "请简述当你发现同事使用脚本批量交卷时，应如何上报并保留证据。", "合规", 2, 10, nil),
	}

	paper := model.Paper{TenantID: hq.ID, Title: "入职安全合规认证卷 A", ShuffleQuestions: true, ShuffleOptions: true, CreatedAt: model.Time(now)}
	db.Create(&paper)
	var total float64
	for i, q := range qs {
		db.Create(&model.PaperItem{PaperID: paper.ID, QuestionID: q.ID, Score: q.Score, GroupName: groupOf(q.Type), SortNo: i})
		total += q.Score
	}
	paper.TotalScore = total
	db.Save(&paper)

	exam := model.Exam{
		TenantID: hq.ID, PaperID: paper.ID, Title: "2026 入职安全合规认证",
		StartAt: model.Time(now.Add(-24 * time.Hour)), EndAt: model.Time(now.Add(30 * 24 * time.Hour)),
		DurationSec: 1800, PassScore: 60, MaxAttempts: 2, Status: "published", CreatedAt: model.Time(now),
	}
	db.Create(&exam)
	db.Create(&model.ExamAssignment{TenantID: hq.ID, ExamID: exam.ID, Scope: "all", CreatedAt: model.Time(now)})

	db.Create(&model.CertProgram{
		TenantID: hq.ID, Name: "安全员一级认证", PositionID: p1.ID, Level: 1,
		RequireCourse: course.ID, RequireExam: exam.ID, MinScore: 60, ValidDays: 365, CreatedAt: model.Time(now),
	})

	// Cross-tenant decoy exam so isolation tests have a foreign ID.
	xhPaper := model.Paper{TenantID: xh.ID, Title: "星河内部卷", CreatedAt: model.Time(now)}
	db.Create(&xhPaper)
	db.Create(&model.Exam{
		TenantID: xh.ID, PaperID: xhPaper.ID, Title: "星河保密考试",
		StartAt: model.Time(now.Add(-time.Hour)), EndAt: model.Time(now.Add(24 * time.Hour)),
		DurationSec: 600, PassScore: 70, MaxAttempts: 1, Status: "published", CreatedAt: model.Time(now),
	})

	logger.Info("seed completed", "hq", hq.ID, "xh", xh.ID)
	return nil
}

func groupOf(t string) string {
	if t == model.QEssay {
		return "主观题"
	}
	return "客观题"
}
