package render

import (
	"testing"

	"github.com/miniqxt/backend/internal/model"
)

func TestDisplayLetterStable(t *testing.T) {
	opts := []model.QuestionOption{{ID: 10}, {ID: 11}, {ID: 12}, {ID: 13}}
	a := DisplayLetter(7, 3, opts, 11, true)
	b := DisplayLetter(7, 3, opts, 11, true)
	if a != b || a == "?" {
		t.Fatalf("%s %s", a, b)
	}
	c := DisplayLetter(8, 3, opts, 11, true)
	if a == c {
		// possible but rare; just ensure function returns a letter
		if len(c) != 1 {
			t.Fatal(c)
		}
	}
}

func TestQuestionPosition(t *testing.T) {
	items := []model.PaperItem{{QuestionID: 1}, {QuestionID: 2}, {QuestionID: 3}}
	p1 := QuestionPosition(1, items, 2, true)
	p2 := QuestionPosition(1, items, 2, true)
	if p1 != p2 || p1 == 0 {
		t.Fatal(p1, p2)
	}
}
