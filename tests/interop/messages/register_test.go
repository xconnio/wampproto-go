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

func registersEqual(msg1 *messages.Register, msg2 *messages.Register) bool {
	return msg1.RequestID() == msg2.RequestID() &&
		msg1.Procedure() == msg2.Procedure() &&
		reflect.DeepEqual(msg1.Options(), msg2.Options())
}

func testRegisterMessage(t *testing.T, serializerStr string, serializer serializers.Serializer) {
	var message = messages.NewRegister(1, map[string]any{"abc": "xyz"}, "test")
	command := fmt.Sprintf("message register 1 test -o abc=xyz --serializer %s --output hex", serializerStr)

	msg := tests.RunCommandAndDeserialize(t, command, serializer)
	require.True(t, registersEqual(message, msg.(*messages.Register)))
}

func TestRegisterMessage(t *testing.T) {
	t.Run("JSONSerializer", func(t *testing.T) {
		serializer := &serializers.JSONSerializer{}
		testRegisterMessage(t, "json", serializer)
	})

	t.Run("CBORSerializer", func(t *testing.T) {
		serializer := &serializers.CBORSerializer{}
		testRegisterMessage(t, "cbor", serializer)
	})

	t.Run("MsgPackSerializer", func(t *testing.T) {
		serializer := &serializers.MsgPackSerializer{}
		testRegisterMessage(t, "msgpack", serializer)
	})
}
