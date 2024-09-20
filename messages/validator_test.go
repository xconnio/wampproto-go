package messages_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/xconnio/wampproto-go/messages"
)

func TestValidateArgs(t *testing.T) {
	t.Run("NoArgs", func(t *testing.T) {
		wampMsg := []any{1, "io.xconn.test", map[string]any{}}
		index := 3
		var fields messages.Fields

		err := messages.ValidateArgs(wampMsg, index, &fields)
		require.NoError(t, err)
		require.Nil(t, fields.Args)
	})

	t.Run("ValidType", func(t *testing.T) {
		wampMsg := []any{1, "io.xconn.test", map[string]any{}, []any{"abc", 123, true}}
		index := 3
		var fields messages.Fields

		err := messages.ValidateArgs(wampMsg, index, &fields)
		require.NoError(t, err)
		require.Equal(t, []any{"abc", 123, true}, fields.Args)
	})

	t.Run("InvalidType", func(t *testing.T) {
		wampMsg := []any{1, "io.xconn.test", map[string]any{}, "invalidType"}
		index := 3
		var fields messages.Fields

		err := messages.ValidateArgs(wampMsg, index, &fields)
		require.EqualError(t, err, "item at index 3 must be of type []any but was string")
	})
}

func TestValidateSessionID(t *testing.T) {
	t.Run("ValidID", func(t *testing.T) {
		wampMsg := []any{1, "io.xconn.test", map[string]any{}}
		index := 0
		var fields messages.Fields

		err := messages.ValidateSessionID(wampMsg, index, &fields)
		require.NoError(t, err)
		require.Equal(t, int64(1), fields.SessionID)
	})

	t.Run("InvalidType", func(t *testing.T) {
		wampMsg := []any{"invalidType", "io.xconn.test", map[string]any{}}
		index := 0
		var fields messages.Fields

		err := messages.ValidateSessionID(wampMsg, index, &fields)
		require.EqualError(t, err, "item at index 0 must be of type int64 but was string")
	})
}

func TestValidateMessageType(t *testing.T) {
	t.Run("ValidType", func(t *testing.T) {
		wampMsg := []any{1, 1, map[string]any{}}
		index := 1
		var fields messages.Fields

		err := messages.ValidateMessageType(wampMsg, index, &fields)
		require.NoError(t, err)
		require.Equal(t, int64(1), fields.MessageType)
	})

	t.Run("InvalidType", func(t *testing.T) {
		wampMsg := []any{1, "invalid", map[string]any{}}
		index := 1
		var fields messages.Fields

		err := messages.ValidateMessageType(wampMsg, index, &fields)
		require.EqualError(t, err, "item at index 1 must be of type int64 but was string")
	})
}

func TestValidateRequestID(t *testing.T) {
	t.Run("ValidType", func(t *testing.T) {
		wampMsg := []any{1, 1, map[string]any{}}
		index := 0
		var fields messages.Fields

		err := messages.ValidateRequestID(wampMsg, index, &fields)
		require.NoError(t, err)
		require.Equal(t, int64(1), fields.RequestID)
	})

	t.Run("InvalidType", func(t *testing.T) {
		wampMsg := []any{"invalid", 1, map[string]any{}}
		index := 0
		var fields messages.Fields

		err := messages.ValidateRequestID(wampMsg, index, &fields)
		require.EqualError(t, err, "item at index 0 must be of type int64 but was string")
	})
}

func TestValidateRegistrationID(t *testing.T) {
	t.Run("ValidType", func(t *testing.T) {
		wampMsg := []any{1, 1, map[string]any{}}
		index := 1
		var fields messages.Fields

		err := messages.ValidateRegistrationID(wampMsg, index, &fields)
		require.NoError(t, err)
		require.Equal(t, int64(1), fields.RegistrationID)
	})

	t.Run("InvalidType", func(t *testing.T) {
		wampMsg := []any{1, "invalid", map[string]any{}}
		index := 1
		var fields messages.Fields

		err := messages.ValidateRegistrationID(wampMsg, index, &fields)
		require.EqualError(t, err, "item at index 1 must be of type int64 but was string")
	})
}

func TestValidatePublicationID(t *testing.T) {
	t.Run("ValidType", func(t *testing.T) {
		wampMsg := []any{1, 1, map[string]any{}}
		index := 1
		var fields messages.Fields

		err := messages.ValidatePublicationID(wampMsg, index, &fields)
		require.NoError(t, err)
		require.Equal(t, int64(1), fields.PublicationID)
	})

	t.Run("InvalidType", func(t *testing.T) {
		wampMsg := []any{1, "invalid", map[string]any{}}
		index := 1
		var fields messages.Fields

		err := messages.ValidatePublicationID(wampMsg, index, &fields)
		require.EqualError(t, err, "item at index 1 must be of type int64 but was string")
	})
}

