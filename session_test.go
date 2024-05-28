package wampproto_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/xconnio/wampproto-go"
	"github.com/xconnio/wampproto-go/messages"
	"github.com/xconnio/wampproto-go/serializers"
)

func registerProc(t *testing.T, callee *wampproto.Session, uri string) {
	// send register message
	register := messages.NewRegister(1, nil, uri)
	payload, err := callee.SendMessage(register)
	require.NoError(t, err)
	require.NotNil(t, payload)

	// confirm registration
	registered := messages.NewRegistered(1, 1)
	reg, err := callee.ReceiveMessage(registered)
	require.NoError(t, err)
	require.NotNil(t, reg)
}

func callProc(t *testing.T, caller, callee *wampproto.Session, uri string) {
	callRequest := int64(2)

	call := messages.NewCall(messages.NewCallFields(callRequest, nil, uri, nil, nil))
	payload, err := caller.SendMessage(call)
	require.NoError(t, err)
	require.NotNil(t, payload)

	// send invocation to the callee
	invocation := messages.NewInvocation(messages.NewInvocationFields(callRequest, 1, nil, nil, nil))
	toSend, err := callee.ReceiveMessage(invocation)
	require.NoError(t, err)
	require.NotNil(t, toSend)

	// send yield to the caller
	result := messages.NewResult(messages.NewResultFields(callRequest, nil, nil, nil))
	rslt, err := caller.ReceiveMessage(result)
	require.NoError(t, err)
	require.NotNil(t, rslt)
	require.Equal(t, rslt, result)
}

func registerAndCall(t *testing.T, procedure string, serializer serializers.Serializer) {
	caller := wampproto.NewSession(serializer)
	callee := wampproto.NewSession(serializer)

	registerProc(t, callee, procedure)
	callProc(t, caller, callee, procedure)
}

func TestSessionCall(t *testing.T) {
	procedure := "foo.bar"
	t.Run("JSON", func(t *testing.T) {
		serializer := &serializers.JSONSerializer{}
		registerAndCall(t, procedure, serializer)
	})

	t.Run("CBOR", func(t *testing.T) {
		serializer := &serializers.CBORSerializer{}
		registerAndCall(t, procedure, serializer)
	})

	t.Run("MSGPACK", func(t *testing.T) {
		serializer := &serializers.MsgPackSerializer{}
		registerAndCall(t, procedure, serializer)
	})
}

func subscribeTopic(t *testing.T, subscriber *wampproto.Session, uri string) {
	subscribeID := int64(1)
	subscribe := messages.NewSubscribe(subscribeID, nil, uri)
	payload, err := subscriber.SendMessage(subscribe)
	require.NoError(t, err)
	require.NotNil(t, payload)

	subscribed := messages.NewSubscribed(subscribeID, 1)
	_, err = subscriber.ReceiveMessage(subscribed)
	require.NoError(t, err)
}

func publishTopic(t *testing.T, publisher, subscriber *wampproto.Session, uri string) {
	publish := messages.NewPublish(2, uri, nil, nil)
	_, err := publisher.SendMessage(publish)
	require.NoError(t, err)

	event := messages.NewEvent(1, 2, nil, nil, nil)
	_, err = subscriber.ReceiveMessage(event)
	require.NoError(t, err)
}

func subscribeAndPublish(t *testing.T, topic string, serializer serializers.Serializer) {
	publisher := wampproto.NewSession(serializer)
	subscriber := wampproto.NewSession(serializer)

	subscribeTopic(t, subscriber, topic)
	publishTopic(t, publisher, subscriber, topic)
}

func TestSessionPublish(t *testing.T) {
	topic := "foo.bar"
	t.Run("JSON", func(t *testing.T) {
		serializer := &serializers.JSONSerializer{}
		subscribeAndPublish(t, topic, serializer)
	})

	t.Run("CBOR", func(t *testing.T) {
		serializer := &serializers.CBORSerializer{}
		subscribeAndPublish(t, topic, serializer)
	})

	t.Run("MSGPACK", func(t *testing.T) {
		serializer := &serializers.MsgPackSerializer{}
		subscribeAndPublish(t, topic, serializer)
	})
}
