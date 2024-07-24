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

func errorsEqual(msg1 *messages.Error, msg2 *messages.Error) bool {
	return msg1.RequestID() == msg2.RequestID() &&
		msg1.URI() == msg2.URI() &&
		msg1.MessageType() == msg2.MessageType() &&
		reflect.DeepEqual(msg1.Details(), msg2.Details()) &&
		reflect.DeepEqual(msg1.Args(), msg2.Args()) &&
		reflect.DeepEqual(msg1.KwArgs(), msg2.KwArgs())
}

func testErrorMessage(t *testing.T, serializerStr string, serializer serializers.Serializer) {
	var message = messages.NewError(1, 1, map[string]any{"foo": "bar"}, "test", []any{"abc"}, map[string]any{"abc": "xyz"})
	command := fmt.Sprintf("message error 1 1 test abc -k abc=xyz -d foo=bar --serializer %s --output hex", serializerStr)

	msg := tests.RunCommandAndDeserialize(t, command, serializer)
	require.True(t, errorsEqual(message, msg.(*messages.Error)))
}

func TestErrorMessage(t *testing.T) {
	t.Run("JSONSerializer", func(t *testing.T) {
		serializer := &serializers.JSONSerializer{}
		testErrorMessage(t, "json", serializer)
	})

	t.Run("CBORSerializer", func(t *testing.T) {
		serializer := &serializers.CBORSerializer{}
		testErrorMessage(t, "cbor", serializer)
	})

	t.Run("MsgpackSerializer", func(t *testing.T) {
		serializer := &serializers.MsgPackSerializer{}
		testErrorMessage(t, "msgpack", serializer)
	})
}
