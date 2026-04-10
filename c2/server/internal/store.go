package internal

import (
	"sync"
	"time"

	"github.com/ianwong123/go-attack/c2/protocol"
)

// hold everything the server knows about connected agent
type AgentRecord struct {
	LastSeen time.Time
	Beacon   protocol.Beacon
}

// in-memory store for agents and their task
type Store struct {
	mu     sync.RWMutex
	agents map[string]*AgentRecord
	tasks  map[string][]protocol.Task
}

func New() *Store {
	return &Store{
		agents: make(map[string]*AgentRecord),
		tasks:  make(map[string][]protocol.Task),
	}
}

func (s *Store) Upsert(b protocol.Beacon) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.agents[b.AgentID] = &AgentRecord{
		LastSeen: time.Now(),
		Beacon:   b,
	}
}

// dequeue / pop the next task for the agent
func (s *Store) Dequeue(agentID string) *protocol.Task {
	s.mu.Lock()
	defer s.mu.Unlock()
	queue := s.tasks[agentID]
	if len(queue) == 0 {
		return nil
	}
	task := queue[0]
	s.tasks[agentID] = queue[1:]
	return &task
}

// Enqueue / add a task to agent's queue.
func (s *Store) Enqueue(agentID string, task protocol.Task) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks[agentID] = append(s.tasks[agentID], task)
}

func (s *Store) All() map[string]*AgentRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()
	copy := make(map[string]*AgentRecord, len(s.agents))
	for k, v := range s.agents {
		copy[k] = v
	}
	return copy
}
