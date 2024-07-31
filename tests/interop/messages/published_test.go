package messages_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/xconnio/wampproto-go/messages"
	"github.com/xconnio/wampproto-go/serializers"
	"github.com/xconnio/wampproto-go/tests"
)

func publishedEqual(msg1 *messages.Published, msg2 *messages.Published) bool {
	return msg1.RequestID() == msg2.RequestID() &&
		msg1.PublicationID() == msg2.PublicationID()
}

func testPublishedMessage(t *testing.T, serializerStr string, serializer serializers.Serializer) {
	var message = messages.NewPublished(1, 1)
	command := fmt.Sprintf("message published 1 1 --serializer %s --output hex", serializerStr)

	msg := tests.RunCommandAndDeserialize(t, command, serializer)
	require.True(t, publishedEqual(message, msg.(*messages.Published)))
}

func TestPublishedMessage(t *testing.T) {
	t.Run("JSONSerializer", func(t *testing.T) {
		serializer := &serializers.JSONSerializer{}
		testPublishedMessage(t, "json", serializer)
	})

	t.Run("CBORSerializer", func(t *testing.T) {
		serializer := &serializers.CBORSerializer{}
		testPublishedMessage(t, "cbor", serializer)
	})

	t.Run("MsgPackSerializer", func(t *testing.T) {
		serializer := &serializers.MsgPackSerializer{}
		testPublishedMessage(t, "msgpack", serializer)
	})
}
