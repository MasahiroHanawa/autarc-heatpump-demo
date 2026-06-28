// Command server is the entry point for the Heat Pump Planning Assistant API.
// It is the composition root: loads configuration, wires dependencies, and
// manages the HTTP server and background worker lifecycle.
package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/MasahiroHanawa/autarc-heatpump-demo/apps/api/internal/analyzer"
	"github.com/MasahiroHanawa/autarc-heatpump-demo/apps/api/internal/jobs"
	"github.com/MasahiroHanawa/autarc-heatpump-demo/apps/api/internal/platform"
	"github.com/MasahiroHanawa/autarc-heatpump-demo/apps/api/internal/server"
	"github.com/MasahiroHanawa/autarc-heatpump-demo/apps/api/internal/worker"
)

func main() {
	cfg := platform.Load()
	logger := platform.NewLogger(cfg.LogLevel)

	// Dependency wiring — each layer only knows about the interface below it.
	store := jobs.NewInMemoryStore()
	svc := jobs.NewService(store)
	ana := analyzer.NewMock()
	wrk := worker.New(store, ana, logger)
	jobsHandler := server.NewJobsHandler(svc)
	router := server.NewRouter(logger, jobsHandler)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	// Start the background worker in its own goroutine.
	workerCtx, stopWorker := context.WithCancel(context.Background())
	go wrk.Run(workerCtx)

	// Start the HTTP server in its own goroutine.
	go func() {
		logger.Info("server starting", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	// Block until we receive an interrupt (Ctrl-C) or terminate signal.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	// Graceful shutdown: stop worker first, then drain in-flight HTTP requests.
	logger.Info("server shutting down")
	stopWorker()

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("graceful shutdown failed", "error", err)
		os.Exit(1)
	}
	logger.Info("server stopped")
}
