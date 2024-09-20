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
		details := wampproto.NewSessionDetails(1, "realm", "authid", "anonymous", false)
		err := dealer.AddSession(details)
		require.NoError(t, err)

		err = dealer.RemoveSession(details.ID())
		require.NoError(t, err)

		err = dealer.RemoveSession(details.ID())
		require.Error(t, err)
	})
}

func TestDealerRegisterUnregister(t *testing.T) {
	dealer := wampproto.NewDealer()

	callee := wampproto.NewSessionDetails(1, "realm", "authid", "anonymous", false)
	err := dealer.AddSession(callee)
	require.NoError(t, err)

	var registerationID int64

	t.Run("Register", func(t *testing.T) {
		register := messages.NewRegister(1, nil, "foo.bar")
		msg, err := dealer.ReceiveMessage(callee.ID(), register)
		require.NoError(t, err)
		require.NotNil(t, msg)
		require.Equal(t, msg.Recipient, callee.ID())
		require.Equal(t, messages.MessageTypeRegistered, msg.Message.Type())

		hasProcedure := dealer.HasProcedure("foo.bar")
		require.True(t, hasProcedure)
		registerationID = msg.Message.(*messages.Registered).RegistrationID()

		t.Run("DuplicateProcedure", func(t *testing.T) {
			register = messages.NewRegister(2, nil, "foo.bar")
			msg, err = dealer.ReceiveMessage(callee.ID(), register)
			require.NoError(t, err)
			require.NotNil(t, msg)
			require.Equal(t, msg.Recipient, callee.ID())
			require.Equal(t, messages.MessageTypeError, msg.Message.Type())
			errMsg := msg.Message.(*messages.Error)
			require.NotNil(t, errMsg)
			require.Equal(t, errMsg.URI(), "wamp.error.procedure_already_exists")
		})

		t.Run("InvalidSessionID", func(t *testing.T) {
			invalidRegister := messages.NewRegister(2, nil, "foo.bar")
			_, err = dealer.ReceiveMessage(5, invalidRegister)
			require.EqualError(t, err, "cannot register procedure for non-existent session 5")
		})
	})

	t.Run("Call", func(t *testing.T) {
		caller := wampproto.NewSessionDetails(2, "realm", "authid", "anonymous", false)
		err := dealer.AddSession(caller)
		require.NoError(t, err)

		call := messages.NewCall(3, map[string]any{}, "foo.bar", []any{"abc"}, nil)
		invWithRecipient, err := dealer.ReceiveMessage(caller.ID(), call)
		require.NoError(t, err)
		require.NotNil(t, invWithRecipient)
		require.Equal(t, callee.ID(), invWithRecipient.Recipient)
		require.Equal(t, messages.MessageTypeInvocation, invWithRecipient.Message.Type())

		// receive yield for invocation
		invocation := invWithRecipient.Message.(*messages.Invocation)
		yield := messages.NewYield(invocation.RequestID(), map[string]any{}, []any{"abc"}, nil)
		yieldWithRecipient, err := dealer.ReceiveMessage(caller.ID(), yield)
		require.NoError(t, err)
		require.NotNil(t, yieldWithRecipient)
		require.Equal(t, caller.ID(), yieldWithRecipient.Recipient)
		require.Equal(t, messages.MessageTypeResult, yieldWithRecipient.Message.Type())

		t.Run("NonExistingProcedure", func(t *testing.T) {
			invalidCallMessage := messages.NewCall(3, map[string]any{}, "invalid", []any{"abc"}, nil)
			errWithRecipient, err := dealer.ReceiveMessage(caller.ID(), invalidCallMessage)
			require.NoError(t, err)
			require.NotNil(t, errWithRecipient)
			require.Equal(t, caller.ID(), errWithRecipient.Recipient)
			require.Equal(t, errWithRecipient.Message.Type(), messages.MessageTypeError)
		})

		t.Run("InvalidYield", func(t *testing.T) {
			_, err = dealer.ReceiveMessage(5, yield)
			require.EqualError(t, err, "yield: not pending calls for session 5")
		})
	})

	t.Run("Unregister", func(t *testing.T) {
		unregister := messages.NewUnregister(callee.ID(), registerationID)
		unregWithRecipient, err := dealer.ReceiveMessage(callee.ID(), unregister)
		require.NoError(t, err)
		require.NotNil(t, unregWithRecipient)
		require.Equal(t, callee.ID(), unregWithRecipient.Recipient)
		require.Equal(t, messages.MessageTypeUnregistered, unregWithRecipient.Message.Type())

		hasProcedure := dealer.HasProcedure("foo.bar")
		require.False(t, hasProcedure)

		t.Run("InvalidRegistration", func(t *testing.T) {
			_, err = dealer.ReceiveMessage(callee.ID(), unregister)
			require.EqualError(t, err, "unregister: session 1 has no registration 1")
		})
	})
}

func TestProgressiveCallResults(t *testing.T) {
	dealer := wampproto.NewDealer()

	callee := wampproto.NewSessionDetails(1, "realm", "authid", "anonymous", false)
	caller := wampproto.NewSessionDetails(2, "realm", "authid", "anonymous", false)

	err := dealer.AddSession(callee)
	require.NoError(t, err)
	err = dealer.AddSession(caller)
	require.NoError(t, err)

	register := messages.NewRegister(1, nil, "foo.bar")
	_, err = dealer.ReceiveMessage(callee.ID(), register)
	require.NoError(t, err)

	call := messages.NewCall(caller.ID(), map[string]any{wampproto.OptionReceiveProgress: true}, "foo.bar", []any{}, nil)
	messageWithRecipient, err := dealer.ReceiveMessage(callee.ID(), call)
	require.NoError(t, err)
	require.Equal(t, callee.ID(), messageWithRecipient.Recipient)
	invocation := messageWithRecipient.Message.(*messages.Invocation)
	require.True(t, invocation.Details()[wampproto.OptionReceiveProgress].(bool))

	for i := 0; i < 10; i++ {
		yield := messages.NewYield(invocation.RequestID(), map[string]any{wampproto.OptionProgress: true}, []any{}, nil)
		messageWithRecipient, err = dealer.ReceiveMessage(callee.ID(), yield)
		require.NoError(t, err)
		require.Equal(t, callee.ID(), messageWithRecipient.Recipient)
		result := messageWithRecipient.Message.(*messages.Result)
		require.Equal(t, call.RequestID(), result.RequestID())
		require.True(t, result.Details()[wampproto.OptionProgress].(bool))
	}

	yield := messages.NewYield(invocation.RequestID(), map[string]any{}, []any{}, nil)
	messageWithRecipient, err = dealer.ReceiveMessage(callee.ID(), yield)
	require.NoError(t, err)
	require.Equal(t, callee.ID(), messageWithRecipient.Recipient)
	result := messageWithRecipient.Message.(*messages.Result)
	require.Equal(t, call.RequestID(), result.RequestID())
	progress, _ := result.Details()[wampproto.OptionReceiveProgress].(bool)
	require.False(t, progress)
}
