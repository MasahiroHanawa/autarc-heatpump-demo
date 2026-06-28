package jobs

import (
	"context"
	"errors"
	"sync"
)

// ErrNotFound is returned when a job ID does not exist in the store.
var ErrNotFound = errors.New("job not found")

// Repository is the persistence contract for jobs. The in-memory
// implementation lives here; a PostgreSQL implementation can swap in later
// without touching any other package.
type Repository interface {
	Create(ctx context.Context, job Job) error
	GetByID(ctx context.Context, id string) (Job, error)
	List(ctx context.Context) ([]Job, error)
	Update(ctx context.Context, job Job) error
}

// InMemoryStore is a thread-safe in-memory implementation of Repository.
type InMemoryStore struct {
	mu   sync.RWMutex
	jobs map[string]Job
}

// NewInMemoryStore returns an initialised, empty store.
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{jobs: make(map[string]Job)}
}

func (s *InMemoryStore) Create(_ context.Context, job Job) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.jobs[job.ID] = job
	return nil
}

func (s *InMemoryStore) GetByID(_ context.Context, id string) (Job, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	job, ok := s.jobs[id]
	if !ok {
		return Job{}, ErrNotFound
	}
	return job, nil
}

func (s *InMemoryStore) List(_ context.Context) ([]Job, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	jobs := make([]Job, 0, len(s.jobs))
	for _, j := range s.jobs {
		jobs = append(jobs, j)
	}
	return jobs, nil
}

func (s *InMemoryStore) Update(_ context.Context, job Job) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.jobs[job.ID]; !ok {
		return ErrNotFound
	}
	s.jobs[job.ID] = job
	return nil
}
