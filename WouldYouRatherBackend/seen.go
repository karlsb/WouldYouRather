package main

import (
	"slices"
	"sync"
)

// SeenStore records which pair IDs each user (by cookie UUID) has been served.
// It is safe for concurrent use by HTTP handlers.
type SeenStore struct {
	mu     sync.Mutex
	byUser map[string][]int
}

func NewSeenStore() *SeenStore {
	return &SeenStore{byUser: make(map[string][]int)}
}

// Snapshot returns a copy of the pair IDs user has been served.
func (s *SeenStore) Snapshot(user string) []int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return slices.Clone(s.byUser[user])
}

// Add records ids as served to user, registering the user if they are new.
func (s *SeenStore) Add(user string, ids ...int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.byUser[user] = append(s.byUser[user], ids...)
}

// ResetIfComplete clears user's history once they have been served total pairs
// and reports whether it did. Other users are not affected, and user stays Known.
func (s *SeenStore) ResetIfComplete(user string, total int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.byUser[user]) < total {
		return false
	}
	s.byUser[user] = s.byUser[user][:0]
	return true
}

// Known reports whether user has fetched pairs since the server started.
func (s *SeenStore) Known(user string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.byUser[user]
	return ok
}
