package wampproto

import (
	"sync"
	"time"

	"golang.org/x/exp/rand"
)

const maxID int64 = 1 << 53

func init() {
	source := rand.NewSource(uint64(time.Now().UnixNano()))
	rand.New(source)
}

// GenerateID generates a random WAMP ID.
func GenerateID() int64 {
	return rand.Int63n(maxID)
}

type SessionScopeIDGenerator struct {
	id int64
	sync.Mutex
}

func (s *SessionScopeIDGenerator) NextID() int64 {
	s.Lock()
	defer s.Unlock()

	if s.id == maxID {
		s.id = 0
	}

	s.id++
	return s.id
}
