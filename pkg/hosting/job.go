package hosting

import (
	"context"
	"log"
	"os/signal"
	"sync"
	"syscall"
)

type JobHost struct {
	Workers []Worker
}

func (h *JobHost) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	wg := &sync.WaitGroup{}

	for _, w := range h.Workers {
		wg.Add(1)
		worker := w
		go worker.Run(ctx, wg.Done)
	}

	doneCh := make(chan struct{})
	go func() {
		wg.Wait()
		close(doneCh)
	}()

	select {
	case <-ctx.Done():
		log.Println("host shutdown started")
		defer log.Println("host shutdown completed")
	case <-doneCh:
		log.Printf("all jobs have finished, host is shutting down")
	}

	wg.Wait()
	stop()
}
