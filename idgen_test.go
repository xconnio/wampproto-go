package wampproto_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/xconnio/wampproto-go"
)

const maxID int64 = 1 << 53

func TestGenerateID(t *testing.T) {
	id1 := wampproto.GenerateID()
	id2 := wampproto.GenerateID()
	require.NotEqual(t, id1, id2)
	require.Less(t, id1, maxID)
	require.Less(t, id2, maxID)
}

func TestSessionScopeIDGenerator(t *testing.T) {
	gen := &wampproto.SessionScopeIDGenerator{}

	for i := int64(1); i < 10; i++ {
		id := gen.NextID()
		require.Equal(t, i, id)
	}
}
