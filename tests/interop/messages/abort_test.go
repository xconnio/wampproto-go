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

func abortsEqual(msg1 *messages.Abort, msg2 *messages.Abort) bool {
	return msg1.Reason() == msg2.Reason() &&
		reflect.DeepEqual(msg1.Details(), msg2.Details()) &&
		reflect.DeepEqual(msg1.Args(), msg2.Args()) &&
		reflect.DeepEqual(msg1.KwArgs(), msg2.KwArgs())
}

func testAbortMessage(t *testing.T, serializerStr string, serializer serializers.Serializer) {
	var message = messages.NewAbort(map[string]any{}, "crash", []any{"abc"}, map[string]any{"abc": "xyz"})
	command := fmt.Sprintf("message abort crash abc -k abc=xyz --serializer %s --output hex", serializerStr)

	msg := tests.RunCommandAndDeserialize(t, command, serializer)
	require.True(t, abortsEqual(message, msg.(*messages.Abort)))
}

func TestAbortMessage(t *testing.T) {
	t.Run("JSONSerializer", func(t *testing.T) {
		serializer := &serializers.JSONSerializer{}
		testAbortMessage(t, "json", serializer)
	})

	t.Run("CBORSerializer", func(t *testing.T) {
		serializer := &serializers.CBORSerializer{}
		testAbortMessage(t, "cbor", serializer)
	})

	t.Run("MsgpackSerializer", func(t *testing.T) {
		serializer := &serializers.MsgPackSerializer{}
		testAbortMessage(t, "msgpack", serializer)
	})
}
