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

func callsEqual(msg1 *messages.Call, msg2 *messages.Call) bool {
	return msg1.RequestID() == msg2.RequestID() &&
		msg1.Procedure() == msg2.Procedure() &&
		reflect.DeepEqual(msg1.Options(), msg2.Options()) &&
		reflect.DeepEqual(msg1.Args(), msg2.Args()) &&
		reflect.DeepEqual(msg1.KwArgs(), msg2.KwArgs())
}

func testCallMessage(t *testing.T, serializerStr string, serializer serializers.Serializer) {
	var message = messages.NewCall(1, map[string]any{}, "test", []any{"abc"}, map[string]any{"abc": "xyz"})
	command := fmt.Sprintf("message call 1 test abc -k abc=xyz --serializer %s --output hex", serializerStr)

	msg := tests.RunCommandAndDeserialize(t, command, serializer)
	require.True(t, callsEqual(message, msg.(*messages.Call)))
}

func TestCallMessage(t *testing.T) {
	t.Run("JSONSerializer", func(t *testing.T) {
		serializer := &serializers.JSONSerializer{}
		testCallMessage(t, "json", serializer)
	})

	t.Run("CBORSerializer", func(t *testing.T) {
		serializer := &serializers.CBORSerializer{}
		testCallMessage(t, "cbor", serializer)
	})

	t.Run("MsgPackSerializer", func(t *testing.T) {
		serializer := &serializers.MsgPackSerializer{}
		testCallMessage(t, "msgpack", serializer)
	})
}
