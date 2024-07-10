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

func callsEqual(msg1 *messages.Call, msg2 *messages.Call) bool {
	return msg1.RequestID() == msg2.RequestID() &&
		msg1.Procedure() == msg2.Procedure() &&
		reflect.DeepEqual(msg1.Options(), msg2.Options()) &&
		reflect.DeepEqual(msg1.Args(), msg2.Args()) &&
		reflect.DeepEqual(msg1.KwArgs(), msg2.KwArgs())
}

func TestCallJSONSerializer(t *testing.T) {
	var message = messages.NewCall(1, map[string]any{}, "test", nil, nil)
	command := fmt.Sprintf("message call %v %s --serializer json --output hex", message.RequestID(), message.Procedure())
	var serializer = &serializers.JSONSerializer{}

	msg := tests.RunCommandAndDeserialize(t, command, serializer)
	require.True(t, callsEqual(message, msg.(*messages.Call)))
}

func TestCallCBORSerializer(t *testing.T) {
	var message = messages.NewCall(1, map[string]any{}, "test", []any{"abc"}, map[string]any{"abc": "xyz"})
	command := fmt.Sprintf("message call %v %s abc -k abc=xyz --serializer cbor --output hex",
		message.RequestID(), message.Procedure())
	var serializer = &serializers.CBORSerializer{}

	msg := tests.RunCommandAndDeserialize(t, command, serializer)
	require.True(t, callsEqual(message, msg.(*messages.Call)))
}

func TestCallMsgPackSerializer(t *testing.T) {
	var message = messages.NewCall(1, map[string]any{}, "test", []any{"abc"}, map[string]any{"abc": "xyz"})
	command := fmt.Sprintf("message call %v %s abc -k abc=xyz --serializer msgpack --output hex",
		message.RequestID(), message.Procedure())
	var serializer = &serializers.MsgPackSerializer{}

	msg := tests.RunCommandAndDeserialize(t, command, serializer)
	require.True(t, callsEqual(message, msg.(*messages.Call)))
}
