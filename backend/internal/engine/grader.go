package engine

import (
	"sort"
	"strconv"
	"strings"

	"github.com/miniqxt/backend/internal/model"
)

// GradeObjective scores a single question. Essays return (0, true).
func GradeObjective(q model.Question, optionIDs []uint64, itemScore float64) (score float64, isEssay bool) {
	if q.Type == model.QEssay {
		return 0, true
	}
	correct := make([]uint64, 0, 4)
	for _, o := range q.Options {
		if o.IsCorrect {
			correct = append(correct, o.ID)
		}
	}
	sort.Slice(correct, func(i, j int) bool { return correct[i] < correct[j] })
	picked := uniqueSorted(optionIDs)
	switch q.Type {
	case model.QSingle, model.QJudge:
		if len(picked) == 1 && len(correct) == 1 && picked[0] == correct[0] {
			return itemScore, false
		}
		return 0, false
	case model.QMulti:
		if sameUint64(picked, correct) {
			return itemScore, false
		}
		if isSubset(picked, correct) && len(picked) > 0 {
			return round2(itemScore * float64(len(picked)) / float64(len(correct))), false
		}
		return 0, false
	default:
		return 0, false
	}
}

func ParseOptionIDs(s string) []uint64 {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]uint64, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		n, err := strconv.ParseUint(p, 10, 64)
		if err == nil {
			out = append(out, n)
		}
	}
	return out
}

func JoinOptionIDs(ids []uint64) string {
	ids = uniqueSorted(ids)
	parts := make([]string, 0, len(ids))
	for _, id := range ids {
		parts = append(parts, strconv.FormatUint(id, 10))
	}
	return strings.Join(parts, ",")
}

func uniqueSorted(ids []uint64) []uint64 {
	seen := map[uint64]struct{}{}
	out := make([]uint64, 0, len(ids))
	for _, id := range ids {
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

func sameUint64(a, b []uint64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func isSubset(a, b []uint64) bool {
	set := map[uint64]struct{}{}
	for _, x := range b {
		set[x] = struct{}{}
	}
	for _, x := range a {
		if _, ok := set[x]; !ok {
			return false
		}
	}
	return true
}

func round2(v float64) float64 {
	return float64(int(v*100+0.5)) / 100
}
