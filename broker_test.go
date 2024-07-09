package wampproto_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/xconnio/wampproto-go"
	"github.com/xconnio/wampproto-go/messages"
)

func TestBrokerAddRemoveSession(t *testing.T) {
	broker := wampproto.NewBroker()

	t.Run("RemoveNonSession", func(t *testing.T) {
		err := broker.RemoveSession(1)
		require.Error(t, err)
	})

	t.Run("AddRemove", func(t *testing.T) {
		details := wampproto.NewSessionDetails(1, "realm", "authid", "anonymous", false)
		err := broker.AddSession(details)
		require.NoError(t, err)

		err = broker.RemoveSession(details.ID())
		require.NoError(t, err)

		err = broker.RemoveSession(details.ID())
		require.Error(t, err)
	})
}

func TestBrokerPublish(t *testing.T) {
	broker := wampproto.NewBroker()

	details := wampproto.NewSessionDetails(1, "realm", "authid", "anonymous", false)
	err := broker.AddSession(details)
	require.NoError(t, err)

	args := []any{1, 2}
	kwArgs := map[string]any{"name": "alex"}
	options := map[string]any{wampproto.OptAcknowledge: true}

	t.Run("NoSubscriber", func(t *testing.T) {
		publish := messages.NewPublish(1, options, "foo.bar", args, kwArgs)
		publication, err := broker.ReceivePublish(details.ID(), publish)
		require.NoError(t, err)
		require.NotNil(t, publication)

		require.Equal(t, publication.Ack.Recipient, details.ID())
		require.Equal(t, publication.Ack.Message.Type(), messages.MessageTypePublished)
		require.Nil(t, publication.Event)
		require.Len(t, publication.Recipients, 0)
	})

	t.Run("WithSubscriber", func(t *testing.T) {
		subDetails := wampproto.NewSessionDetails(2, "realm", "authid", "anonymous", false)
		err = broker.AddSession(subDetails)
		require.NoError(t, err)

		subscribe := messages.NewSubscribe(2, nil, "foo.bar")
		msgWithRecipient, err := broker.ReceiveMessage(subDetails.ID(), subscribe)
		require.NoError(t, err)
		require.Equal(t, msgWithRecipient.Recipient, subDetails.ID())
		require.Equal(t, msgWithRecipient.Message.Type(), messages.MessageTypeSubscribed)

		publish := messages.NewPublish(3, options, "foo.bar", args, kwArgs)
		publication, err := broker.ReceivePublish(details.ID(), publish)
		require.NoError(t, err)
		require.NotNil(t, publication)

		require.Equal(t, publication.Ack.Recipient, details.ID())
		require.Equal(t, publication.Ack.Message.Type(), messages.MessageTypePublished)
		require.NotNil(t, publication.Event)
		require.Len(t, publication.Recipients, 1)
	})

	t.Run("WithoutAcknowledge", func(t *testing.T) {
		publish := messages.NewPublish(4, map[string]any{}, "foo.bar", args, kwArgs)
		publication, err := broker.ReceivePublish(details.ID(), publish)
		require.NoError(t, err)
		require.NotNil(t, publication)
		require.Nil(t, publication.Ack)
	})

	t.Run("InvalidSessionID", func(t *testing.T) {
		publish := messages.NewPublish(1, options, "foo.bar", args, kwArgs)
		_, err = broker.ReceivePublish(5, publish)
		require.EqualError(t, err, "broker: cannot publish, session 5 doesn't exist")
	})
}

func TestBrokerSubscribeUnsubscribe(t *testing.T) {
	broker := wampproto.NewBroker()

	subDetails := wampproto.NewSessionDetails(1, "realm", "authid", "anonymous", false)
	err := broker.AddSession(subDetails)
	require.NoError(t, err)

	var subscriptionID int64
	t.Run("Subscribe", func(t *testing.T) {
		subscribe := messages.NewSubscribe(1, nil, "foo.bar")
		msgWithRecipient, err := broker.ReceiveMessage(subDetails.ID(), subscribe)
		require.NoError(t, err)
		require.Equal(t, msgWithRecipient.Recipient, subDetails.ID())
		require.Equal(t, msgWithRecipient.Message.Type(), messages.MessageTypeSubscribed)

		subscriptionID = msgWithRecipient.Message.(*messages.Subscribed).SubscriptionID()
	})

	t.Run("PublishAndReceiveEvent", func(t *testing.T) {
		pubDetails := wampproto.NewSessionDetails(2, "realm", "authid", "anonymous", false)
		err = broker.AddSession(pubDetails)
		require.NoError(t, err)

		publish := messages.NewPublish(2, map[string]any{wampproto.OptAcknowledge: true}, "foo.bar", []any{1, 2}, nil)
		publication, err := broker.ReceivePublish(pubDetails.ID(), publish)
		require.NoError(t, err)
		require.NotNil(t, publication)

		require.Equal(t, publication.Ack.Recipient, pubDetails.ID())
		require.Equal(t, publication.Ack.Message.Type(), messages.MessageTypePublished)
		require.NotNil(t, publication.Event)
		require.Len(t, publication.Recipients, 1)
	})

	t.Run("Unsubscribe", func(t *testing.T) {
		unsubscribe := messages.NewUnsubscribe(3, subscriptionID)

		msgWithRecipient, err := broker.ReceiveMessage(subDetails.ID(), unsubscribe)
		require.NoError(t, err)
		require.Equal(t, msgWithRecipient.Recipient, subDetails.ID())
		require.Equal(t, msgWithRecipient.Message.Type(), messages.MessageTypeUnsubscribed)
	})

	t.Run("SubscribeInvalidSessionID", func(t *testing.T) {
		subscribe := messages.NewSubscribe(4, nil, "foo.bar")
		_, err = broker.ReceiveMessage(5, subscribe)
		require.EqualError(t, err, "broker: cannot subscribe, session 5 doesn't exist")
	})

	t.Run("UnsubscribeInvalidSessionID", func(t *testing.T) {
		unsubscribe := messages.NewUnsubscribe(5, subscriptionID)

		_, err = broker.ReceiveMessage(5, unsubscribe)
		require.EqualError(t, err, "broker: cannot unsubscribe, session 5 doesn't exist")
	})

	t.Run("UnsubscribeInvalidSubscriptionID", func(t *testing.T) {
		unsubscribe := messages.NewUnsubscribe(3, 5)

		_, err = broker.ReceiveMessage(subDetails.ID(), unsubscribe)
		require.EqualError(t, err, "broker: cannot unsubscribe non-existent subscription 5")
	})
}
