// Command server is the entry point for the Heat Pump Planning Assistant API.
// It is the composition root: it loads configuration, builds shared
// dependencies (logger), wires up the HTTP server, and manages its lifecycle.
package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/MasahiroHanawa/autarc-heatpump-demo/apps/api/internal/platform"
)

func main() {
	// Configuration and the logger come first — everything else depends on
	// them. This is the start of the dependency wiring that grows in #14.
	cfg := platform.Load()
	logger := platform.NewLogger(cfg.LogLevel)

	// A minimal router for now. The real jobs API arrives in later issues
	// (#11/#12). The /healthz endpoint lets us — and Docker/Kubernetes — verify
	// the process is alive. The "GET /healthz" pattern uses method-based
	// routing from net/http (Go 1.22+).
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}

	// Run the server in its own goroutine so main can block waiting for an OS
	// shutdown signal. ListenAndServe blocks until the server is closed.
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

	// Graceful shutdown: stop accepting new connections and give in-flight
	// requests up to ShutdownTimeout to finish before forcing exit.
	logger.Info("server shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("graceful shutdown failed", "error", err)
		os.Exit(1)
	}
	logger.Info("server stopped")
}
