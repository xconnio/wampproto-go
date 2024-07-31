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

func interruptsEqual(msg1 *messages.Interrupt, msg2 *messages.Interrupt) bool {
	return msg1.RequestID() == msg2.RequestID() &&
		reflect.DeepEqual(msg1.Options(), msg2.Options())
}

func testInterruptMessage(t *testing.T, serializerStr string, serializer serializers.Serializer) {
	var message = messages.NewInterrupt(1, map[string]any{"foo": "bar"})
	command := fmt.Sprintf("message interrupt 1 -o foo=bar --serializer %s --output hex", serializerStr)

	msg := tests.RunCommandAndDeserialize(t, command, serializer)
	require.True(t, interruptsEqual(message, msg.(*messages.Interrupt)))
}

func TestInterruptMessage(t *testing.T) {
	t.Run("JSONSerializer", func(t *testing.T) {
		serializer := &serializers.JSONSerializer{}
		testInterruptMessage(t, "json", serializer)
	})

	t.Run("CBORSerializer", func(t *testing.T) {
		serializer := &serializers.CBORSerializer{}
		testInterruptMessage(t, "cbor", serializer)
	})

	t.Run("MsgpackSerializer", func(t *testing.T) {
		serializer := &serializers.MsgPackSerializer{}
		testInterruptMessage(t, "msgpack", serializer)
	})
}
