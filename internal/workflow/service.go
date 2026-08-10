package workflow

import (
	"sync"
	"time"
)

type Workflow struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Version   string    `json:"version"`
	SpecYAML  string    `json:"spec_yaml"`
	CreatedAt time.Time `json:"created_at"`
}

type Store struct {
	mu        sync.RWMutex
	workflows map[string]*Workflow
}

func NewStore() *Store {
	return &Store{workflows: make(map[string]*Workflow)}
}

func (s *Store) Create(wf *Workflow) {
	s.mu.Lock()
	defer s.mu.Unlock()
	wf.CreatedAt = time.Now().UTC()
	s.workflows[wf.ID] = wf
}

func (s *Store) Get(id string) (*Workflow, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	wf, ok := s.workflows[id]
	return wf, ok
}

