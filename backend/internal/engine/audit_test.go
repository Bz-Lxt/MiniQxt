package engine

import (
	"testing"
	"time"

	"github.com/miniqxt/backend/internal/model"
	"github.com/miniqxt/backend/internal/timeutil"
)

func ts(sec int) model.Time {
	base := time.Date(2026, 8, 24, 10, 0, 0, 0, timeutil.Beijing)
	return model.Time(base.Add(time.Duration(sec) * time.Second))
}

func TestR1Rhythm(t *testing.T) {
	var traces []model.AnswerTrace
	for i := 0; i < 12; i++ {
		traces = append(traces, model.AnswerTrace{Seq: i + 1, OccurredAt: ts(i * 2), QuestionID: uint64(i + 1)})
	}
	hits := Analyze(traces, nil, nil, nil)
	if !hasRule(hits, "R1") {
		t.Fatalf("expected R1, got %#v", hits)
	}
}

func TestR2InstantHard(t *testing.T) {
	qmap := map[uint64]model.Question{
		1: {ID: 1, Type: model.QSingle, Difficulty: 3},
		2: {ID: 2, Type: model.QSingle, Difficulty: 3},
	}
	answers := []model.SubmissionAnswer{{QuestionID: 1, AutoScore: 5}, {QuestionID: 2, AutoScore: 5}}
	traces := []model.AnswerTrace{
		{QuestionID: 1, Seq: 1, OccurredAt: ts(0)},
		{QuestionID: 1, Seq: 2, OccurredAt: ts(1)},
		{QuestionID: 2, Seq: 3, OccurredAt: ts(2)},
		{QuestionID: 2, Seq: 4, OccurredAt: ts(3)},
	}
	if !hasRule(Analyze(traces, nil, answers, qmap), "R2") {
		t.Fatal("expected R2")
	}
}

func TestR3Avalanche(t *testing.T) {
	from, to := uint64(1), uint64(2)
	qmap := map[uint64]model.Question{}
	var traces []model.AnswerTrace
	for i := 1; i <= 5; i++ {
		qid := uint64(i)
		qmap[qid] = model.Question{ID: qid, Type: model.QSingle, Options: []model.QuestionOption{{ID: 2, IsCorrect: true}, {ID: 1}}}
		f, tt := from, to
		traces = append(traces, model.AnswerTrace{QuestionID: qid, FromOptionID: &f, ToOptionID: &tt, Seq: i, OccurredAt: ts(i)})
	}
	if !hasRule(Analyze(traces, nil, nil, qmap), "R3") {
		t.Fatal("expected R3")
	}
}

func TestR4BlurBoost(t *testing.T) {
	to := uint64(2)
	qmap := map[uint64]model.Question{1: {ID: 1, Type: model.QSingle}, 2: {ID: 2, Type: model.QSingle}, 3: {ID: 3, Type: model.QSingle}}
	answers := []model.SubmissionAnswer{{QuestionID: 1, AutoScore: 5}, {QuestionID: 2, AutoScore: 5}, {QuestionID: 3, AutoScore: 5}}
	cheats := []model.AntiCheatEvent{{EventType: "blur", OccurredAt: ts(10)}}
	traces := []model.AnswerTrace{
		{QuestionID: 1, ToOptionID: &to, OccurredAt: ts(11), Seq: 1},
		{QuestionID: 2, ToOptionID: &to, OccurredAt: ts(12), Seq: 2},
		{QuestionID: 3, ToOptionID: &to, OccurredAt: ts(13), Seq: 3},
	}
	if !hasRule(Analyze(traces, cheats, answers, qmap), "R4") {
		t.Fatal("expected R4")
	}
}

func hasRule(hits []FlagHit, code string) bool {
	for _, h := range hits {
		if h.Rule == code {
			return true
		}
	}
	return false
}
