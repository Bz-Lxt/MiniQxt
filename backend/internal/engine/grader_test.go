package engine

import (
	"testing"

	"github.com/miniqxt/backend/internal/model"
)

func q(typ string, correct ...uint64) model.Question {
	qq := model.Question{Type: typ, Score: 10}
	for i := uint64(1); i <= 4; i++ {
		ok := false
		for _, c := range correct {
			if c == i {
				ok = true
			}
		}
		qq.Options = append(qq.Options, model.QuestionOption{ID: i, IsCorrect: ok})
	}
	return qq
}

func TestGradeSingle(t *testing.T) {
	qq := q(model.QSingle, 2)
	got, essay := GradeObjective(qq, []uint64{2}, 5)
	if essay || got != 5 {
		t.Fatalf("%v %v", got, essay)
	}
	got, _ = GradeObjective(qq, []uint64{1}, 5)
	if got != 0 {
		t.Fatal(got)
	}
}

func TestGradeMultiPartial(t *testing.T) {
	qq := q(model.QMulti, 1, 2, 4)
	got, _ := GradeObjective(qq, []uint64{1, 2, 4}, 9)
	if got != 9 {
		t.Fatal(got)
	}
	got, _ = GradeObjective(qq, []uint64{1, 2}, 9)
	if got != 6 {
		t.Fatal(got)
	}
	got, _ = GradeObjective(qq, []uint64{1, 3}, 9)
	if got != 0 {
		t.Fatal(got)
	}
}

func TestGradeEssay(t *testing.T) {
	qq := model.Question{Type: model.QEssay}
	_, essay := GradeObjective(qq, nil, 10)
	if !essay {
		t.Fatal("essay")
	}
}

func TestJoinParse(t *testing.T) {
	s := JoinOptionIDs([]uint64{3, 1, 1, 2})
	if s != "1,2,3" {
		t.Fatal(s)
	}
	ids := ParseOptionIDs(s)
	if len(ids) != 3 || ids[0] != 1 {
		t.Fatal(ids)
	}
}

func TestIdempotentScore(t *testing.T) {
	qq := q(model.QJudge, 2)
	a, _ := GradeObjective(qq, []uint64{2}, 4)
	b, _ := GradeObjective(qq, []uint64{2}, 4)
	if a != b {
		t.Fatal(a, b)
	}
}
