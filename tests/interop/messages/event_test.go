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

func eventsEqual(msg1 *messages.Event, msg2 *messages.Event) bool {
	return msg1.SubscriptionID() == msg2.SubscriptionID() &&
		msg1.PublicationID() == msg2.PublicationID() &&
		reflect.DeepEqual(msg1.Details(), msg2.Details()) &&
		reflect.DeepEqual(msg1.Args(), msg2.Args()) &&
		reflect.DeepEqual(msg1.KwArgs(), msg2.KwArgs())
}

func testEventMessage(t *testing.T, serializerStr string, serializer serializers.Serializer) {
	var message = messages.NewEvent(1, 1, map[string]any{"foo": true}, []any{"abc"}, map[string]any{"abc": "xyz"})
	command := fmt.Sprintf("message event 1 1 abc -d foo=true -k abc=xyz --serializer %s --output hex", serializerStr)

	msg := tests.RunCommandAndDeserialize(t, command, serializer)
	require.True(t, eventsEqual(message, msg.(*messages.Event)))
}

func TestEventMessage(t *testing.T) {
	t.Run("JSONSerializer", func(t *testing.T) {
		serializer := &serializers.JSONSerializer{}
		testEventMessage(t, "json", serializer)
	})

	t.Run("CBORSerializer", func(t *testing.T) {
		serializer := &serializers.CBORSerializer{}
		testEventMessage(t, "cbor", serializer)
	})

	t.Run("MsgSerializer", func(t *testing.T) {
		serializer := &serializers.MsgPackSerializer{}
		testEventMessage(t, "msgpack", serializer)
	})
}
