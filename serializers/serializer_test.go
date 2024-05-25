package serializers_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/xconnio/wampproto-go/messages"
	"github.com/xconnio/wampproto-go/serializers"
)

func serializeDeserialize(t *testing.T, serializer serializers.Serializer) {
	details := map[string]any{}
	reason := "hello"
	var args []any
	kwArgs := map[string]any{}
	message := messages.NewAbort(details, reason, args, kwArgs)

	data, err := serializer.Serialize(message)
	require.NoError(t, err)
	require.NotNil(t, data)

	deserialized, err := serializer.Deserialize(data)
	require.NoError(t, err)
	require.NotNil(t, deserialized)
	abort := deserialized.(messages.Abort)

	require.Equal(t, message.Reason(), abort.Reason())
	require.Equal(t, message.Args(), abort.Args())
	require.Equal(t, message.KwArgs(), abort.KwArgs())
	require.Equal(t, message.Details(), abort.Details())
}

func TestSerializers(t *testing.T) {
	t.Run("CBOR", func(t *testing.T) {
		serializer := &serializers.CBORSerializer{}
		serializeDeserialize(t, serializer)
	})

	t.Run("MSGPACK", func(t *testing.T) {
		serializer := &serializers.MsgPackSerializer{}
		serializeDeserialize(t, serializer)
	})

	t.Run("JSON", func(t *testing.T) {
		serializer := &serializers.JSONSerializer{}
		serializeDeserialize(t, serializer)
	})
}
