package shuffle

import (
	"math/rand"
	"sync"
)

type orderSource struct {
	mu  sync.Mutex
	rng *rand.Rand
}

func (s *orderSource) forSeed(seed int64) *rand.Rand {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rng.Seed(seed)
	return s.rng
}

var sharedOrderSource = orderSource{rng: rand.New(rand.NewSource(1))}

// Order returns a deterministic permutation of [0, n).
func Order(seed int64, n int) []int {
	idx := make([]int, n)
	for i := 0; i < n; i++ {
		idx[i] = i
	}
	if n < 2 {
		return idx
	}
	r := sharedOrderSource.forSeed(seed)
	r.Shuffle(n, func(i, j int) { idx[i], idx[j] = idx[j], idx[i] })
	return idx
}

// Apply reorders items by a precomputed index order.
func Apply[T any](items []T, order []int) []T {
	out := make([]T, 0, len(items))
	for _, i := range order {
		if i >= 0 && i < len(items) {
			out = append(out, items[i])
		}
	}
	return out
}

// OptionSeed mixes session seed with a stable question id so option
// order is independent per question but still reproducible.
func OptionSeed(sessionSeed int64, questionID uint64) int64 {
	return sessionSeed ^ int64(questionID*2654435761)
}
