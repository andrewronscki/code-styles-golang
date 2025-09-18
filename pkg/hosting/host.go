package hosting

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type Host struct {
	Router  *gin.Engine
	Workers []Worker
	Addr    string
	server  *http.Server
}

func (h *Host) Start() {
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

	if h.Router != nil {
		if h.Addr == "" {
			h.Addr = ":8080"
		}

		h.server = &http.Server{
			Addr:    h.Addr,
			Handler: h.Router,
		}

		go func() {
			if err := h.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Printf("http server stopped")
				cancel()
			}
		}()
	}

	<-ctx.Done()
	log.Println("host shutdown started")
	defer log.Println("host shutdown completed")
	wg.Wait()

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	if err := h.server.Shutdown(ctx); err != nil {
		log.Printf("HTTP server was forced to shutdown")
	}
}
