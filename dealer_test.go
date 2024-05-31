package wampproto_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/xconnio/wampproto-go"
	"github.com/xconnio/wampproto-go/messages"
)

func TestDealerAddRemoveSession(t *testing.T) {
	dealer := wampproto.NewDealer()

	t.Run("RemoveNonSession", func(t *testing.T) {
		err := dealer.RemoveSession(1)
		require.Error(t, err)
	})

	t.Run("AddRemove", func(t *testing.T) {
		details := wampproto.NewSessionDetails(1, "realm", "authid", "anonymous")
		err := dealer.AddSession(details)
		require.NoError(t, err)

		err = dealer.RemoveSession(details.ID())
		require.NoError(t, err)

		err = dealer.RemoveSession(details.ID())
		require.Error(t, err)
	})
}

func TestDealerRegisterUnRegister(t *testing.T) {
	dealer := wampproto.NewDealer()

	t.Run("Register", func(t *testing.T) {
		details := wampproto.NewSessionDetails(1, "realm", "authid", "anonymous")
		err := dealer.AddSession(details)
		require.NoError(t, err)

		register := messages.NewRegister(1, nil, "foo.bar")
		msg, err := dealer.ReceiveMessage(details.ID(), register)
		require.NoError(t, err)
		require.NotNil(t, msg)
		require.Equal(t, msg.Recipient, details.ID())
		require.Equal(t, msg.Message.Type(), messages.MessageTypeRegistered)

		register = messages.NewRegister(2, nil, "foo.bar")
		msg, err = dealer.ReceiveMessage(details.ID(), register)
		require.NoError(t, err)
		require.NotNil(t, msg)
		require.Equal(t, msg.Recipient, details.ID())
		require.Equal(t, msg.Message.Type(), messages.MessageTypeError)
		errMsg := msg.Message.(*messages.Error)
		require.NotNil(t, errMsg)
		require.Equal(t, errMsg.URI(), "wamp.error.procedure_already_exists")
	})
}
