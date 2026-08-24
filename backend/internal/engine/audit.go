package engine

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/miniqxt/backend/internal/model"
)

type FlagHit struct {
	Rule     string
	Title    string
	Severity string
	Evidence map[string]any
}

// Analyze applies the four interpretable anti-cheat rules (R1–R4).
func Analyze(traces []model.AnswerTrace, cheats []model.AntiCheatEvent, answers []model.SubmissionAnswer, questions map[uint64]model.Question) []FlagHit {
	var hits []FlagHit
	if h := ruleRhythm(traces); h != nil {
		hits = append(hits, *h)
	}
	if h := ruleInstantHard(traces, answers, questions); h != nil {
		hits = append(hits, *h)
	}
	if h := ruleAnswerAvalanche(traces, answers, questions); h != nil {
		hits = append(hits, *h)
	}
	if h := ruleBlurBoost(traces, cheats, answers, questions); h != nil {
		hits = append(hits, *h)
	}
	return hits
}

func ruleRhythm(traces []model.AnswerTrace) *FlagHit {
	if len(traces) < 8 {
		return nil
	}
	sort.Slice(traces, func(i, j int) bool { return traces[i].Seq < traces[j].Seq })
	gaps := make([]float64, 0, len(traces)-1)
	for i := 1; i < len(traces); i++ {
		d := traces[i].OccurredAt.Time().Sub(traces[i-1].OccurredAt.Time()).Seconds()
		if d <= 0 {
			d = 0.01
		}
		gaps = append(gaps, d)
	}
	mean, sd := meanStd(gaps)
	if mean > 0 && sd < 0.35 && mean < 4 {
		return &FlagHit{
			Rule: "R1", Title: "机器人节律", Severity: "high",
			Evidence: map[string]any{"mean_gap_sec": round2(mean), "stdev_sec": round2(sd), "n": len(gaps)},
		}
	}
	return nil
}

func ruleInstantHard(traces []model.AnswerTrace, answers []model.SubmissionAnswer, questions map[uint64]model.Question) *FlagHit {
	first := map[uint64]time.Time{}
	last := map[uint64]time.Time{}
	for _, t := range traces {
		ot := t.OccurredAt.Time()
		if _, ok := first[t.QuestionID]; !ok {
			first[t.QuestionID] = ot
		}
		last[t.QuestionID] = ot
	}
	type rec struct {
		QID uint64
		Sec float64
		Diff int
	}
	var hits []rec
	for _, a := range answers {
		q, ok := questions[a.QuestionID]
		if !ok || q.Type == model.QEssay || q.Difficulty < 3 {
			continue
		}
		if a.AutoScore <= 0 {
			continue
		}
		f, ok1 := first[a.QuestionID]
		l, ok2 := last[a.QuestionID]
		if !ok1 || !ok2 {
			continue
		}
		sec := l.Sub(f).Seconds()
		if sec < 2.5 {
			hits = append(hits, rec{QID: a.QuestionID, Sec: sec, Diff: q.Difficulty})
		}
	}
	if len(hits) >= 2 {
		return &FlagHit{
			Rule: "R2", Title: "秒答高分", Severity: "high",
			Evidence: map[string]any{"count": len(hits), "samples": hits},
		}
	}
	return nil
}

func ruleAnswerAvalanche(traces []model.AnswerTrace, answers []model.SubmissionAnswer, questions map[uint64]model.Question) *FlagHit {
	correctSet := map[uint64]map[uint64]struct{}{}
	for id, q := range questions {
		s := map[uint64]struct{}{}
		for _, o := range q.Options {
			if o.IsCorrect {
				s[o.ID] = struct{}{}
			}
		}
		correctSet[id] = s
	}
	sort.Slice(traces, func(i, j int) bool { return traces[i].Seq < traces[j].Seq })
	window := 45 * time.Second
	best := 0
	var bestAt time.Time
	for i := 0; i < len(traces); i++ {
		end := traces[i].OccurredAt.Time().Add(window)
		n := 0
		for j := i; j < len(traces); j++ {
			if traces[j].OccurredAt.Time().After(end) {
				break
			}
			if traces[j].FromOptionID == nil || traces[j].ToOptionID == nil {
				continue
			}
			set := correctSet[traces[j].QuestionID]
			_, fromOK := set[*traces[j].FromOptionID]
			_, toOK := set[*traces[j].ToOptionID]
			if !fromOK && toOK {
				n++
			}
		}
		if n > best {
			best = n
			bestAt = traces[i].OccurredAt.Time()
		}
	}
	if best >= 4 {
		return &FlagHit{
			Rule: "R3", Title: "答案雪崩", Severity: "high",
			Evidence: map[string]any{"wrong_to_right": best, "window_sec": 45, "from": bestAt.Format(time.RFC3339)},
		}
	}
	return nil
}

func ruleBlurBoost(traces []model.AnswerTrace, cheats []model.AntiCheatEvent, answers []model.SubmissionAnswer, questions map[uint64]model.Question) *FlagHit {
	if len(cheats) == 0 {
		return nil
	}
	correctAt := map[uint64]bool{}
	for _, a := range answers {
		q := questions[a.QuestionID]
		if q.Type == model.QEssay {
			continue
		}
		correctAt[a.QuestionID] = a.AutoScore > 0
	}
	boost := 0
	for _, ev := range cheats {
		if ev.EventType != "blur" && ev.EventType != "hidden" {
			continue
		}
		t0 := ev.OccurredAt.Time()
		t1 := t0.Add(20 * time.Second)
		for _, tr := range traces {
			ot := tr.OccurredAt.Time()
			if ot.Before(t0) || ot.After(t1) || tr.ToOptionID == nil {
				continue
			}
			if correctAt[tr.QuestionID] {
				boost++
			}
		}
	}
	if boost >= 3 {
		return &FlagHit{
			Rule: "R4", Title: "切屏关联正确率跃升", Severity: "medium",
			Evidence: map[string]any{"post_blur_correct_changes": boost, "blur_events": len(cheats)},
		}
	}
	return nil
}

func meanStd(xs []float64) (float64, float64) {
	if len(xs) == 0 {
		return 0, 0
	}
	var s float64
	for _, x := range xs {
		s += x
	}
	m := s / float64(len(xs))
	var v float64
	for _, x := range xs {
		d := x - m
		v += d * d
	}
	return m, math.Sqrt(v / float64(len(xs)))
}

func EvidenceJSON(m map[string]any) string {
	b, err := json.Marshal(m)
	if err != nil {
		return fmt.Sprintf("%v", m)
	}
	return string(b)
}
