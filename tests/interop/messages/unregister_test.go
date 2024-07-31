package messages_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/xconnio/wampproto-go/messages"
	"github.com/xconnio/wampproto-go/serializers"
	"github.com/xconnio/wampproto-go/tests"
)

func unregistersEqual(msg1 *messages.Unregister, msg2 *messages.Unregister) bool {
	return msg1.RequestID() == msg2.RequestID() &&
		msg1.RegistrationID() == msg2.RegistrationID()
}

func testUnregisterMessage(t *testing.T, serializerStr string, serializer serializers.Serializer) {
	var message = messages.NewUnregister(1, 1)
	command := fmt.Sprintf("message unregister 1 1 --serializer %s --output hex", serializerStr)

	msg := tests.RunCommandAndDeserialize(t, command, serializer)
	require.True(t, unregistersEqual(message, msg.(*messages.Unregister)))
}

func TestUnregisterMessage(t *testing.T) {
	t.Run("JSONSerializer", func(t *testing.T) {
		serializer := &serializers.JSONSerializer{}
		testUnregisterMessage(t, "json", serializer)
	})

	t.Run("CBORSerializer", func(t *testing.T) {
		serializer := &serializers.CBORSerializer{}
		testUnregisterMessage(t, "cbor", serializer)
	})

	t.Run("MsgPackSerializer", func(t *testing.T) {
		serializer := &serializers.MsgPackSerializer{}
		testUnregisterMessage(t, "msgpack", serializer)
	})
}
