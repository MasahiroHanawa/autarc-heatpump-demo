package jobs

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"
)

// Service contains the business logic for managing analysis jobs.
// It sits between the HTTP handlers and the repository, keeping both sides
// free of business rules.
type Service struct {
	repo Repository
}

// NewService wires the service to a Repository implementation.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// CreateJob persists a new job in the pending state and returns it.
func (s *Service) CreateJob(ctx context.Context, input Input) (Job, error) {
	now := time.Now().UTC()
	job := Job{
		ID:        newID(),
		Input:     input,
		Status:    StatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.repo.Create(ctx, job); err != nil {
		return Job{}, err
	}
	return job, nil
}

// GetJob retrieves a single job by ID. Returns ErrNotFound if it does not exist.
func (s *Service) GetJob(ctx context.Context, id string) (Job, error) {
	return s.repo.GetByID(ctx, id)
}

// ListJobs returns all jobs in the store, in no guaranteed order.
func (s *Service) ListJobs(ctx context.Context) ([]Job, error) {
	return s.repo.List(ctx)
}

// newID generates a random UUID v4 style identifier using crypto/rand.
func newID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	// Set version (4) and variant bits per RFC 4122.
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