func TestValidateSubscriptionID(t *testing.T) {
	t.Run("ValidType", func(t *testing.T) {
		wampMsg := []any{1, 1, map[string]any{}}
		index := 1
		var fields messages.Fields

		err := messages.ValidateSubscriptionID(wampMsg, index, &fields)
		require.NoError(t, err)
		require.Equal(t, int64(1), fields.SubscriptionID)
	})

	t.Run("InvalidType", func(t *testing.T) {
		wampMsg := []any{1, "invalid", map[string]any{}}
		index := 1
		var fields messages.Fields

		err := messages.ValidateSubscriptionID(wampMsg, index, &fields)
		require.EqualError(t, err, "item at index 1 must be of type int64 but was string")
	})
}

func TestValidateSignature(t *testing.T) {
	t.Run("ValidType", func(t *testing.T) {
		wampMsg := []any{1, "abcd", map[string]any{}}
		index := 1
		var fields messages.Fields

		err := messages.ValidateSignature(wampMsg, index, &fields)
		require.NoError(t, err)
		require.Equal(t, "abcd", fields.Signature)
	})

	t.Run("InvalidType", func(t *testing.T) {
		wampMsg := []any{1, 1, map[string]any{}}
		index := 1
		var fields messages.Fields

		err := messages.ValidateSignature(wampMsg, index, &fields)
		require.EqualError(t, err, "item at index 1 must be of type string but was int")
	})
}

func TestValidateURI(t *testing.T) {
	t.Run("ValidType", func(t *testing.T) {
		wampMsg := []any{1, "io.xconn.test", map[string]any{}}
		index := 1
		var fields messages.Fields

		err := messages.ValidateURI(wampMsg, index, &fields)
		require.NoError(t, err)
		require.Equal(t, "io.xconn.test", fields.URI)
	})

	t.Run("InvalidType", func(t *testing.T) {
		wampMsg := []any{1, 1, map[string]any{}}
		index := 1
		var fields messages.Fields

		err := messages.ValidateURI(wampMsg, index, &fields)
		require.EqualError(t, err, "item at index 1 must be of type string but was int")
	})
}

func TestValidateRealm(t *testing.T) {
	t.Run("ValidType", func(t *testing.T) {
		wampMsg := []any{1, "io.xconn.test", map[string]any{}}
		index := 1
		var fields messages.Fields

		err := messages.ValidateRealm(wampMsg, index, &fields)
		require.NoError(t, err)
		require.Equal(t, "io.xconn.test", fields.Realm)
	})

	t.Run("InvalidType", func(t *testing.T) {
		wampMsg := []any{1, 1, map[string]any{}}
		index := 1
		var fields messages.Fields

		err := messages.ValidateRealm(wampMsg, index, &fields)
		require.EqualError(t, err, "item at index 1 must be of type string but was int")
	})
}

func TestValidateAuthMethod(t *testing.T) {
	t.Run("ValidType", func(t *testing.T) {
		wampMsg := []any{1, "anonymous", map[string]any{}}
		index := 1
		var fields messages.Fields

		err := messages.ValidateAuthMethod(wampMsg, index, &fields)
		require.NoError(t, err)
		require.Equal(t, "anonymous", fields.AuthMethod)
	})

	t.Run("InvalidType", func(t *testing.T) {
		wampMsg := []any{1, 1, map[string]any{}}
		index := 1
		var fields messages.Fields

		err := messages.ValidateAuthMethod(wampMsg, index, &fields)
		require.EqualError(t, err, "item at index 1 must be of type string but was int")
	})
}

func TestValidateReason(t *testing.T) {
	t.Run("ValidType", func(t *testing.T) {
		wampMsg := []any{1, "wamp.error.unknown", map[string]any{}}
		index := 1
		var fields messages.Fields

		err := messages.ValidateReason(wampMsg, index, &fields)
		require.NoError(t, err)
		require.Equal(t, "wamp.error.unknown", fields.Reason)
	})

	t.Run("InvalidType", func(t *testing.T) {
		wampMsg := []any{1, 1, map[string]any{}}
		index := 1
		var fields messages.Fields

		err := messages.ValidateReason(wampMsg, index, &fields)
		require.EqualError(t, err, "item at index 1 must be of type string but was int")
	})
}

func TestValidateExtra(t *testing.T) {
	t.Run("ValidType", func(t *testing.T) {
		wampMsg := []any{1, "io.xconn.test", map[string]any{"pubKey": "jhysggvvjygvb"}}
		index := 2
		var fields messages.Fields

		err := messages.ValidateExtra(wampMsg, index, &fields)
		require.NoError(t, err)
		require.Equal(t, map[string]any{"pubKey": "jhysggvvjygvb"}, fields.Extra)
	})

	t.Run("InvalidType", func(t *testing.T) {
		wampMsg := []any{1, "io.xconn.test", 1}
		index := 2
		var fields messages.Fields

		err := messages.ValidateExtra(wampMsg, index, &fields)
		require.EqualError(t, err, "item at index 2 must be of type map[string]any but was int")
	})
}

