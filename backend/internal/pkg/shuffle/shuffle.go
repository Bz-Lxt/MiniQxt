package shuffle

import "math/rand"

// Order returns a deterministic permutation of [0, n).
//
// A fresh *rand.Rand is created for every call so that concurrent
// invocations with the same seed produce byte-identical permutations
// without any shared mutable state or lock contention. Reusing a single
// *rand.Rand across goroutines (even under a mutex) lets one Shuffle
// call corrupt the stream of another, which made the same shuffle_seed
// yield different question/option order under concurrent paper renders.
func Order(seed int64, n int) []int {
	idx := make([]int, n)
	for i := 0; i < n; i++ {
		idx[i] = i
	}
	if n < 2 {
		return idx
	}
	r := rand.New(rand.NewSource(seed))
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
