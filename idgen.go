package wampproto

import (
	"sync"
	"time"

	"golang.org/x/exp/rand"

	"github.com/xconnio/wampproto-go/util"
)

const maxID uint64 = 1 << 53

func init() {
	source := rand.NewSource(uint64(time.Now().UnixNano())) // #nosec
	rand.New(source)
}

// GenerateID generates a random WAMP ID.
func GenerateID() uint64 {
	id, _ := util.AsUInt64(rand.Int63n(int64(maxID)))
	return id
}

type SessionScopeIDGenerator struct {
	id uint64
	sync.Mutex
}

func (s *SessionScopeIDGenerator) NextID() uint64 {
	s.Lock()
	defer s.Unlock()

	if s.id == maxID {
		s.id = 0
	}

	s.id++
	return s.id
}
