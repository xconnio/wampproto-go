package internal_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/xconnio/wampproto-go/internal"
)

func TestMapDelete(t *testing.T) {
	m := &internal.Map[string, int]{}
	m.Store("key1", 100)

	m.Delete("key1")

	_, ok := m.Load("key1")
	require.False(t, ok)
}

func TestMapLoad(t *testing.T) {
	m := &internal.Map[string, bool]{}
	m.Store("key1", true)

	value, ok := m.Load("key1")
	require.True(t, ok)
	require.True(t, value)

	_, ok = m.Load("key2")
	require.False(t, ok)
}

func TestMapLoadAndDelete(t *testing.T) {
	m := &internal.Map[string, int]{}
	m.Store("key1", 100)

	value, loaded := m.LoadAndDelete("key1")
	require.True(t, loaded)
	require.Equal(t, 100, value)

	_, loaded = m.Load("key1")
	require.False(t, loaded)

	_, loaded = m.LoadAndDelete("key2")
	require.False(t, loaded)
}

func TestMapLoadOrStore(t *testing.T) {
	m := &internal.Map[string, string]{}

	actual, loaded := m.LoadOrStore("key1", "foo")
	require.False(t, loaded)
	require.Equal(t, "foo", actual)

	actual, loaded = m.LoadOrStore("key1", "bar")
	require.True(t, loaded)
	require.Equal(t, "foo", actual)
}

func TestMapRange(t *testing.T) {
	m := &internal.Map[string, int]{}
	m.Store("key1", 100)
	m.Store("key2", 200)

	keys := make(map[string]bool)
	values := make(map[int]bool)

	m.Range(func(key string, value int) bool {
		keys[key] = true
		values[value] = true
		return true
	})

	require.True(t, keys["key1"])
	require.True(t, keys["key2"])
	require.True(t, values[100])
	require.True(t, values[200])
}

func TestMapStore(t *testing.T) {
	m := &internal.Map[string, int]{}
	m.Store("key1", 100)

	value, ok := m.Load("key1")
	require.True(t, ok)
	require.Equal(t, 100, value)
}
