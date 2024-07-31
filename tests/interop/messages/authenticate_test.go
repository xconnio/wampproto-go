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

func authenticatesEqual(msg1 *messages.Authenticate, msg2 *messages.Authenticate) bool {
	return msg1.Signature() == msg2.Signature() &&
		reflect.DeepEqual(msg1.Extra(), msg2.Extra())
}

func testAuthenticateMessage(t *testing.T, serializerStr string, serializer serializers.Serializer) {
	var message = messages.NewAuthenticate("anonymous", map[string]any{"foo": "bar"})
	command := fmt.Sprintf("message authenticate anonymous -e foo=bar --serializer %s --output hex", serializerStr)

	msg := tests.RunCommandAndDeserialize(t, command, serializer)
	require.True(t, authenticatesEqual(message, msg.(*messages.Authenticate)))
}

func TestAuthenticateMessage(t *testing.T) {
	t.Run("JSONSerializer", func(t *testing.T) {
		serializer := &serializers.JSONSerializer{}
		testAuthenticateMessage(t, "json", serializer)
	})

	t.Run("CBORSerializer", func(t *testing.T) {
		serializer := &serializers.CBORSerializer{}
		testAuthenticateMessage(t, "cbor", serializer)
	})

	t.Run("MsgpackSerializer", func(t *testing.T) {
		serializer := &serializers.MsgPackSerializer{}
		testAuthenticateMessage(t, "msgpack", serializer)
	})
}