func TestValidateOptions(t *testing.T) {
	t.Run("ValidType", func(t *testing.T) {
		wampMsg := []any{1, "io.xconn.test", map[string]any{"invoke": "roundrobin"}}
		index := 2
		var fields messages.Fields

		err := messages.ValidateOptions(wampMsg, index, &fields)
		require.NoError(t, err)
		require.Equal(t, map[string]any{"invoke": "roundrobin"}, fields.Options)
	})

	t.Run("InvalidType", func(t *testing.T) {
		wampMsg := []any{1, "io.xconn.test", 1}
		index := 2
		var fields messages.Fields

		err := messages.ValidateOptions(wampMsg, index, &fields)
		require.EqualError(t, err, "item at index 2 must be of type map[string]any but was int")
	})
}

func TestValidateDetails(t *testing.T) {
	t.Run("ValidType", func(t *testing.T) {
		wampMsg := []any{1, "io.xconn.test", map[string]any{"subscription": 1}}
		index := 2
		var fields messages.Fields

		err := messages.ValidateDetails(wampMsg, index, &fields)
		require.NoError(t, err)
		require.Equal(t, map[string]any{"subscription": 1}, fields.Details)
	})

	t.Run("InvalidType", func(t *testing.T) {
		wampMsg := []any{1, "io.xconn.test", 1}
		index := 2
		var fields messages.Fields

		err := messages.ValidateDetails(wampMsg, index, &fields)
		require.EqualError(t, err, "item at index 2 must be of type map[string]any but was int")
	})
}

func TestValidateKwArgs(t *testing.T) {
	t.Run("NoKwargs", func(t *testing.T) {
		wampMsg := []any{1, "io.xconn.test", map[string]any{}}
		index := 4
		var fields messages.Fields

		err := messages.ValidateKwArgs(wampMsg, index, &fields)
		require.NoError(t, err)
		require.Nil(t, fields.KwArgs)
	})

	t.Run("ValidType", func(t *testing.T) {
		wampMsg := []any{1, "io.xconn.test", map[string]any{}, []string{}, map[string]any{"abc": 123}}
		index := 4
		var fields messages.Fields

		err := messages.ValidateDetails(wampMsg, index, &fields)
		require.NoError(t, err)
		require.Equal(t, map[string]any{"abc": 123}, fields.Details)
	})

	t.Run("InvalidType", func(t *testing.T) {
		wampMsg := []any{1, "io.xconn.test", map[string]any{}, []string{}, 1}
		index := 4
		var fields messages.Fields

		err := messages.ValidateDetails(wampMsg, index, &fields)
		require.EqualError(t, err, "item at index 4 must be of type map[string]any but was int")
	})
}

func TestValidateMessage(t *testing.T) {
	spec := messages.ValidationSpec{
		MinLength: 3,
		MaxLength: 5,
		Spec: messages.Spec{
			0: messages.ValidateSessionID,
			1: messages.ValidateURI,
			2: messages.ValidateOptions,
			3: messages.ValidateArgs,
			4: messages.ValidateKwArgs,
		},
	}

	t.Run("ValidMessage", func(t *testing.T) {
		wampMsg := []any{1, "io.xconn.test", map[string]any{}, []any{"abc", 1}}
		fields, err := messages.ValidateMessage(wampMsg, spec)

		require.NoError(t, err)
		require.NotNil(t, fields)

	})

	t.Run("InvalidMessageLength", func(t *testing.T) {
		wampMsg := []any{1}
		_, err := messages.ValidateMessage(wampMsg, spec)

		require.Error(t, err)
		require.Contains(t, err.Error(), "unexpected message length")
	})

	t.Run("ValidatorError", func(t *testing.T) {
		wampMsg := []any{1, "io.xconn.test", map[string]any{}, "invalidType"}
		_, err := messages.ValidateMessage(wampMsg, spec)

		require.EqualError(t, err, "item at index 3 must be of type []any but was string")
	})

	t.Run("MultipleErrors", func(t *testing.T) {
		wampMsg := []any{1, "io.xconn.test", map[string]any{}, "invalidType", "extra"}
		_, err := messages.ValidateMessage(wampMsg, spec)

		require.Contains(t, []string{
			`item at index 3 must be of type []any but was string
item at index 4 must be of type map[string]any but was string`,
			`item at index 4 must be of type map[string]any but was string
item at index 3 must be of type []any but was string`,
		}, err.Error())
	})
}

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
			result, ok := messages.AsInt64(test.input)
			require.True(t, ok)
			require.Equal(t, test.expected, result)
		}
	})

	t.Run("InvalidConversion", func(t *testing.T) {
		result, ok := messages.AsInt64("invalid")
		require.False(t, ok)
		require.Equal(t, int64(0), result)
	})
}

func TestAnysToStrings(t *testing.T) {
	t.Run("ValidConversion", func(t *testing.T) {
		input := []any{"foo", "bar", "helloo"}

		result, err := messages.AnysToStrings(input)
		require.NoError(t, err)
		require.Equal(t, []string{"foo", "bar", "helloo"}, result)
	})

	t.Run("InvalidConversion", func(t *testing.T) {
		input := []any{"foo", 123, "bar"}

		_, err := messages.AnysToStrings(input)
		require.Error(t, err)
		require.Contains(t, err.Error(), "element 123 is not a string")
	})
}
