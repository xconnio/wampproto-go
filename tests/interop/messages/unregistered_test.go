package messages_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/xconnio/wampproto-go/messages"
	"github.com/xconnio/wampproto-go/serializers"
	"github.com/xconnio/wampproto-go/tests"
)

func unregisteredEqual(msg1 *messages.Unregistered, msg2 *messages.Unregistered) bool {
	return msg1.RequestID() == msg2.RequestID()
}

func testUnregisteredMessage(t *testing.T, serializerStr string, serializer serializers.Serializer) {
	var message = messages.NewUnregistered(1)
	command := fmt.Sprintf("message unregistered 1 --serializer %s --output hex", serializerStr)

	msg := tests.RunCommandAndDeserialize(t, command, serializer)
	require.True(t, unregisteredEqual(message, msg.(*messages.Unregistered)))
}

func TestUnregisteredMessage(t *testing.T) {
	t.Run("JSONSerializer", func(t *testing.T) {
		serializer := &serializers.JSONSerializer{}
		testUnregisteredMessage(t, "json", serializer)
	})

	t.Run("CBORSerializer", func(t *testing.T) {
		serializer := &serializers.CBORSerializer{}
		testUnregisteredMessage(t, "cbor", serializer)
	})

	t.Run("MsgpackSerializer", func(t *testing.T) {
		serializer := &serializers.MsgPackSerializer{}
		testUnregisteredMessage(t, "msgpack", serializer)
	})
}
