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
	})
}
