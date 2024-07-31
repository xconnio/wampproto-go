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

func hellosEqual(msg1 *messages.Hello, msg2 *messages.Hello) bool {
	return msg1.AuthID() == msg2.AuthID() &&
		msg1.Realm() == msg2.Realm() &&
		reflect.DeepEqual(msg1.AuthMethods(), msg2.AuthMethods()) &&
		reflect.DeepEqual(msg1.AuthExtra(), msg2.AuthExtra()) &&
		reflect.DeepEqual(msg1.Roles(), msg2.Roles())
}

func testHelloMessage(t *testing.T, serializerStr string, serializer serializers.Serializer) {
	var message = messages.NewHello("realm1", "foo", map[string]any{"foo": "bar"},
		map[string]any{"callee": true}, []string{"anonymous"})
	command := fmt.Sprintf(
		"message hello realm1 anonymous --authid foo -r callee=true -e foo:bar --serializer %s --output hex",
		serializerStr)

	msg := tests.RunCommandAndDeserialize(t, command, serializer)
	require.True(t, hellosEqual(message, msg.(*messages.Hello)))
}

func TestHelloMessage(t *testing.T) {
	t.Run("JSONSerializer", func(t *testing.T) {
		serializer := &serializers.JSONSerializer{}
		testHelloMessage(t, "json", serializer)
	})

	t.Run("CBORSerializer", func(t *testing.T) {
		serializer := &serializers.CBORSerializer{}
		testHelloMessage(t, "cbor", serializer)
	})

	t.Run("MsgpackSerializer", func(t *testing.T) {
		serializer := &serializers.MsgPackSerializer{}
		testHelloMessage(t, "msgpack", serializer)
	})
}
