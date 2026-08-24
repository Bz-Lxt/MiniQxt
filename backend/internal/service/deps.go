package service

import (
	"github.com/miniqxt/backend/internal/config"
	"github.com/miniqxt/backend/internal/queue"
	"gorm.io/gorm"
)

type App struct {
	DB     *gorm.DB
	Cfg    config.Config
	Grader *queue.SubmitQueue
	Traces *queue.TraceQueue
	Met    *queue.Metrics
}
