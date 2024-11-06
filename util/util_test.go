package util_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/xconnio/wampproto-go/util"
)

func TestAsInt64(t *testing.T) {
	t.Run("ValidConversion", func(t *testing.T) {
		tests := []struct {
			input    interface{}
			expected int64
		}{
			{input: int64(123), expected: 123},
			{input: uint64(456), expected: 456},
			{input: uint8(7), expected: 7},
			{input: 890, expected: 890},
			{input: int8(-12), expected: -12},
			{input: int32(345), expected: 345},
			{input: uint(678), expected: 678},
			{input: uint16(901), expected: 901},
			{input: uint32(234), expected: 234},
			{input: 56.78, expected: 56},
			{input: float32(9.01), expected: 9},
		}

		for _, test := range tests {
			result, ok := util.AsInt64(test.input)
			require.True(t, ok)
			require.Equal(t, test.expected, result)
		}
	})

	t.Run("InvalidConversion", func(t *testing.T) {
		result, ok := util.AsInt64("invalid")
		require.False(t, ok)
		require.Equal(t, int64(0), result)
	})
}

func TestAsFloat64(t *testing.T) {
	t.Run("ValidConversion", func(t *testing.T) {
		tests := []struct {
			input    any
			expected float64
		}{
			{input: float64(123.45), expected: 123.45},
			{input: float32(67.89), expected: 67.88999938964844},
			{input: int64(123), expected: 123.0},
			{input: uint64(456), expected: 456.0},
			{input: 789, expected: 789.0},
			{input: int8(-12), expected: -12.0},
			{input: int32(345), expected: 345.0},
			{input: uint(678), expected: 678.0},
			{input: uint32(234), expected: 234.0},
			{input: uint8(7), expected: 7.0},
			{input: uint16(90), expected: 90.0},
		}

		for _, test := range tests {
			result, ok := util.AsFloat64(test.input)
			require.True(t, ok)
			require.Equal(t, test.expected, result)
		}
	})

	t.Run("InvalidConversion", func(t *testing.T) {
		result, ok := util.AsFloat64("invalid")
		require.False(t, ok)
		require.Equal(t, float64(0), result)
	})
}

func TestAsBool(t *testing.T) {
	t.Run("ValidConversion", func(t *testing.T) {
		result, ok := util.AsBool(true)
		require.True(t, ok)
		require.True(t, result)

		result, ok = util.AsBool(false)
		require.True(t, ok)
		require.False(t, result)
	})

	t.Run("InvalidConversion", func(t *testing.T) {
		result, ok := util.AsBool(123)
		require.False(t, ok)
		require.False(t, result)
	})
}

func TestToBool(t *testing.T) {
	require.True(t, util.ToBool(true))
	require.False(t, util.ToBool(false))
	require.False(t, util.ToBool(123))
}

func TestAsString(t *testing.T) {
	t.Run("ValidConversion", func(t *testing.T) {
		result, ok := util.AsString("hello")
		require.True(t, ok)
		require.Equal(t, "hello", result)
	})

	t.Run("InvalidConversion", func(t *testing.T) {
		result, ok := util.AsString(123)
		require.False(t, ok)
		require.Equal(t, "", result)
	})
}

func TestToString(t *testing.T) {
	require.Equal(t, "hello", util.ToString("hello"))
	require.Equal(t, "", util.ToString(123))
}

func TestAnysToStrings(t *testing.T) {
	t.Run("ValidConversion", func(t *testing.T) {
		input := []any{"foo", "bar", "helloo"}

		result, err := util.AnysToStrings(input)
		require.NoError(t, err)
		require.Equal(t, []string{"foo", "bar", "helloo"}, result)
	})

	t.Run("InvalidConversion", func(t *testing.T) {
		input := []any{"foo", 123, "bar"}

		_, err := util.AnysToStrings(input)
		require.Error(t, err)
		require.Contains(t, err.Error(), "element 123 is not a string")
	})
}
