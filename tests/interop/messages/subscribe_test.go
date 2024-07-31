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

func subscribesEqual(msg1 *messages.Subscribe, msg2 *messages.Subscribe) bool {
	return msg1.RequestID() == msg2.RequestID() &&
		msg1.Topic() == msg2.Topic() &&
		reflect.DeepEqual(msg1.Options(), msg2.Options())
}

func testSubscribeMessage(t *testing.T, serializerStr string, serializer serializers.Serializer) {
	var message = messages.NewSubscribe(1, map[string]any{"abc": "xyz"}, "test")
	command := fmt.Sprintf("message subscribe 1 test -o abc=xyz --serializer %s --output hex", serializerStr)

	msg := tests.RunCommandAndDeserialize(t, command, serializer)
	require.True(t, subscribesEqual(message, msg.(*messages.Subscribe)))
}

func TestSubscribeMessage(t *testing.T) {
	t.Run("JSONSerializer", func(t *testing.T) {
		serializer := &serializers.JSONSerializer{}
		testSubscribeMessage(t, "json", serializer)
	})

	t.Run("CBORSerializer", func(t *testing.T) {
		serializer := &serializers.CBORSerializer{}
		testSubscribeMessage(t, "cbor", serializer)
	})

	t.Run("MsgPackSerializer", func(t *testing.T) {
		serializer := &serializers.MsgPackSerializer{}
		testSubscribeMessage(t, "msgpack", serializer)
	})
}
