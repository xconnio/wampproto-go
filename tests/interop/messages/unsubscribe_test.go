package messages_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/xconnio/wampproto-go/messages"
	"github.com/xconnio/wampproto-go/serializers"
	"github.com/xconnio/wampproto-go/tests"
)

func unsubscribeEqual(msg1 *messages.Unsubscribe, msg2 *messages.Unsubscribe) bool {
	return msg1.RequestID() == msg2.RequestID() &&
		msg1.SubscriptionID() == msg2.SubscriptionID()
}

func testUnsubscribeMessage(t *testing.T, serializerStr string, serializer serializers.Serializer) {
	var message = messages.NewUnsubscribe(1, 1)
	command := fmt.Sprintf("message unsubscribe 1 1 --serializer %s --output hex", serializerStr)

	msg := tests.RunCommandAndDeserialize(t, command, serializer)
	require.True(t, unsubscribeEqual(message, msg.(*messages.Unsubscribe)))
}

func TestUnsubscribeMessage(t *testing.T) {
	t.Run("JSONSerializer", func(t *testing.T) {
		serializer := &serializers.JSONSerializer{}
		testUnsubscribeMessage(t, "json", serializer)
	})

	t.Run("CBORSerializer", func(t *testing.T) {
		serializer := &serializers.CBORSerializer{}
		testUnsubscribeMessage(t, "cbor", serializer)
	})

	t.Run("MsgPackSerializer", func(t *testing.T) {
		serializer := &serializers.MsgPackSerializer{}
		testUnsubscribeMessage(t, "msgpack", serializer)
	})
}
