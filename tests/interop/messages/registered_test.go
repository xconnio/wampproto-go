package messages_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/xconnio/wampproto-go/messages"
	"github.com/xconnio/wampproto-go/serializers"
	"github.com/xconnio/wampproto-go/tests"
)

func registeredEqual(msg1 *messages.Registered, msg2 *messages.Registered) bool {
	return msg1.RequestID() == msg2.RequestID() &&
		msg1.RegistrationID() == msg2.RegistrationID()
}

func testRegisteredMessage(t *testing.T, serializerStr string, serializer serializers.Serializer) {
	var message = messages.NewRegistered(1, 1)
	command := fmt.Sprintf("message registered 1 1 --serializer %s --output hex", serializerStr)

	msg := tests.RunCommandAndDeserialize(t, command, serializer)
	require.True(t, registeredEqual(message, msg.(*messages.Registered)))
}

func TestRegisteredMessage(t *testing.T) {
	t.Run("JSONSerializer", func(t *testing.T) {
		serializer := &serializers.JSONSerializer{}
		testRegisteredMessage(t, "json", serializer)
	})

	t.Run("CBORSerializer", func(t *testing.T) {
		serializer := &serializers.CBORSerializer{}
		testRegisteredMessage(t, "cbor", serializer)
	})

	t.Run("MsgPackSerializer", func(t *testing.T) {
		serializer := &serializers.MsgPackSerializer{}
		testRegisteredMessage(t, "msgpack", serializer)
	})
}
