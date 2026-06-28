// Package worker provides the background processing loop that picks up
// pending jobs and drives them through the analyzer.
package worker

import (
	"context"
	"log/slog"
	"time"

	"github.com/MasahiroHanawa/autarc-heatpump-demo/apps/api/internal/analyzer"
	"github.com/MasahiroHanawa/autarc-heatpump-demo/apps/api/internal/jobs"
)

// Worker polls the repository for pending jobs on a fixed interval and
// processes each one through the analyzer.
type Worker struct {
	repo     jobs.Repository
	analyzer analyzer.Analyzer
	logger   *slog.Logger
	interval time.Duration
}

func New(repo jobs.Repository, a analyzer.Analyzer, logger *slog.Logger) *Worker {
	return &Worker{
		repo:     repo,
		analyzer: a,
		logger:   logger,
		interval: 500 * time.Millisecond,
	}
}

// Run starts the polling loop and blocks until ctx is cancelled.
func (w *Worker) Run(ctx context.Context) {
	w.logger.Info("worker started")
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			w.logger.Info("worker stopped")
			return
		case <-ticker.C:
			w.tick(ctx)
		}
	}
}

func (w *Worker) tick(ctx context.Context) {
	all, err := w.repo.List(ctx)
	if err != nil {
		w.logger.Error("worker: list failed", "error", err)
		return
	}
	for _, job := range all {
		if job.Status == jobs.StatusPending {
			w.process(ctx, job)
		}
	}
}

func (w *Worker) process(ctx context.Context, job jobs.Job) {
	w.logger.Info("processing job", "id", job.ID)

	job.Status = jobs.StatusProcessing
	job.UpdatedAt = time.Now().UTC()
	if err := w.repo.Update(ctx, job); err != nil {
		w.logger.Error("worker: failed to mark processing", "id", job.ID, "error", err)
		return
	}

	output, err := w.analyzer.Analyze(ctx, job.Input)
	if err != nil {
		w.logger.Error("worker: analyzer failed", "id", job.ID, "error", err)
		job.Status = jobs.StatusFailed
		job.UpdatedAt = time.Now().UTC()
		_ = w.repo.Update(ctx, job)
		return
	}

	job.Status = jobs.StatusCompleted
	job.Output = &output
	job.UpdatedAt = time.Now().UTC()
	if err := w.repo.Update(ctx, job); err != nil {
		w.logger.Error("worker: failed to mark completed", "id", job.ID, "error", err)
		return
	}

	w.logger.Info("job completed", "id", job.ID, "suitability", output.Suitability)
}
