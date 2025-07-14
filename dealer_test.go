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
		msgs, err := dealer.ReceiveMessage(callee.ID(), register)
		require.NoError(t, err)
		require.NotNil(t, msgs)
		require.Len(t, msgs, 1)
		require.Equal(t, msgs[0].Recipient, callee.ID())
		require.Equal(t, messages.MessageTypeRegistered, msgs[0].Message.Type())

		hasProcedure := dealer.HasProcedure("foo.bar")
		require.True(t, hasProcedure)
		registerationID = msgs[0].Message.(*messages.Registered).RegistrationID()

		t.Run("DuplicateProcedure", func(t *testing.T) {
			register = messages.NewRegister(2, nil, "foo.bar")
			msgs, err = dealer.ReceiveMessage(callee.ID(), register)
			require.NoError(t, err)
			require.NotNil(t, msgs)
			require.Len(t, msgs, 1)
			require.Equal(t, msgs[0].Recipient, callee.ID())
			require.Equal(t, messages.MessageTypeError, msgs[0].Message.Type())
			errMsg := msgs[0].Message.(*messages.Error)
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
		msgs, err := dealer.ReceiveMessage(caller.ID(), call)
		require.NoError(t, err)
		require.Len(t, msgs, 1)
		invWithRecipient := msgs[0]
		require.Equal(t, callee.ID(), invWithRecipient.Recipient)
		require.Equal(t, messages.MessageTypeInvocation, invWithRecipient.Message.Type())

		// receive yield for invocation
		invocation := invWithRecipient.Message.(*messages.Invocation)
		yield := messages.NewYield(invocation.RequestID(), map[string]any{}, []any{"abc"}, nil)
		msgs, err = dealer.ReceiveMessage(caller.ID(), yield)
		require.NoError(t, err)
		require.Len(t, msgs, 1)
		yieldWithRecipient := msgs[0]
		require.Equal(t, caller.ID(), yieldWithRecipient.Recipient)
		require.Equal(t, messages.MessageTypeResult, yieldWithRecipient.Message.Type())

		t.Run("NonExistingProcedure", func(t *testing.T) {
			invalidCallMessage := messages.NewCall(3, map[string]any{}, "invalid", []any{"abc"}, nil)
			msgs, err = dealer.ReceiveMessage(caller.ID(), invalidCallMessage)
			require.NoError(t, err)
			require.Len(t, msgs, 1)
			errWithRecipient := msgs[0]
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
		msgs, err := dealer.ReceiveMessage(callee.ID(), unregister)
		require.NoError(t, err)
		require.Len(t, msgs, 1)
		unregWithRecipient := msgs[0]
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
	messagesWithRecipient, err := dealer.ReceiveMessage(callee.ID(), call)
	require.NoError(t, err)
	require.Equal(t, callee.ID(), messagesWithRecipient[0].Recipient)
	invocation := messagesWithRecipient[0].Message.(*messages.Invocation)
	require.True(t, invocation.Details()[wampproto.OptionReceiveProgress].(bool))

	for i := 0; i < 10; i++ {
		yield := messages.NewYield(invocation.RequestID(), map[string]any{wampproto.OptionProgress: true}, []any{}, nil)
		messagesWithRecipient, err = dealer.ReceiveMessage(callee.ID(), yield)
		require.NoError(t, err)
		require.Equal(t, callee.ID(), messagesWithRecipient[0].Recipient)
		result := messagesWithRecipient[0].Message.(*messages.Result)
		require.Equal(t, call.RequestID(), result.RequestID())
		require.True(t, result.Details()[wampproto.OptionProgress].(bool))
	}

	yield := messages.NewYield(invocation.RequestID(), map[string]any{}, []any{}, nil)
	messagesWithRecipient, err = dealer.ReceiveMessage(callee.ID(), yield)
	require.NoError(t, err)
	require.Equal(t, callee.ID(), messagesWithRecipient[0].Recipient)
	result := messagesWithRecipient[0].Message.(*messages.Result)
	require.Equal(t, call.RequestID(), result.RequestID())
	progress, _ := result.Details()[wampproto.OptionReceiveProgress].(bool)
	require.False(t, progress)
}

func TestProgressiveCallInvocations(t *testing.T) {
	dealer := wampproto.NewDealer()

	callee := wampproto.NewSessionDetails(1, "realm", "authid", "anonymous", false)
	caller := wampproto.NewSessionDetails(2, "realm", "authid", "anonymous", false)

	err := dealer.AddSession(callee)
	require.NoError(t, err)
	err = dealer.AddSession(caller)
	require.NoError(t, err)

	register := messages.NewRegister(3, nil, "foo.bar")
	_, err = dealer.ReceiveMessage(callee.ID(), register)
	require.NoError(t, err)

	call := messages.NewCall(4, map[string]any{wampproto.OptionProgress: true}, "foo.bar", []any{}, nil)
	messageWithRecipient, err := dealer.ReceiveMessage(callee.ID(), call)
	require.NoError(t, err)
	require.Equal(t, callee.ID(), messageWithRecipient[0].Recipient)

	invMessage := messageWithRecipient[0].Message.(*messages.Invocation)
	require.True(t, invMessage.Details()[wampproto.OptionProgress].(bool))

	invRequestID := invMessage.RequestID()
	for i := 0; i < 10; i++ {
		call = messages.NewCall(4, map[string]any{wampproto.OptionProgress: true}, "foo.bar", []any{}, nil)
		messageWithRecipient, err = dealer.ReceiveMessage(callee.ID(), call)
		require.NoError(t, err)

		invMessage = messageWithRecipient[0].Message.(*messages.Invocation)
		require.True(t, invMessage.Details()[wampproto.OptionProgress].(bool))
		require.Equal(t, invRequestID, invMessage.RequestID())
	}

	finalCall := messages.NewCall(4, map[string]any{}, "foo.bar", []any{}, nil)
	messageWithRecipient, err = dealer.ReceiveMessage(callee.ID(), finalCall)
	require.NoError(t, err)
	require.Equal(t, callee.ID(), messageWithRecipient[0].Recipient)

	invocation := messageWithRecipient[0].Message.(*messages.Invocation)
	inProgress, _ := invocation.Details()[wampproto.OptionProgress].(bool)
	require.False(t, inProgress)
}

func TestDealerPrefixRegistration(t *testing.T) {
	testDealerRegistrationAndCall(t,
		wampproto.MatchPrefix,
		"foo.bar.",
		"foo.bar.test",
	)
}

func TestDealerWildcardRegistration(t *testing.T) {
	testDealerRegistrationAndCall(t,
		wampproto.MatchWildcard,
		"foo.bar*test",
		"foo.bar.alpha.test",
	)
}

func testDealerRegistrationAndCall(t *testing.T, matchType, procedure, callURI string) {
	dealer := wampproto.NewDealer()

	callee := wampproto.NewSessionDetails(1, "realm", "authid", "anonymous", false)
	err := dealer.AddSession(callee)
	require.NoError(t, err)

	t.Run("Register", func(t *testing.T) {
		register := messages.NewRegister(1, map[string]any{
			wampproto.OptionMatch: matchType,
		}, procedure)
		msg, err := dealer.ReceiveMessage(callee.ID(), register)
		require.NoError(t, err)
		require.NotNil(t, msg)
		require.Equal(t, msg[0].Recipient, callee.ID())
		require.Equal(t, messages.MessageTypeRegistered, msg[0].Message.Type())
		require.True(t, dealer.HasProcedure(procedure))
	})

	t.Run("Call", func(t *testing.T) {
		caller := wampproto.NewSessionDetails(2, "realm", "authid", "anonymous", false)
		err := dealer.AddSession(caller)
		require.NoError(t, err)

		call := messages.NewCall(3, map[string]any{}, callURI, []any{"abc"}, nil)
		invWithRecipient, err := dealer.ReceiveMessage(caller.ID(), call)
		require.NoError(t, err)
		require.NotNil(t, invWithRecipient)
		require.Equal(t, callee.ID(), invWithRecipient[0].Recipient)
		require.Equal(t, messages.MessageTypeInvocation, invWithRecipient[0].Message.Type())

		// receive yield for invocation
		invocation := invWithRecipient[0].Message.(*messages.Invocation)
		yield := messages.NewYield(invocation.RequestID(), map[string]any{}, []any{"abc"}, nil)
		yieldWithRecipient, err := dealer.ReceiveMessage(caller.ID(), yield)
		require.NoError(t, err)
		require.NotNil(t, yieldWithRecipient)
		require.Equal(t, caller.ID(), yieldWithRecipient[0].Recipient)
		require.Equal(t, messages.MessageTypeResult, yieldWithRecipient[0].Message.Type())
	})
}

func TestDealerCancelMessage(t *testing.T) {
	dealer := wampproto.NewDealer()

	caller := wampproto.NewSessionDetails(1, "realm", "authid", "anonymous", false)
	callee := wampproto.NewSessionDetails(2, "realm", "authid", "anonymous", false)

	require.NoError(t, dealer.AddSession(caller))
	require.NoError(t, dealer.AddSession(callee))

	const procedure = "foo.bar"
	register := messages.NewRegister(1, nil, procedure)
	msgs, err := dealer.ReceiveMessage(callee.ID(), register)
	require.NoError(t, err)
	require.Len(t, msgs, 1)
	require.Equal(t, messages.MessageTypeRegistered, msgs[0].Message.Type())

	callAndCancel := func(requestID int64, cancelMode string) []*wampproto.MessageWithRecipient {
		call := messages.NewCall(requestID, nil, procedure, nil, nil)
		msgs, err = dealer.ReceiveMessage(caller.ID(), call)
		require.NoError(t, err)
		require.Len(t, msgs, 1)
		require.Equal(t, messages.MessageTypeInvocation, msgs[0].Message.Type())

		cancel := messages.NewCancel(requestID, map[string]any{wampproto.OptionMode: cancelMode})
		msgs, err = dealer.ReceiveMessage(caller.ID(), cancel)
		require.NoError(t, err)

		return msgs
	}

	validateErrorMessage := func(msg *wampproto.MessageWithRecipient) {
		require.Equal(t, caller.ID(), msg.Recipient)
		require.Equal(t, messages.MessageTypeError, msg.Message.Type())
		errorMsg := msg.Message.(*messages.Error)
		require.Equal(t, wampproto.ErrCanceled, errorMsg.URI())
	}

	validateInterruptMessage := func(msg *wampproto.MessageWithRecipient) {
		require.Equal(t, callee.ID(), msg.Recipient)
		require.Equal(t, messages.MessageTypeInterrupt, msg.Message.Type())
		interrupt := msg.Message.(*messages.Interrupt)
		require.Equal(t, wampproto.ErrCanceled, interrupt.Options()[wampproto.OptionReason])
	}

	t.Run("CancelModeSkip", func(t *testing.T) {
		msgs = callAndCancel(1, wampproto.CancelModeSkip)
		require.Len(t, msgs, 1)
		validateErrorMessage(msgs[0])
	})

	t.Run("CancelModeKill", func(t *testing.T) {
		msgs = callAndCancel(2, wampproto.CancelModeKill)
		require.Len(t, msgs, 1)
		validateInterruptMessage(msgs[0])
	})

	t.Run("CancelModeKillNoWait", func(t *testing.T) {
		msgs = callAndCancel(3, wampproto.CancelModeKillNoWait)
		require.Len(t, msgs, 2)
		validateInterruptMessage(msgs[0])
		validateErrorMessage(msgs[1])
	})

	t.Run("CancelInvalidInvocation", func(t *testing.T) {
		cancelInvalid := messages.NewCancel(999, nil)
		msgs, err = dealer.ReceiveMessage(caller.ID(), cancelInvalid)
		require.Error(t, err)
		require.Contains(t, err.Error(), "no pending invocation to cancel")
	})
}
