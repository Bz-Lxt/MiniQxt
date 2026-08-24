package engine

import (
	"time"

	"github.com/miniqxt/backend/internal/model"
)

// IntegrityFromTelemetry implements C-01: the backend is the referee.
// A submission with no client telemetry cannot claim a clean exam.
func IntegrityFromTelemetry(sess model.ExamSession, traces int64, cheats int64) string {
	if sess.HeartbeatN == 0 && traces == 0 && cheats == 0 {
		return model.IntegrityUnverified
	}
	if sess.Suspicious || sess.BlurCount >= 3 {
		return model.IntegritySuspicious
	}
	return model.IntegrityOK
}

// TimelineGap reports whether the session wall clock is missing telemetry
// coverage. Used as an additional integrity signal, not a hard reject.
func TimelineGap(started, ended time.Time, heartbeats int, traces []model.AnswerTrace) bool {
	if ended.Before(started) || ended.Equal(started) {
		return false
	}
	dur := ended.Sub(started)
	if dur > 3*time.Minute && heartbeats == 0 && len(traces) == 0 {
		return true
	}
	if len(traces) < 2 {
		return false
	}
	var maxGap time.Duration
	for i := 1; i < len(traces); i++ {
		g := traces[i].OccurredAt.Time().Sub(traces[i-1].OccurredAt.Time())
		if g > maxGap {
			maxGap = g
		}
	}
	return maxGap > 10*time.Minute
}
