package wampproto_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/xconnio/wampproto-go"
	"github.com/xconnio/wampproto-go/messages"
)

const (
	testArg    = "abc"
	testAuthID = "authid"
	testRole   = "anonymous"
	testURI    = "foo.bar"
	testInvoke = "invoke"
)

func TestDealerAddRemoveSession(t *testing.T) {
	dealer := wampproto.NewDealer()

	t.Run("RemoveNonSession", func(t *testing.T) {
		err := dealer.RemoveSession(1)
		require.Error(t, err)
	})

	t.Run("AddRemove", func(t *testing.T) {
		details := wampproto.NewSessionDetails(1, "realm", testAuthID, testRole, "", false, wampproto.RouterRoles, nil)
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

	callee := wampproto.NewSessionDetails(1, "realm", testAuthID, testRole, "", false, wampproto.RouterRoles, nil)
	err := dealer.AddSession(callee)
	require.NoError(t, err)

	var registerationID uint64

	t.Run("Register", func(t *testing.T) {
		register := messages.NewRegister(1, nil, testURI)
		msg, err := dealer.ReceiveMessage(callee.ID(), register)
		require.NoError(t, err)
		require.NotNil(t, msg)
		require.Equal(t, msg.Recipient, callee.ID())
		require.Equal(t, messages.MessageTypeRegistered, msg.Message.Type())

		hasProcedure := dealer.HasProcedure(testURI)
		require.True(t, hasProcedure)
		registerationID = msg.Message.(*messages.Registered).RegistrationID()

		t.Run("DuplicateProcedure", func(t *testing.T) {
			register = messages.NewRegister(2, nil, testURI)
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
			invalidRegister := messages.NewRegister(2, nil, testURI)
			_, err = dealer.ReceiveMessage(5, invalidRegister)
			require.EqualError(t, err, "cannot register procedure for non-existent session 5")
		})
	})

	t.Run("Call", func(t *testing.T) {
		caller := wampproto.NewSessionDetails(2, "realm", testAuthID, testRole, "", false, wampproto.RouterRoles, nil)
		err := dealer.AddSession(caller)
		require.NoError(t, err)

		call := messages.NewCall(3, map[string]any{}, testURI, []any{testArg}, nil)
		invWithRecipient, err := dealer.ReceiveMessage(caller.ID(), call)
		require.NoError(t, err)
		require.NotNil(t, invWithRecipient)
		require.Equal(t, callee.ID(), invWithRecipient.Recipient)
		require.Equal(t, messages.MessageTypeInvocation, invWithRecipient.Message.Type())

		// receive yield for invocation
		invocation := invWithRecipient.Message.(*messages.Invocation)
		yield := messages.NewYield(invocation.RequestID(), map[string]any{}, []any{testArg}, nil)
		yieldWithRecipient, err := dealer.ReceiveMessage(caller.ID(), yield)
		require.NoError(t, err)
		require.NotNil(t, yieldWithRecipient)
		require.Equal(t, caller.ID(), yieldWithRecipient.Recipient)
		require.Equal(t, messages.MessageTypeResult, yieldWithRecipient.Message.Type())

		t.Run("NonExistingProcedure", func(t *testing.T) {
			invalidCallMessage := messages.NewCall(3, map[string]any{}, "invalid", []any{testArg}, nil)
			errWithRecipient, err := dealer.ReceiveMessage(caller.ID(), invalidCallMessage)
			require.NoError(t, err)
			require.NotNil(t, errWithRecipient)
			require.Equal(t, caller.ID(), errWithRecipient.Recipient)
			require.Equal(t, errWithRecipient.Message.Type(), messages.MessageTypeError)
		})

		t.Run("InvalidYield", func(t *testing.T) {
			_, err = dealer.ReceiveMessage(5, yield)
			require.Error(t, err)
		})
	})

	t.Run("Unregister", func(t *testing.T) {
		unregister := messages.NewUnregister(callee.ID(), registerationID)
		unregWithRecipient, err := dealer.ReceiveMessage(callee.ID(), unregister)
		require.NoError(t, err)
		require.NotNil(t, unregWithRecipient)
		require.Equal(t, callee.ID(), unregWithRecipient.Recipient)
		require.Equal(t, messages.MessageTypeUnregistered, unregWithRecipient.Message.Type())

		hasProcedure := dealer.HasProcedure(testURI)
		require.False(t, hasProcedure)

		t.Run("InvalidRegistration", func(t *testing.T) {
			_, err = dealer.ReceiveMessage(callee.ID(), unregister)
			require.EqualError(t, err, "unregister: session 1 has no registration 1")
		})
	})
}

func TestProgressiveCallResults(t *testing.T) {
	dealer := wampproto.NewDealer()

	callee := wampproto.NewSessionDetails(1, "realm", testAuthID, testRole, "", false, wampproto.RouterRoles, nil)
	caller := wampproto.NewSessionDetails(2, "realm", testAuthID, testRole, "", false, wampproto.RouterRoles, nil)

	err := dealer.AddSession(callee)
	require.NoError(t, err)
	err = dealer.AddSession(caller)
	require.NoError(t, err)

	register := messages.NewRegister(1, nil, testURI)
	_, err = dealer.ReceiveMessage(callee.ID(), register)
	require.NoError(t, err)

	call := messages.NewCall(caller.ID(), map[string]any{wampproto.OptionReceiveProgress: true}, testURI, []any{}, nil)
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

func TestProgressiveCallInvocations(t *testing.T) {
	dealer := wampproto.NewDealer()

	callee := wampproto.NewSessionDetails(1, "realm", testAuthID, testRole, "", false, wampproto.RouterRoles, nil)
	caller := wampproto.NewSessionDetails(2, "realm", testAuthID, testRole, "", false, wampproto.RouterRoles, nil)

	err := dealer.AddSession(callee)
	require.NoError(t, err)
	err = dealer.AddSession(caller)
	require.NoError(t, err)

	register := messages.NewRegister(3, nil, testURI)
	_, err = dealer.ReceiveMessage(callee.ID(), register)
	require.NoError(t, err)

	call := messages.NewCall(4, map[string]any{wampproto.OptionProgress: true}, testURI, []any{}, nil)
	messageWithRecipient, err := dealer.ReceiveMessage(callee.ID(), call)
	require.NoError(t, err)
	require.Equal(t, callee.ID(), messageWithRecipient.Recipient)

	invMessage := messageWithRecipient.Message.(*messages.Invocation)
	require.True(t, invMessage.Details()[wampproto.OptionProgress].(bool))

	invRequestID := invMessage.RequestID()
	for i := 0; i < 10; i++ {
		call = messages.NewCall(4, map[string]any{wampproto.OptionProgress: true}, testURI, []any{}, nil)
		messageWithRecipient, err = dealer.ReceiveMessage(callee.ID(), call)
		require.NoError(t, err)

		invMessage = messageWithRecipient.Message.(*messages.Invocation)
		require.True(t, invMessage.Details()[wampproto.OptionProgress].(bool))
		require.Equal(t, invRequestID, invMessage.RequestID())
	}

	finalCall := messages.NewCall(4, map[string]any{}, testURI, []any{}, nil)
	messageWithRecipient, err = dealer.ReceiveMessage(callee.ID(), finalCall)
	require.NoError(t, err)
	require.Equal(t, callee.ID(), messageWithRecipient.Recipient)

	invocation := messageWithRecipient.Message.(*messages.Invocation)
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

	callee := wampproto.NewSessionDetails(1, "realm", testAuthID, testRole, "", false, wampproto.RouterRoles, nil)
	err := dealer.AddSession(callee)
	require.NoError(t, err)

	t.Run("Register", func(t *testing.T) {
		register := messages.NewRegister(1, map[string]any{
			wampproto.OptionMatch: matchType,
		}, procedure)
		msg, err := dealer.ReceiveMessage(callee.ID(), register)
		require.NoError(t, err)
		require.NotNil(t, msg)
		require.Equal(t, msg.Recipient, callee.ID())
		require.Equal(t, messages.MessageTypeRegistered, msg.Message.Type())
		require.True(t, dealer.HasProcedure(procedure))
	})

	t.Run("Call", func(t *testing.T) {
		caller := wampproto.NewSessionDetails(2, "realm", testAuthID, testRole, "", false, wampproto.RouterRoles, nil)
		err := dealer.AddSession(caller)
		require.NoError(t, err)

		call := messages.NewCall(3, map[string]any{}, callURI, []any{testArg}, nil)
		invWithRecipient, err := dealer.ReceiveMessage(caller.ID(), call)
		require.NoError(t, err)
		require.NotNil(t, invWithRecipient)
		require.Equal(t, callee.ID(), invWithRecipient.Recipient)
		require.Equal(t, messages.MessageTypeInvocation, invWithRecipient.Message.Type())

		// receive yield for invocation
		invocation := invWithRecipient.Message.(*messages.Invocation)
		yield := messages.NewYield(invocation.RequestID(), map[string]any{}, []any{testArg}, nil)
		yieldWithRecipient, err := dealer.ReceiveMessage(caller.ID(), yield)
		require.NoError(t, err)
		require.NotNil(t, yieldWithRecipient)
		require.Equal(t, caller.ID(), yieldWithRecipient.Recipient)
		require.Equal(t, messages.MessageTypeResult, yieldWithRecipient.Message.Type())
	})
}

func newCancelTestDealer(t *testing.T) (*wampproto.Dealer, *wampproto.SessionDetails, *wampproto.SessionDetails) {
	t.Helper()
	d := wampproto.NewDealer()
	caller := wampproto.NewSessionDetails(1, "realm", "authid", "anonymous", "", false, wampproto.RouterRoles, nil)
	callee := wampproto.NewSessionDetails(2, "realm", "authid", "anonymous", "", false, wampproto.RouterRoles, nil)
	require.NoError(t, d.AddSession(caller))
	require.NoError(t, d.AddSession(callee))
	reg := messages.NewRegister(1, nil, "foo.bar")
	msg, err := d.ReceiveMessage(callee.ID(), reg)
	require.NoError(t, err)
	require.Equal(t, messages.MessageTypeRegistered, msg.Message.Type())
	return d, caller, callee
}

func TestDealerCancelMessage(t *testing.T) {
	t.Run("CancelModeSkip", func(t *testing.T) {
		dealer, caller, _ := newCancelTestDealer(t)

		call := messages.NewCall(1, nil, "foo.bar", nil, nil)
		_, err := dealer.ReceiveMessage(caller.ID(), call)
		require.NoError(t, err)

		cancel := messages.NewCancel(1, map[string]any{wampproto.OptionMode: wampproto.CancelModeSkip})
		msgs, err := dealer.ReceiveCancel(caller.ID(), cancel)
		require.NoError(t, err)
		require.Len(t, msgs, 1)
		require.Equal(t, caller.ID(), msgs[0].Recipient)
		require.Equal(t, messages.MessageTypeError, msgs[0].Message.Type())
		require.Equal(t, wampproto.ErrCanceled, msgs[0].Message.(*messages.Error).URI())
	})

	t.Run("CancelModeKillNoWait", func(t *testing.T) {
		dealer, caller, callee := newCancelTestDealer(t)

		call := messages.NewCall(1, nil, "foo.bar", nil, nil)
		_, err := dealer.ReceiveMessage(caller.ID(), call)
		require.NoError(t, err)

		cancel := messages.NewCancel(1, map[string]any{wampproto.OptionMode: wampproto.CancelModeKillNoWait})
		msgs, err := dealer.ReceiveCancel(caller.ID(), cancel)
		require.NoError(t, err)
		require.Len(t, msgs, 2)

		require.Equal(t, callee.ID(), msgs[0].Recipient)
		require.Equal(t, messages.MessageTypeInterrupt, msgs[0].Message.Type())
		require.Equal(t, wampproto.ErrCanceled, msgs[0].Message.(*messages.Interrupt).Options()[wampproto.OptionReason])

		require.Equal(t, caller.ID(), msgs[1].Recipient)
		require.Equal(t, messages.MessageTypeError, msgs[1].Message.Type())
		require.Equal(t, wampproto.ErrCanceled, msgs[1].Message.(*messages.Error).URI())
	})

	t.Run("CancelModeKillCalleeRespondsWithError", func(t *testing.T) {
		dealer, caller, callee := newCancelTestDealer(t)

		call := messages.NewCall(1, nil, "foo.bar", nil, nil)
		inv, err := dealer.ReceiveMessage(caller.ID(), call)
		require.NoError(t, err)
		invocationID := inv.Message.(*messages.Invocation).RequestID()

		cancel := messages.NewCancel(1, map[string]any{wampproto.OptionMode: wampproto.CancelModeKill})
		msgs, err := dealer.ReceiveCancel(caller.ID(), cancel)
		require.NoError(t, err)
		require.Len(t, msgs, 1)
		require.Equal(t, callee.ID(), msgs[0].Recipient)
		require.Equal(t, messages.MessageTypeInterrupt, msgs[0].Message.Type())

		// callee responds with ERROR after receiving INTERRUPT
		calleeErr := messages.NewError(messages.MessageTypeInvocation, invocationID, nil, wampproto.ErrCanceled, nil, nil)
		result, err := dealer.ReceiveMessage(callee.ID(), calleeErr)
		require.NoError(t, err)
		require.Equal(t, caller.ID(), result.Recipient)
		require.Equal(t, messages.MessageTypeError, result.Message.Type())
		require.Equal(t, wampproto.ErrCanceled, result.Message.(*messages.Error).URI())
	})

	t.Run("CancelModeKillCalleeIgnoresInterruptAndYields", func(t *testing.T) {
		dealer, caller, callee := newCancelTestDealer(t)

		call := messages.NewCall(1, nil, "foo.bar", nil, nil)
		inv, err := dealer.ReceiveMessage(caller.ID(), call)
		require.NoError(t, err)
		invocationID := inv.Message.(*messages.Invocation).RequestID()

		cancel := messages.NewCancel(1, map[string]any{wampproto.OptionMode: wampproto.CancelModeKill})
		_, err = dealer.ReceiveCancel(caller.ID(), cancel)
		require.NoError(t, err)

		// callee ignores INTERRUPT and sends YIELD anyway
		yield := messages.NewYield(invocationID, nil, []any{"result"}, nil)
		result, err := dealer.ReceiveMessage(callee.ID(), yield)
		require.NoError(t, err)
		require.Equal(t, caller.ID(), result.Recipient)
		require.Equal(t, messages.MessageTypeError, result.Message.Type())
		require.Equal(t, wampproto.ErrCanceled, result.Message.(*messages.Error).URI())
	})

	t.Run("CancelInvalidMode", func(t *testing.T) {
		dealer, caller, _ := newCancelTestDealer(t)

		call := messages.NewCall(1, nil, "foo.bar", nil, nil)
		_, err := dealer.ReceiveMessage(caller.ID(), call)
		require.NoError(t, err)

		cancel := messages.NewCancel(1, map[string]any{wampproto.OptionMode: "bogus"})
		msgs, err := dealer.ReceiveCancel(caller.ID(), cancel)
		require.NoError(t, err)
		require.Len(t, msgs, 1)
		require.Equal(t, caller.ID(), msgs[0].Recipient)
		require.Equal(t, messages.MessageTypeError, msgs[0].Message.Type())
		require.Equal(t, wampproto.ErrInvalidArgument, msgs[0].Message.(*messages.Error).URI())
	})

	t.Run("CancelUnknownRequest", func(t *testing.T) {
		dealer, caller, _ := newCancelTestDealer(t)

		cancel := messages.NewCancel(999, nil)
		msgs, err := dealer.ReceiveCancel(caller.ID(), cancel)
		require.NoError(t, err)
		require.Len(t, msgs, 1)
		require.Equal(t, caller.ID(), msgs[0].Recipient)
		require.Equal(t, messages.MessageTypeError, msgs[0].Message.Type())
		require.Equal(t, wampproto.ErrInvalidArgument, msgs[0].Message.(*messages.Error).URI())
	})
}

func TestDealerDiscloseCallerDetails(t *testing.T) {
	dealer := wampproto.NewDealer()

	callee := wampproto.NewSessionDetails(1, "realm", testAuthID, testRole, "", false, wampproto.RouterRoles, nil)
	err := dealer.AddSession(callee)
	require.NoError(t, err)

	register := messages.NewRegister(1, nil, testURI)
	_, err = dealer.ReceiveMessage(callee.ID(), register)
	require.NoError(t, err)

	caller := wampproto.NewSessionDetails(2, "realm", testAuthID, testRole, "", false, wampproto.RouterRoles, nil)
	err = dealer.AddSession(caller)
	require.NoError(t, err)

	t.Run("DisabledByDefault", func(t *testing.T) {
		call := messages.NewCall(3, map[string]any{}, testURI, []any{testArg}, nil)
		invWithRecipient, err := dealer.ReceiveMessage(caller.ID(), call)
		require.NoError(t, err)
		invocation := invWithRecipient.Message.(*messages.Invocation)
		require.Equal(t, map[string]any{}, invocation.Details())
	})

	t.Run("Enable", func(t *testing.T) {
		dealer.AutoDiscloseCaller(true)
		call := messages.NewCall(4, map[string]any{}, testURI, []any{testArg}, nil)
		invWithRecipient, err := dealer.ReceiveMessage(caller.ID(), call)
		require.NoError(t, err)
		invocation := invWithRecipient.Message.(*messages.Invocation)
		expectedDetails := map[string]any{"caller": uint64(2), "caller_authid": testAuthID,
			"caller_authrole": testRole, "procedure": testURI}
		require.Equal(t, expectedDetails, invocation.Details())
	})

	t.Run("Disable", func(t *testing.T) {
		dealer.AutoDiscloseCaller(false)
		call := messages.NewCall(4, map[string]any{}, testURI, []any{testArg}, nil)
		invWithRecipient, err := dealer.ReceiveMessage(caller.ID(), call)
		require.NoError(t, err)
		invocation := invWithRecipient.Message.(*messages.Invocation)
		require.Equal(t, map[string]any{}, invocation.Details())
	})
}

func TestDealerInvocationOptions(t *testing.T) {
	dealer := wampproto.NewDealer()

	callee1 := wampproto.NewSessionDetails(1, "realm", testAuthID, testRole, "", false, wampproto.RouterRoles, nil)
	callee2 := wampproto.NewSessionDetails(2, "realm", testAuthID, testRole, "", false, wampproto.RouterRoles, nil)
	require.NoError(t, dealer.AddSession(callee1))
	require.NoError(t, dealer.AddSession(callee2))

	caller := wampproto.NewSessionDetails(3, "realm", testAuthID, testRole, "", false, wampproto.RouterRoles, nil)
	require.NoError(t, dealer.AddSession(caller))

	registerProcedures := func(proc, policy string) {
		for _, callee := range []uint64{callee1.ID(), callee2.ID()} {
			register := messages.NewRegister(callee, map[string]any{testInvoke: policy}, proc)
			msgWithRecipient, err := dealer.ReceiveMessage(callee, register)
			require.NoError(t, err)
			require.Equal(t, messages.MessageTypeRegistered, msgWithRecipient.Message.Type())
		}
	}

	t.Run("First", func(t *testing.T) {
		registerProcedures("first.proc", "first")

		for i := 0; i < 3; i++ {
			call := messages.NewCall(uint64(20+i), nil, "first.proc", nil, nil)
			inv, err := dealer.ReceiveMessage(caller.ID(), call)
			require.NoError(t, err)
			require.Equal(t, callee1.ID(), inv.Recipient)
		}
	})

	t.Run("Last", func(t *testing.T) {
		registerProcedures("last.proc", "last")

		for i := 0; i < 3; i++ {
			call := messages.NewCall(uint64(30+i), nil, "last.proc", nil, nil)
			inv, err := dealer.ReceiveMessage(caller.ID(), call)
			require.NoError(t, err)
			require.Equal(t, callee2.ID(), inv.Recipient)
		}
	})

	t.Run("RegisterFirstAndThenLast", func(t *testing.T) {
		register := messages.NewRegister(callee1.ID(), map[string]any{testInvoke: "first"}, "io.xconn.test")
		_, err := dealer.ReceiveMessage(callee1.ID(), register)
		require.NoError(t, err)

		register1 := messages.NewRegister(callee2.ID(), map[string]any{testInvoke: "last"}, "io.xconn.test")
		msgWithRecipient, err := dealer.ReceiveMessage(callee2.ID(), register1)
		require.NoError(t, err)
		require.Equal(t, messages.MessageTypeError, msgWithRecipient.Message.Type())
	})

	t.Run("RoundRobin", func(t *testing.T) {
		registerProcedures("roundrobin.proc", "roundrobin")

		expectedRecipients := []uint64{callee1.ID(), callee2.ID(), callee1.ID(), callee2.ID()}
		for i, expected := range expectedRecipients {
			call := messages.NewCall(uint64(10+i), nil, "roundrobin.proc", nil, nil)
			inv, err := dealer.ReceiveMessage(caller.ID(), call)
			require.NoError(t, err)
			require.NotNil(t, inv)
			require.Equal(t, expected, inv.Recipient)
		}
	})

	t.Run("Random", func(t *testing.T) {
		registerProcedures("random.proc", "random")

		recipients := map[uint64]bool{}
		for i := 0; i < 10; i++ {
			call := messages.NewCall(uint64(40+i), nil, "random.proc", nil, nil)
			inv, err := dealer.ReceiveMessage(caller.ID(), call)
			require.NoError(t, err)
			recipients[inv.Recipient] = true
		}
		require.Len(t, recipients, 2)
	})

	t.Run("DisconnectOneCalleeAndCall", func(t *testing.T) {
		require.NoError(t, dealer.RemoveSession(callee1.ID()))

		for i := 0; i < 3; i++ {
			call := messages.NewCall(uint64(10+i), nil, "roundrobin.proc", nil, nil)
			inv, err := dealer.ReceiveMessage(caller.ID(), call)
			require.NoError(t, err)
			require.Equal(t, callee2.ID(), inv.Recipient)
		}

		for i := 0; i < 3; i++ {
			call := messages.NewCall(uint64(20+i), nil, "first.proc", nil, nil)
			inv, err := dealer.ReceiveMessage(caller.ID(), call)
			require.NoError(t, err)
			require.Equal(t, callee2.ID(), inv.Recipient)
		}

		for i := 0; i < 3; i++ {
			call := messages.NewCall(uint64(30+i), nil, "last.proc", nil, nil)
			inv, err := dealer.ReceiveMessage(caller.ID(), call)
			require.NoError(t, err)
			require.Equal(t, callee2.ID(), inv.Recipient)
		}

		for i := 0; i < 3; i++ {
			call := messages.NewCall(uint64(40+i), nil, "random.proc", nil, nil)
			inv, err := dealer.ReceiveMessage(caller.ID(), call)
			require.NoError(t, err)
			require.Equal(t, callee2.ID(), inv.Recipient)
		}
	})
}
