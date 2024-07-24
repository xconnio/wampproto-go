package messages_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/xconnio/wampproto-go/messages"
	"github.com/xconnio/wampproto-go/serializers"
	"github.com/xconnio/wampproto-go/tests"
)

func unsubscribedEqual(msg1 *messages.Unsubscribed, msg2 *messages.Unsubscribed) bool {
	return msg1.RequestID() == msg2.RequestID()
}

func testUnsubscribedMessage(t *testing.T, serializerStr string, serializer serializers.Serializer) {
	var message = messages.NewUnsubscribed(1)
	command := fmt.Sprintf("message unsubscribed 1 --serializer %s --output hex", serializerStr)

	msg := tests.RunCommandAndDeserialize(t, command, serializer)
	require.True(t, unsubscribedEqual(message, msg.(*messages.Unsubscribed)))
}

func TestUnsubscribedMessage(t *testing.T) {
	t.Run("JSONSerializer", func(t *testing.T) {
		serializer := &serializers.JSONSerializer{}
		testUnsubscribedMessage(t, "json", serializer)
	})

	t.Run("CBORSerializer", func(t *testing.T) {
		serializer := &serializers.CBORSerializer{}
		testUnsubscribedMessage(t, "cbor", serializer)
	})

	t.Run("MsgPackSerializer", func(t *testing.T) {
		serializer := &serializers.MsgPackSerializer{}
		testUnsubscribedMessage(t, "msgpack", serializer)
	})
}
