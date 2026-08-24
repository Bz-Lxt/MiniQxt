package queue

import (
	"context"
	"time"

	"github.com/miniqxt/backend/internal/logger"
	"github.com/miniqxt/backend/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TraceQueue struct {
	db      *gorm.DB
	ch      chan model.AnswerTrace
	batch   int
	flush   time.Duration
	Metrics *Metrics
}

func NewTraceQueue(db *gorm.DB, batch, flushMS int, m *Metrics) *TraceQueue {
	if batch < 50 {
		batch = 50
	}
	if flushMS < 50 {
		flushMS = 200
	}
	return &TraceQueue{
		db: db, ch: make(chan model.AnswerTrace, 8192),
		batch: batch, flush: time.Duration(flushMS) * time.Millisecond, Metrics: m,
	}
}

func (q *TraceQueue) Start(ctx context.Context) {
	go q.loop(ctx)
}

func (q *TraceQueue) Push(evs []model.AnswerTrace) {
	for _, e := range evs {
		select {
		case q.ch <- e:
		default:
			logger.Warn("trace channel full, dropping one event")
		}
	}
}

func (q *TraceQueue) loop(ctx context.Context) {
	buf := make([]model.AnswerTrace, 0, q.batch)
	tick := time.NewTicker(q.flush)
	defer tick.Stop()
	flush := func() {
		if len(buf) == 0 {
			return
		}
		rows := make([]model.AnswerTrace, len(buf))
		copy(rows, buf)
		buf = buf[:0]
		if err := q.db.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(rows, q.batch).Error; err != nil {
			logger.Error("trace batch insert", "err", err, "n", len(rows))
			return
		}
		q.Metrics.SQLInserts.Add(1)
		q.Metrics.TraceRows.Add(int64(len(rows)))
		q.Metrics.TraceFlush.Add(1)
	}
	for {
		select {
		case <-ctx.Done():
			flush()
			return
		case ev := <-q.ch:
			buf = append(buf, ev)
			if len(buf) >= q.batch {
				flush()
			}
		case <-tick.C:
			flush()
		}
	}
}
