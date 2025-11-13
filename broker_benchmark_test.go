package wampproto_test

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/xconnio/wampproto-go"
	"github.com/xconnio/wampproto-go/messages"
)

func BenchmarkBrokerConcurrentSubscriptions(b *testing.B) {
	broker := wampproto.NewBroker()
	session := wampproto.NewSessionDetails(1, "realm", "authid", "anonymous", "", false, wampproto.RouterRoles, nil)
	require.NoError(b, broker.AddSession(session))

	const topic = "io.xconn.test"

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			subscribe := messages.NewSubscribe(rand.Uint64(), map[string]any{
				wampproto.OptionInvoke: wampproto.InvokeRoundRobin,
			}, topic)

			msg, err := broker.ReceiveMessage(session.ID(), subscribe)
			require.NoError(b, err)
			require.NotNil(b, msg)
			require.Equal(b, session.ID(), msg.Recipient)
			require.Equal(b, messages.MessageTypeSubscribed, msg.Message.Type())
		}
	})
}

func performPublish(b *testing.B, broker *wampproto.Broker, publisher *wampproto.SessionDetails, topic string) {
	publish := messages.NewPublish(rand.Uint64(), map[string]any{"acknowledge": true}, topic, []any{"arg1"}, nil)

	publication, err := broker.ReceivePublish(publisher.ID(), publish)
	require.NoError(b, err)
	require.NotNil(b, publication)

	require.Equal(b, publication.Ack.Recipient, publisher.ID())
	require.Equal(b, publication.Ack.Message.Type(), messages.MessageTypePublished)
	require.NotNil(b, publication.Event)
}

func BenchmarkBrokerConcurrentPublish(b *testing.B) {
	broker := wampproto.NewBroker()

	subscriber := wampproto.NewSessionDetails(1, "realm", "authid", "anonymous", "", false, wampproto.RouterRoles, nil)
	publisher := wampproto.NewSessionDetails(2, "realm", "authid", "anonymous", "", false, wampproto.RouterRoles, nil)

	require.NoError(b, broker.AddSession(subscriber))
	require.NoError(b, broker.AddSession(publisher))

	const subCount = 1000
	for i := 0; i < subCount; i++ {
		topic := "io.xconn.test." + strconv.Itoa(i)
		subscribe := messages.NewSubscribe(uint64(i+1), nil, topic)
		_, err := broker.ReceiveMessage(subscriber.ID(), subscribe)
		require.NoError(b, err)
	}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			topic := "io.xconn.test." + strconv.Itoa(rand.Intn(subCount))
			performPublish(b, broker, publisher, topic)
		}
	})
}

func BenchmarkBrokerConcurrentPrefixPublish(b *testing.B) {
	broker := wampproto.NewBroker()

	subscriber := wampproto.NewSessionDetails(1, "realm", "authid", "anonymous", "", false, wampproto.RouterRoles, nil)
	publisher := wampproto.NewSessionDetails(2, "realm", "authid", "anonymous", "", false, wampproto.RouterRoles, nil)

	require.NoError(b, broker.AddSession(subscriber))
	require.NoError(b, broker.AddSession(publisher))

	const subCount = 1000
	for i := 0; i < subCount; i++ {
		topic := "io.xconn.test." + strconv.Itoa(i)
		subscribe := messages.NewSubscribe(uint64(i+1), map[string]any{
			wampproto.OptionMatch: wampproto.MatchPrefix,
		}, topic)
		_, err := broker.ReceiveMessage(subscriber.ID(), subscribe)
		require.NoError(b, err)
	}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			prefixNumber := strconv.Itoa(rand.Intn(subCount))
			topic := "io.xconn.test." + prefixNumber + ".sub" + prefixNumber
			performPublish(b, broker, publisher, topic)
		}
	})
}

func BenchmarkBrokerConcurrentWildcardPublish(b *testing.B) {
	broker := wampproto.NewBroker()

	subscriber := wampproto.NewSessionDetails(1, "realm", "authid", "anonymous", "", false, wampproto.RouterRoles, nil)
	publisher := wampproto.NewSessionDetails(2, "realm", "authid", "anonymous", "", false, wampproto.RouterRoles, nil)

	require.NoError(b, broker.AddSession(subscriber))
	require.NoError(b, broker.AddSession(publisher))

	const subCount = 1000
	for i := 0; i < subCount; i++ {
		topic := "io.xconn.test.*." + strconv.Itoa(i)
		subscribe := messages.NewSubscribe(uint64(i+1), map[string]any{
			wampproto.OptionMatch: wampproto.MatchWildcard,
		}, topic)
		_, err := broker.ReceiveMessage(subscriber.ID(), subscribe)
		require.NoError(b, err)
	}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			topic := "io.xconn.test.sub." + strconv.Itoa(rand.Intn(subCount))
			performPublish(b, broker, publisher, topic)
		}
	})
}
