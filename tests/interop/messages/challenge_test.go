package messages_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/xconnio/wampproto-go/messages"
	"github.com/xconnio/wampproto-go/serializers"
	"github.com/xconnio/wampproto-go/tests"
)

func challengesEqual(msg1 *messages.Challenge, msg2 *messages.Challenge) bool {
	return msg1.AuthMethod() == msg2.AuthMethod() &&
		reflect.DeepEqual(msg1.Extra(), msg2.Extra())
}

func testChallengeMessage(t *testing.T, serializerStr string, serializer serializers.Serializer) {
	var message = messages.NewChallenge("anonymous", map[string]any{"foo": "bar"})
	command := fmt.Sprintf("message challenge anonymous -e foo=bar --serializer %s --output hex", serializerStr)

	msg := tests.RunCommandAndDeserialize(t, command, serializer)
	require.True(t, challengesEqual(message, msg.(*messages.Challenge)))
}

func TestChallengeMessage(t *testing.T) {
	t.Run("JSONSerializer", func(t *testing.T) {
		serializer := &serializers.JSONSerializer{}
		testChallengeMessage(t, "json", serializer)
	})

	t.Run("CBORSerializer", func(t *testing.T) {
		serializer := &serializers.CBORSerializer{}
		testChallengeMessage(t, "cbor", serializer)
	})

	t.Run("MsgpackSerializer", func(t *testing.T) {
		serializer := &serializers.MsgPackSerializer{}
		testChallengeMessage(t, "msgpack", serializer)
	})
}
