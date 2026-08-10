package run

import (
	"sync"
	"time"
)

type Status string

const (
	StatusPending  Status = "PENDING"
	StatusRunning  Status = "RUNNING"
	StatusSuccess  Status = "SUCCEEDED"
	StatusFailed   Status = "FAILED"
	StatusCanceled Status = "CANCELED"
)

type Run struct {
	ID          string            `json:"run_id"`
	WorkflowID  string            `json:"workflow_id"`
	Version     string            `json:"version"`
	EnvID       string            `json:"env_id"`
	Status      Status            `json:"status"`
	Params      map[string]any    `json:"params"`
	CreatedAt   time.Time         `json:"created_at"`
	StartedAt   *time.Time        `json:"started_at,omitempty"`
	FinishedAt  *time.Time        `json:"finished_at,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

type Store struct {
	mu   sync.RWMutex
	runs map[string]*Run
}

func NewStore() *Store {
	return &Store{runs: make(map[string]*Run)}
}

func (s *Store) Create(r *Run) {
	s.mu.Lock()
	defer s.mu.Unlock()
	r.CreatedAt = time.Now().UTC()
	s.runs[r.ID] = r
}

func (s *Store) Get(id string) (*Run, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.runs[id]
	return r, ok
}
