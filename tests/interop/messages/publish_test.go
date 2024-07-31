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

func publishEqual(msg1 *messages.Publish, msg2 *messages.Publish) bool {
	return msg1.RequestID() == msg2.RequestID() &&
		msg1.Topic() == msg2.Topic() &&
		reflect.DeepEqual(msg1.Options(), msg2.Options()) &&
		reflect.DeepEqual(msg1.Args(), msg2.Args()) &&
		reflect.DeepEqual(msg1.KwArgs(), msg2.KwArgs())
}

func testPublishMessage(t *testing.T, serializerStr string, serializer serializers.Serializer) {
	var message = messages.NewPublish(1, map[string]any{"abc": "xyz"}, "test", []any{"abc"}, map[string]any{"abc": "xyz"})
	command := fmt.Sprintf("message publish 1 test abc -o abc:xyz -k abc=xyz --serializer %s --output hex", serializerStr)

	msg := tests.RunCommandAndDeserialize(t, command, serializer)
	require.True(t, publishEqual(message, msg.(*messages.Publish)))
}

func TestPublishMessage(t *testing.T) {
	t.Run("JSONSerializer", func(t *testing.T) {
		serializer := &serializers.JSONSerializer{}
		testPublishMessage(t, "json", serializer)
	})

	t.Run("CBORSerializer", func(t *testing.T) {
		serializer := &serializers.CBORSerializer{}
		testPublishMessage(t, "cbor", serializer)
	})

	t.Run("MsgPackSerializer", func(t *testing.T) {
		serializer := &serializers.MsgPackSerializer{}
		testPublishMessage(t, "msgpack", serializer)
	})
}
