package auth_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/xconnio/wampproto-go/auth"
)

func TestISO8601(t *testing.T) {
	// UTC time
	t1 := time.Date(2024, 7, 10, 12, 30, 0, 0, time.UTC)
	expected1 := "2024-07-10T12:30:00Z"
	require.Equal(t, expected1, auth.ISO8601(t1))

	// Time with negative timezone offset
	t2 := time.Date(2024, 7, 10, 12, 30, 0, 0, time.FixedZone("CET", -1*60*60))
	expected2 := "2024-07-10T12:30:00-0100"
	require.Equal(t, expected2, auth.ISO8601(t2))

	// Time with positive timezone offset
	t3 := time.Date(2024, 7, 10, 12, 30, 0, 0, time.FixedZone("JST", 9*60*60))
	expected3 := "2024-07-10T12:30:00+0900"
	require.Equal(t, expected3, auth.ISO8601(t3))
}
