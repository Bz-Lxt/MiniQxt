package service_test

import (
	"sync"
	"testing"

	"github.com/miniqxt/backend/internal/model"
	"github.com/miniqxt/backend/internal/service"
)

func TestRenderPaperConcurrentSessionsStayDeterministic(t *testing.T) {
	app := &service.App{}
	paper := concurrentRenderPaper(24, 6)
	seeds := []int64{1103, 2909}
	want := make(map[int64][]uint64, len(seeds))
	for _, seed := range seeds {
		want[seed] = renderedOrder(app.RenderPaper(paper, seed, false))
	}

	const workers = 16
	const iterations = 12
	start := make(chan struct{})
	failures := make(chan renderFailure, workers)
	var wg sync.WaitGroup
	for worker := 0; worker < workers; worker++ {
		worker := worker
		seed := seeds[worker%len(seeds)]
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			for iteration := 0; iteration < iterations; iteration++ {
				got := renderedOrder(app.RenderPaper(paper, seed, false))
				if !sameOrder(got, want[seed]) {
					failures <- renderFailure{worker: worker, iteration: iteration, seed: seed}
					return
				}
			}
		}()
	}

	close(start)
	wg.Wait()
	select {
	case failure := <-failures:
		t.Fatalf("worker %d iteration %d rendered seed %d with a different question or option order", failure.worker, failure.iteration, failure.seed)
	default:
	}
}

type renderFailure struct {
	worker    int
	iteration int
	seed      int64
}

func concurrentRenderPaper(questionN, optionN int) model.Paper {
	items := make([]model.PaperItem, 0, questionN)
	for question := 1; question <= questionN; question++ {
		options := make([]model.QuestionOption, 0, optionN)
		for option := 1; option <= optionN; option++ {
			options = append(options, model.QuestionOption{
				ID:         uint64(question*100 + option),
				QuestionID: uint64(question),
				SortNo:     option - 1,
			})
		}
		items = append(items, model.PaperItem{
			QuestionID: uint64(question),
			Score:      1,
			SortNo:     question - 1,
			Question: model.Question{
				ID:      uint64(question),
				Type:    model.QSingle,
				Options: options,
			},
		})
	}
	return model.Paper{
		ShuffleQuestions: true,
		ShuffleOptions:   true,
		Items:            items,
	}
}

func renderedOrder(items []service.RenderedItem) []uint64 {
	order := make([]uint64, 0, len(items)*7)
	for _, item := range items {
		order = append(order, item.QuestionID)
		for _, option := range item.Options {
			order = append(order, option.ID)
		}
	}
	return order
}

func sameOrder(left, right []uint64) bool {
	if len(left) != len(right) {
		return false
	}
	for i := range left {
		if left[i] != right[i] {
			return false
		}
	}
	return true
}
