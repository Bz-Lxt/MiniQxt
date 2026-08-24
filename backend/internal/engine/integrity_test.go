package engine

import (
	"testing"
	"time"

	"github.com/miniqxt/backend/internal/model"
)

func TestIntegrityUnverified(t *testing.T) {
	got := IntegrityFromTelemetry(model.ExamSession{}, 0, 0)
	if got != model.IntegrityUnverified {
		t.Fatal(got)
	}
}

func TestIntegritySuspicious(t *testing.T) {
	got := IntegrityFromTelemetry(model.ExamSession{HeartbeatN: 3, BlurCount: 4}, 2, 1)
	if got != model.IntegritySuspicious {
		t.Fatal(got)
	}
}

func TestTimelineGap(t *testing.T) {
	st := time.Now()
	en := st.Add(5 * time.Minute)
	if !TimelineGap(st, en, 0, nil) {
		t.Fatal("expected gap")
	}
}
