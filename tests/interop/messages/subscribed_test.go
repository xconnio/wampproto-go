package messages_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/xconnio/wampproto-go/messages"
	"github.com/xconnio/wampproto-go/serializers"
	"github.com/xconnio/wampproto-go/tests"
)

func subscribedEqual(msg1 *messages.Subscribed, msg2 *messages.Subscribed) bool {
	return msg1.RequestID() == msg2.RequestID() &&
		msg1.SubscriptionID() == msg2.SubscriptionID()
}

func testSubscribedMessage(t *testing.T, serializerStr string, serializer serializers.Serializer) {
	var message = messages.NewSubscribed(1, 1)
	command := fmt.Sprintf("message subscribed 1 1 --serializer %s --output hex", serializerStr)

	msg := tests.RunCommandAndDeserialize(t, command, serializer)
	require.True(t, subscribedEqual(message, msg.(*messages.Subscribed)))
}

func TestSubscribedMessage(t *testing.T) {
	t.Run("JSONSerializer", func(t *testing.T) {
		serializer := &serializers.JSONSerializer{}
		testSubscribedMessage(t, "json", serializer)
	})

	t.Run("CBORSerializer", func(t *testing.T) {
		serializer := &serializers.CBORSerializer{}
		testSubscribedMessage(t, "cbor", serializer)
	})

	t.Run("MsgPackSerializer", func(t *testing.T) {
		serializer := &serializers.MsgPackSerializer{}
		testSubscribedMessage(t, "msgpack", serializer)
	})
}
