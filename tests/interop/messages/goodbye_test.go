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

func goodByesEqual(msg1 *messages.GoodBye, msg2 *messages.GoodBye) bool {
	return msg1.Reason() == msg2.Reason() &&
		reflect.DeepEqual(msg1.Details(), msg2.Details())
}

func testGoodByeMessage(t *testing.T, serializerStr string, serializer serializers.Serializer) {
	var message = messages.NewGoodBye("disconnect", map[string]any{"foo": "bar"})
	command := fmt.Sprintf("message goodbye disconnect -d foo=bar --serializer %s --output hex", serializerStr)

	msg := tests.RunCommandAndDeserialize(t, command, serializer)
	require.True(t, goodByesEqual(message, msg.(*messages.GoodBye)))
}

func TestGoodByeMessage(t *testing.T) {
	t.Run("JSONSerializer", func(t *testing.T) {
		serializer := &serializers.JSONSerializer{}
		testGoodByeMessage(t, "json", serializer)
	})

	t.Run("CBORSerializer", func(t *testing.T) {
		serializer := &serializers.CBORSerializer{}
		testGoodByeMessage(t, "cbor", serializer)
	})

	t.Run("MsgPackSerializer", func(t *testing.T) {
		serializer := &serializers.MsgPackSerializer{}
		testGoodByeMessage(t, "msgpack", serializer)
	})
}
