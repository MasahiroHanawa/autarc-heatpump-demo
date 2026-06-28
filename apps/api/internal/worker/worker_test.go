package worker_test

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/MasahiroHanawa/autarc-heatpump-demo/apps/api/internal/analyzer"
	"github.com/MasahiroHanawa/autarc-heatpump-demo/apps/api/internal/jobs"
	"github.com/MasahiroHanawa/autarc-heatpump-demo/apps/api/internal/worker"
)

func TestWorker_ProcessesPendingJob(t *testing.T) {
	store := jobs.NewInMemoryStore()
	svc := jobs.NewService(store)
	ana := analyzer.NewMock()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	wrk := worker.New(store, ana, logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go wrk.Run(ctx)

	// Create a job and wait for the worker to process it.
	job, _ := svc.CreateJob(ctx, jobs.Input{
		BuildingType:     jobs.BuildingTypeDetachedHouse,
		LivingAreaM2:     120,
		ConstructionYear: 1995,
		InsulationLevel:  jobs.InsulationLevelGood,
		CurrentHeatingSystem: jobs.HeatingSystemGasBoiler,
	})

	// Poll until completed or timeout.
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		updated, err := svc.GetJob(ctx, job.ID)
		if err != nil {
			t.Fatalf("GetJob failed: %v", err)
		}
		if updated.Status == jobs.StatusCompleted {
			if updated.Output == nil {
				t.Error("expected non-nil output on completed job")
			}
			t.Logf("job completed: suitability=%s heat_demand=%.0f kwh",
				updated.Output.Suitability, updated.Output.EstimatedHeatDemandKwh)
			return
		}
		time.Sleep(100 * time.Millisecond)
	}

	t.Fatal("job did not complete within 3 seconds")
}

func TestWorker_MarksJobFailedOnAnalyzerError(t *testing.T) {
	store := jobs.NewInMemoryStore()
	svc := jobs.NewService(store)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	wrk := worker.New(store, &errorAnalyzer{}, logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go wrk.Run(ctx)

	job, _ := svc.CreateJob(ctx, jobs.Input{
		BuildingType:         jobs.BuildingTypeApartment,
		LivingAreaM2:         60,
		ConstructionYear:     2000,
		InsulationLevel:      jobs.InsulationLevelAverage,
		CurrentHeatingSystem: jobs.HeatingSystemElectric,
	})

	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		updated, _ := svc.GetJob(ctx, job.ID)
		if updated.Status == jobs.StatusFailed {
			t.Logf("job correctly marked as failed")
			return
		}
		time.Sleep(100 * time.Millisecond)
	}

	t.Fatal("job was not marked failed within 3 seconds")
}

// errorAnalyzer always returns an error to simulate an analyzer failure.
type errorAnalyzer struct{}

func (e *errorAnalyzer) Analyze(_ context.Context, _ jobs.Input) (jobs.Output, error) {
	return jobs.Output{}, fmt.Errorf("simulated analyzer failure")
}
