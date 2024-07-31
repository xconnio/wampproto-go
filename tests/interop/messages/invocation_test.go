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

func invocationsEqual(msg1 *messages.Invocation, msg2 *messages.Invocation) bool {
	return msg1.RequestID() == msg2.RequestID() &&
		msg1.RegistrationID() == msg2.RegistrationID() &&
		reflect.DeepEqual(msg1.Details(), msg2.Details()) &&
		reflect.DeepEqual(msg1.Args(), msg2.Args()) &&
		reflect.DeepEqual(msg1.KwArgs(), msg2.KwArgs())
}

func testInvocationMessage(t *testing.T, serializerStr string, serializer serializers.Serializer) {
	var message = messages.NewInvocation(1, 1, map[string]any{"foo": true}, []any{"abc"}, map[string]any{"abc": "xyz"})
	command := fmt.Sprintf("message invocation 1 1 abc -d foo=true -k abc=xyz --serializer %s --output hex", serializerStr)

	msg := tests.RunCommandAndDeserialize(t, command, serializer)
	require.True(t, invocationsEqual(message, msg.(*messages.Invocation)))
}

func TestInvocationMessage(t *testing.T) {
	t.Run("JSONSerializer", func(t *testing.T) {
		serializer := &serializers.JSONSerializer{}
		testInvocationMessage(t, "json", serializer)
	})

	t.Run("CBORSerializer", func(t *testing.T) {
		serializer := &serializers.CBORSerializer{}
		testInvocationMessage(t, "cbor", serializer)
	})

	t.Run("MsgPackSerializer", func(t *testing.T) {
		serializer := &serializers.MsgPackSerializer{}
		testInvocationMessage(t, "msgpack", serializer)
	})
}
