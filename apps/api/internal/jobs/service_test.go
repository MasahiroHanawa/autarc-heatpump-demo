package jobs_test

import (
	"context"
	"testing"

	"github.com/MasahiroHanawa/autarc-heatpump-demo/apps/api/internal/jobs"
)

func TestCreateAndGetJob(t *testing.T) {
	ctx := context.Background()

	store := jobs.NewInMemoryStore()
	svc := jobs.NewService(store)

	input := jobs.Input{
		BuildingType:               jobs.BuildingTypeDetachedHouse,
		LivingAreaM2:               120.0,
		ConstructionYear:           1990,
		InsulationLevel:            jobs.InsulationLevelAverage,
		CurrentHeatingSystem:       jobs.HeatingSystemGasBoiler,
		AnnualEnergyConsumptionKwh: 18000,
		InstallerNotes:             "south-facing roof, no basement",
	}

	// Create a job and verify initial state
	job, err := svc.CreateJob(ctx, input)
	if err != nil {
		t.Fatalf("CreateJob failed: %v", err)
	}
	if job.ID == "" {
		t.Error("expected a non-empty ID")
	}
	if job.Status != jobs.StatusPending {
		t.Errorf("expected status pending, got %s", job.Status)
	}
	if job.Output != nil {
		t.Error("expected Output to be nil for a new job")
	}

	t.Logf("created job id=%s status=%s", job.ID, job.Status)

	// Retrieve it back and verify it matches
	fetched, err := svc.GetJob(ctx, job.ID)
	if err != nil {
		t.Fatalf("GetJob failed: %v", err)
	}
	if fetched.ID != job.ID {
		t.Errorf("expected id %s, got %s", job.ID, fetched.ID)
	}

	t.Logf("fetched job id=%s status=%s", fetched.ID, fetched.Status)
}

func TestGetJobNotFound(t *testing.T) {
	ctx := context.Background()

	store := jobs.NewInMemoryStore()
	svc := jobs.NewService(store)

	_, err := svc.GetJob(ctx, "non-existent-id")
	if err == nil {
		t.Fatal("expected ErrNotFound, got nil")
	}
	t.Logf("got expected error: %v", err)
}

func TestListJobs(t *testing.T) {
	ctx := context.Background()

	store := jobs.NewInMemoryStore()
	svc := jobs.NewService(store)

	input := jobs.Input{BuildingType: jobs.BuildingTypeApartment, LivingAreaM2: 60}

	svc.CreateJob(ctx, input)
	svc.CreateJob(ctx, input)

	list, err := svc.ListJobs(ctx)
	if err != nil {
		t.Fatalf("ListJobs failed: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 jobs, got %d", len(list))
	}

	t.Logf("listed %d jobs", len(list))
}
