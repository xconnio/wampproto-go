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

func welcomesEqual(msg1 *messages.Welcome, msg2 *messages.Welcome) bool {
	return msg1.SessionID() == msg2.SessionID() &&
		reflect.DeepEqual(msg1.Details(), msg2.Details())
}

func testWelcomeMessage(t *testing.T, serializerStr string, serializer serializers.Serializer) {
	var message = messages.NewWelcome(1, map[string]any{
		"roles":      map[string]any{"callee": true},
		"authrole":   "anonymous",
		"authid":     "foo",
		"authmethod": "anonymous",
		"authextra":  map[string]any{"foo": "bar"},
	})
	command := fmt.Sprintf("message welcome 1 --authmethod=anonymous --authid=foo --authrole=anonymous "+
		"--roles callee=true -e foo=bar --serializer %s --output hex", serializerStr)

	msg := tests.RunCommandAndDeserialize(t, command, serializer)
	require.True(t, welcomesEqual(message, msg.(*messages.Welcome)))
}

func TestWelcomeMessage(t *testing.T) {
	t.Run("JSONSerializer", func(t *testing.T) {
		serializer := &serializers.JSONSerializer{}
		testWelcomeMessage(t, "json", serializer)
	})

	t.Run("CBORSerializer", func(t *testing.T) {
		serializer := &serializers.CBORSerializer{}
		testWelcomeMessage(t, "cbor", serializer)
	})

	t.Run("MsgpackSerializer", func(t *testing.T) {
		serializer := &serializers.MsgPackSerializer{}
		testWelcomeMessage(t, "msgpack", serializer)
	})
}
