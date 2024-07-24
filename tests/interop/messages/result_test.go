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

func resultsEqual(msg1 *messages.Result, msg2 *messages.Result) bool {
	return msg1.RequestID() == msg2.RequestID() &&
		reflect.DeepEqual(msg1.Details(), msg2.Details()) &&
		reflect.DeepEqual(msg1.Args(), msg2.Args()) &&
		reflect.DeepEqual(msg1.KwArgs(), msg2.KwArgs())
}

func testResultMessage(t *testing.T, serializerStr string, serializer serializers.Serializer) {
	var message = messages.NewResult(1, map[string]any{"foo": true}, []any{"abc"}, map[string]any{"abc": "xyz"})
	command := fmt.Sprintf("message result 1 abc -d foo=true -k abc=xyz --serializer %s --output hex", serializerStr)

	msg := tests.RunCommandAndDeserialize(t, command, serializer)
	require.True(t, resultsEqual(message, msg.(*messages.Result)))
}

func TestResultMessage(t *testing.T) {
	t.Run("JSONSerializer", func(t *testing.T) {
		serializer := &serializers.JSONSerializer{}
		testResultMessage(t, "json", serializer)
	})

	t.Run("CBORSerializer", func(t *testing.T) {
		serializer := &serializers.CBORSerializer{}
		testResultMessage(t, "cbor", serializer)
	})

	t.Run("MsgPackSerializer", func(t *testing.T) {
		serializer := &serializers.MsgPackSerializer{}
		testResultMessage(t, "msgpack", serializer)
	})
}
