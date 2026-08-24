package render

import (
	"github.com/miniqxt/backend/internal/model"
	"github.com/miniqxt/backend/internal/pkg/shuffle"
)

var letters = []string{"A", "B", "C", "D", "E", "F", "G", "H"}

// DisplayLetter maps a stable option_id to the letter this candidate saw.
func DisplayLetter(seed int64, questionID uint64, options []model.QuestionOption, optionID uint64, doShuffle bool) string {
	opts := append([]model.QuestionOption(nil), options...)
	if doShuffle && len(opts) > 1 {
		ord := shuffle.Order(shuffle.OptionSeed(seed, questionID), len(opts))
		opts = shuffle.Apply(opts, ord)
	}
	for i, o := range opts {
		if o.ID == optionID && i < len(letters) {
			return letters[i]
		}
	}
	return "?"
}

// QuestionPosition is 1-based display order for a candidate paper.
func QuestionPosition(seed int64, items []model.PaperItem, questionID uint64, doShuffle bool) int {
	cp := append([]model.PaperItem(nil), items...)
	if doShuffle {
		ord := shuffle.Order(seed, len(cp))
		cp = shuffle.Apply(cp, ord)
	}
	for i, it := range cp {
		if it.QuestionID == questionID {
			return i + 1
		}
	}
	return 0
}
