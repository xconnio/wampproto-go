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

func yieldsEqual(msg1 *messages.Yield, msg2 *messages.Yield) bool {
	return msg1.RequestID() == msg2.RequestID() &&
		reflect.DeepEqual(msg1.Options(), msg2.Options()) &&
		reflect.DeepEqual(msg1.Args(), msg2.Args()) &&
		reflect.DeepEqual(msg1.KwArgs(), msg2.KwArgs())
}

func testYieldMessage(t *testing.T, serializerStr string, serializer serializers.Serializer) {
	var message = messages.NewYield(1, map[string]any{"foo": true}, []any{"abc"}, map[string]any{"abc": "xyz"})
	command := fmt.Sprintf("message yield 1 abc -o foo=true -k abc=xyz --serializer %s --output hex", serializerStr)

	msg := tests.RunCommandAndDeserialize(t, command, serializer)
	require.True(t, yieldsEqual(message, msg.(*messages.Yield)))
}

func TestYieldMessage(t *testing.T) {
	t.Run("JSONSerializer", func(t *testing.T) {
		serializer := &serializers.JSONSerializer{}
		testYieldMessage(t, "json", serializer)
	})

	t.Run("CBORSerializer", func(t *testing.T) {
		serializer := &serializers.CBORSerializer{}
		testYieldMessage(t, "cbor", serializer)
	})

	t.Run("MsgPackSerializer", func(t *testing.T) {
		serializer := &serializers.MsgPackSerializer{}
		testYieldMessage(t, "msgpack", serializer)
	})
}
