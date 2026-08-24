package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/miniqxt/backend/internal/config"
	"github.com/miniqxt/backend/internal/db"
	"github.com/miniqxt/backend/internal/logger"
	"github.com/miniqxt/backend/internal/queue"
	"github.com/miniqxt/backend/internal/router"
	"github.com/miniqxt/backend/internal/seed"
	"github.com/miniqxt/backend/internal/service"
	"github.com/miniqxt/backend/internal/timeutil"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()
	logger.Init(cfg.Env)
	logger.Info("boot", "tz", timeutil.Now().Format(time.RFC3339), "addr", cfg.HTTPAddr)

	conn, err := openRetry(cfg.MySQLDSN)
	if err != nil {
		logger.Error("mysql", "err", err)
		os.Exit(1)
	}
	if err := db.AutoMigrate(conn); err != nil {
		logger.Error("migrate", "err", err)
		os.Exit(1)
	}
	if err := seed.Run(conn); err != nil {
		logger.Error("seed", "err", err)
		os.Exit(1)
	}

	met := &queue.Metrics{}
	grader := queue.NewSubmitQueue(conn, cfg.GraderWorkers, cfg.GraderQueue, met)
	traces := queue.NewTraceQueue(conn, cfg.TraceBatch, cfg.TraceFlushMS, met)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	grader.Start(ctx)
	traces.Start(ctx)
	grader.Recover()

	app := &service.App{DB: conn, Cfg: cfg, Grader: grader, Traces: traces, Met: met}
	engine := router.New(app)
	srv := &http.Server{Addr: cfg.HTTPAddr, Handler: engine, ReadHeaderTimeout: 10 * time.Second}
	go func() {
		logger.Info("listen", "addr", cfg.HTTPAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("http", "err", err)
			os.Exit(1)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	cancel()
	shut, done := context.WithTimeout(context.Background(), 8*time.Second)
	defer done()
	_ = srv.Shutdown(shut)
	logger.Info("shutdown")
}

func openRetry(dsn string) (*gorm.DB, error) {
	var last error
	for i := 0; i < 40; i++ {
		conn, err := db.Open(dsn)
		if err == nil {
			return conn, nil
		}
		last = err
		logger.Warn("mysql retry", "n", i, "err", err)
		time.Sleep(2 * time.Second)
	}
	return nil, last
}
