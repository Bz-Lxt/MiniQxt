package queue

import "sync/atomic"

type Metrics struct {
	Enqueued   atomic.Int64
	Graded     atomic.Int64
	Failed     atomic.Int64
	Recovered  atomic.Int64
	TraceFlush atomic.Int64
	TraceRows  atomic.Int64
	SQLInserts atomic.Int64
}

func (m *Metrics) Snapshot(depth int, workers int) map[string]any {
	return map[string]any{
		"queue_depth":        depth,
		"workers":            workers,
		"enqueued":           m.Enqueued.Load(),
		"graded_total":       m.Graded.Load(),
		"failed":             m.Failed.Load(),
		"recovered":          m.Recovered.Load(),
		"trace_flushes":      m.TraceFlush.Load(),
		"trace_rows":         m.TraceRows.Load(),
		"sql_inserts":        m.SQLInserts.Load(),
		"write_amplification": amp(m.SQLInserts.Load(), m.TraceRows.Load()),
	}
}

func amp(sqls, rows int64) float64 {
	if rows == 0 {
		return 0
	}
	return float64(sqls) / float64(rows)
}
