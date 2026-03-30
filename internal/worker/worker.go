package worker

import (
	"context"
	"sync"

	"dcaiovinicius/httpchecker/internal/config"
	"dcaiovinicius/httpchecker/internal/logger"
	"dcaiovinicius/httpchecker/internal/request"
)

type Worker struct {
	urls   []string
	config *config.Config
}

func NewWorker(urls []string, cfg *config.Config) *Worker {
	return &Worker{
		urls:   urls,
		config: cfg,
	}
}

func (w *Worker) Run() {
	ctx, cancel := context.WithTimeout(context.Background(), w.config.GlobalTimeout)
	defer cancel()

	checker := request.NewChecker(w.config)

	buf := len(w.urls)
	if buf > w.config.Concurrency {
		buf = w.config.Concurrency
	}
	jobs := make(chan string, buf)

	var wg sync.WaitGroup
	for i := 0; i < w.config.Concurrency; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			logger.Info("worker %d started", id)
			for u := range jobs {
				if err := checker.Check(ctx, u); err != nil {
					logger.Info("check %s: %v", u, err)
				}
			}
		}(i)
	}

	for _, u := range w.urls {
		jobs <- u
	}
	close(jobs)
	wg.Wait()
}
