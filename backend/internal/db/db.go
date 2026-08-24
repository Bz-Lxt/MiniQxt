package db

import (
	"time"

	"github.com/miniqxt/backend/internal/logger"
	"github.com/miniqxt/backend/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func Open(dsn string) (*gorm.DB, error) {
	gdb, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Warn),
	})
	if err != nil {
		return nil, err
	}
	sqlDB, err := gdb.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(80)
	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	return gdb, nil
}

func AutoMigrate(gdb *gorm.DB) error {
	err := gdb.AutoMigrate(
		&model.Tenant{},
		&model.User{},
		&model.Department{},
		&model.Position{},
		&model.Course{},
		&model.Chapter{},
		&model.CourseAssignment{},
		&model.LearningProgress{},
		&model.Question{},
		&model.QuestionOption{},
		&model.Paper{},
		&model.PaperItem{},
		&model.Exam{},
		&model.ExamAssignment{},
		&model.ExamSession{},
		&model.Submission{},
		&model.SubmissionAnswer{},
		&model.AnswerTrace{},
		&model.AntiCheatEvent{},
		&model.AuditFlag{},
		&model.GraderMetric{},
		&model.CertProgram{},
		&model.Certificate{},
	)
	if err != nil {
		return err
	}
	logger.Info("schema migrated")
	return nil
}
