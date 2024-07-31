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

func cancelsEqual(msg1 *messages.Cancel, msg2 *messages.Cancel) bool {
	return msg1.RequestID() == msg2.RequestID() &&
		reflect.DeepEqual(msg1.Options(), msg2.Options())
}

func testCancelMessage(t *testing.T, serializerStr string, serializer serializers.Serializer) {
	var message = messages.NewCancel(1, map[string]any{"foo": "bar"})
	command := fmt.Sprintf("message cancel 1 -o foo=bar --serializer %s --output hex", serializerStr)

	msg := tests.RunCommandAndDeserialize(t, command, serializer)
	require.True(t, cancelsEqual(message, msg.(*messages.Cancel)))
}

func TestCancelMessage(t *testing.T) {
	t.Run("JSONSerializer", func(t *testing.T) {
		serializer := &serializers.JSONSerializer{}
		testCancelMessage(t, "json", serializer)
	})

	t.Run("CBORSerializer", func(t *testing.T) {
		serializer := &serializers.CBORSerializer{}
		testCancelMessage(t, "cbor", serializer)
	})

	t.Run("MsgPackSerializer", func(t *testing.T) {
		serializer := &serializers.MsgPackSerializer{}
		testCancelMessage(t, "msgpack", serializer)
	})
}
