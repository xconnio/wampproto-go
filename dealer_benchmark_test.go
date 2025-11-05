package wampproto_test

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/xconnio/wampproto-go"
	"github.com/xconnio/wampproto-go/messages"
)

func BenchmarkDealerConcurrentRegistrations(b *testing.B) {
	dealer := wampproto.NewDealer()
	session := wampproto.NewSessionDetails(1, "realm", "authid", "anonymous", "", false, wampproto.RouterRoles, nil)
	require.NoError(b, dealer.AddSession(session))

	const proc = "io.xconn.test"

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			register := messages.NewRegister(rand.Uint64(), map[string]any{
				wampproto.OptionInvoke: wampproto.InvokeRoundRobin,
			}, proc)

			msg, err := dealer.ReceiveMessage(session.ID(), register)
			require.NoError(b, err)
			require.NotNil(b, msg)
			require.Equal(b, session.ID(), msg.Recipient)
			require.Equal(b, messages.MessageTypeRegistered, msg.Message.Type())
		}
	})
}

func performCall(b *testing.B, dealer *wampproto.Dealer, caller, callee *wampproto.SessionDetails, proc string) {
	call := messages.NewCall(rand.Uint64(), nil, proc, []any{"arg1"}, nil)

	msg, err := dealer.ReceiveMessage(caller.ID(), call)
	require.NoError(b, err)
	require.NotNil(b, msg)
	require.Equal(b, messages.MessageTypeInvocation, msg.Message.Type())

	inv := msg.Message.(*messages.Invocation)
	yield := messages.NewYield(inv.RequestID(), nil, []any{"result"}, nil)
	_, err = dealer.ReceiveMessage(callee.ID(), yield)
	require.NoError(b, err)
}

func BenchmarkDealerConcurrentCalls(b *testing.B) {
	dealer := wampproto.NewDealer()

	callee := wampproto.NewSessionDetails(1, "realm", "authid", "anonymous", "", false, wampproto.RouterRoles, nil)
	caller := wampproto.NewSessionDetails(2, "realm", "authid", "anonymous", "", false, wampproto.RouterRoles, nil)

	require.NoError(b, dealer.AddSession(callee))
	require.NoError(b, dealer.AddSession(caller))

	const procCount = 1000
	for i := 0; i < procCount; i++ {
		proc := "io.xconn.test." + strconv.Itoa(i)
		register := messages.NewRegister(uint64(i+1), nil, proc)
		_, err := dealer.ReceiveMessage(callee.ID(), register)
		require.NoError(b, err)
	}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			proc := "io.xconn.test." + strconv.Itoa(rand.Intn(procCount))
			performCall(b, dealer, caller, callee, proc)
		}
	})
}

func BenchmarkDealerConcurrentPrefixCalls(b *testing.B) {
	dealer := wampproto.NewDealer()

	callee := wampproto.NewSessionDetails(1, "realm", "authid", "anonymous", "", false, wampproto.RouterRoles, nil)
	caller := wampproto.NewSessionDetails(2, "realm", "authid", "anonymous", "", false, wampproto.RouterRoles, nil)

	require.NoError(b, dealer.AddSession(callee))
	require.NoError(b, dealer.AddSession(caller))

	const procCount = 1000
	for i := 0; i < procCount; i++ {
		proc := "io.xconn.test." + strconv.Itoa(i)
		register := messages.NewRegister(uint64(i+1), map[string]any{
			wampproto.OptionMatch: wampproto.MatchPrefix,
		}, proc)
		_, err := dealer.ReceiveMessage(callee.ID(), register)
		require.NoError(b, err)
	}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			prefixNumber := strconv.Itoa(rand.Intn(procCount))
			proc := "io.xconn.test." + prefixNumber + ".sub" + prefixNumber
			performCall(b, dealer, caller, callee, proc)
		}
	})
}

func BenchmarkDealerConcurrentWildcardCalls(b *testing.B) {
	dealer := wampproto.NewDealer()

	callee := wampproto.NewSessionDetails(1, "realm", "authid", "anonymous", "", false, wampproto.RouterRoles, nil)
	caller := wampproto.NewSessionDetails(2, "realm", "authid", "anonymous", "", false, wampproto.RouterRoles, nil)

	require.NoError(b, dealer.AddSession(callee))
	require.NoError(b, dealer.AddSession(caller))

	const procCount = 1000
	for i := 0; i < procCount; i++ {
		proc := "io.xconn.test.*." + strconv.Itoa(i)
		register := messages.NewRegister(uint64(i+1), map[string]any{
			wampproto.OptionMatch: wampproto.MatchWildcard,
		}, proc)
		_, err := dealer.ReceiveMessage(callee.ID(), register)
		require.NoError(b, err)
	}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			proc := "io.xconn.test.sub." + strconv.Itoa(rand.Intn(procCount))
			performCall(b, dealer, caller, callee, proc)
		}
	})
}
