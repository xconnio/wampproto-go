package wampproto_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/xconnio/wampproto-go"
	"github.com/xconnio/wampproto-go/auth"
	"github.com/xconnio/wampproto-go/messages"
	"github.com/xconnio/wampproto-go/serializers"
)

const (
	realm      = "realm1"
	authID     = "foo"
	ticket     = "fooTicket"
	secret     = "barSecret"
	privateKey = "175604dcce3944595dad640da1676d5e1e1a3950f872f177b1269981140f1c5d"
	publicKey  = "8096cadfd3af87662d4c6589605801c1e2841c4e2cf3d6c30fb187c09c76c5ac"
)

type Authenticator struct {
}

func NewAuthenticator() *Authenticator {
	return &Authenticator{}
}

func (a *Authenticator) Methods() []auth.Method {
	return []auth.Method{auth.MethodAnonymous, auth.MethodTicket, auth.MethodCRA, auth.MethodCryptoSign}
}

func (a *Authenticator) Authenticate(request auth.Request) (auth.Response, error) {
	switch request.AuthMethod() {
	case auth.MethodAnonymous:
		if request.Realm() == realm && request.AuthID() == authID {
			return auth.NewResponse(request.AuthID(), "anonymous", 0)
		}

		return nil, fmt.Errorf("invalid realm")

	case auth.MethodTicket:
		ticketRequest, ok := request.(*auth.TicketRequest)
		if !ok {
			return nil, fmt.Errorf("invalid request")
		}
		if ticketRequest.Realm() == realm && ticketRequest.Ticket() == ticket {
			return auth.NewResponse(ticketRequest.AuthID(), "anonymous", 0)
		}

		return nil, fmt.Errorf("invalid ticket")

	case auth.MethodCRA:
		if request.Realm() == realm && request.AuthID() == authID {
			return auth.NewCRAResponse(request.AuthID(), "anonymous", secret, 0), nil
		}

		return nil, fmt.Errorf("invalid realm")

	case auth.MethodCryptoSign:
		cryptosignRequest, ok := request.(*auth.RequestCryptoSign)
		if !ok {
			return nil, fmt.Errorf("invalid request")
		}

		if cryptosignRequest.Realm() == realm && cryptosignRequest.PublicKey() == publicKey {
			return auth.NewResponse(cryptosignRequest.AuthID(), "anonymous", 0)
		}

		return nil, fmt.Errorf("unknown publickey")

	default:
		return nil, fmt.Errorf("unknown authentication method: %v", request.AuthMethod())
	}
}

func testAnonymousAuth(t *testing.T, serializer serializers.Serializer) {
	var authenticator = NewAuthenticator()
	anonymousAuthenticator := auth.NewAnonymousAuthenticator(authID, map[string]any{})
	joiner := wampproto.NewJoiner(realm, serializer, anonymousAuthenticator)
	acceptor := wampproto.NewAcceptor(serializer, authenticator)

	hello, err := joiner.SendHello()
	require.NoError(t, err)

	// Process and verify the HELLO message
	payload, welcomed, err := acceptor.Receive(hello)
	require.NoError(t, err)
	require.True(t, welcomed)

	welcome, err := serializer.Deserialize(payload)
	require.NoError(t, err)
	require.IsType(t, &messages.Welcome{}, welcome)

	// Ensure no additional messages are received
	data, err := joiner.Receive(payload)
	require.NoError(t, err)
	require.Nil(t, data)

	// Verify session details are available
	sessionDetails, err := joiner.SessionDetails()
	require.NoError(t, err)
	require.NotEmpty(t, sessionDetails)
}

func TestAnonymousAuth(t *testing.T) {
	t.Run("JSONSerializer", func(t *testing.T) {
		jsonSerializer := &serializers.JSONSerializer{}
		testAnonymousAuth(t, jsonSerializer)
	})

	t.Run("CBORSerializer", func(t *testing.T) {
		cborSerializer := &serializers.CBORSerializer{}
		testAnonymousAuth(t, cborSerializer)
	})

	t.Run("MsgPackSerializer", func(t *testing.T) {
		msgPackSerializer := &serializers.MsgPackSerializer{}
		testAnonymousAuth(t, msgPackSerializer)
	})
}

func testAuth(t *testing.T, clientAuthenticator auth.ClientAuthenticator, serializer serializers.Serializer) error {
	var authenticator = NewAuthenticator()
	joiner := wampproto.NewJoiner(realm, serializer, clientAuthenticator)
	acceptor := wampproto.NewAcceptor(serializer, authenticator)

	hello, err := joiner.SendHello()
	require.NoError(t, err)

	// Process and verify the HELLO message
	challengePayload, welcomed, err := acceptor.Receive(hello)
	require.NoError(t, err)
	require.False(t, welcomed)

	challenge, err := serializer.Deserialize(challengePayload)
	require.NoError(t, err)
	_, ok := challenge.(*messages.Challenge)
	if !ok {
		abort, err := serializer.Deserialize(challengePayload)
		require.NoError(t, err)
		abortMsg := abort.(*messages.Abort)
		return fmt.Errorf("%s", abortMsg.Reason())
	}

	authenticate, err := joiner.Receive(challengePayload)
	require.NoError(t, err)

	welcomePayload, welcomed, err := acceptor.Receive(authenticate)
	require.NoError(t, err)

	if !welcomed {
		abort, err := serializer.Deserialize(welcomePayload)
		require.NoError(t, err)
		abortMsg := abort.(*messages.Abort)
		return fmt.Errorf("%s", abortMsg.Reason())
	}

	welcome, err := serializer.Deserialize(welcomePayload)
	require.NoError(t, err)
	require.IsType(t, &messages.Welcome{}, welcome)

	// Ensure no additional messages are received
	data, err := joiner.Receive(welcomePayload)
	require.NoError(t, err)
	require.Nil(t, data)

	// Verify session details are available
	sessionDetails, err := joiner.SessionDetails()
	require.NoError(t, err)
	require.NotEmpty(t, sessionDetails)

	return nil
}

func TestTicketAuth(t *testing.T) {
	t.Run("JSONSerializer", func(t *testing.T) {
		jsonSerializer := &serializers.JSONSerializer{}
		ticketAuthenticator := auth.NewTicketAuthenticator(authID, ticket, map[string]any{})
		err := testAuth(t, ticketAuthenticator, jsonSerializer)
		require.NoError(t, err)
	})

	t.Run("CBORSerializer", func(t *testing.T) {
		cborSerializer := &serializers.CBORSerializer{}
		ticketAuthenticator := auth.NewTicketAuthenticator(authID, ticket, map[string]any{})
		err := testAuth(t, ticketAuthenticator, cborSerializer)
		require.NoError(t, err)
	})

	t.Run("MsgPackSerializer", func(t *testing.T) {
		msgPackSerializer := &serializers.MsgPackSerializer{}
		ticketAuthenticator := auth.NewTicketAuthenticator(authID, ticket, map[string]any{})
		err := testAuth(t, ticketAuthenticator, msgPackSerializer)
		require.NoError(t, err)
	})

	t.Run("InvalidTicket", func(t *testing.T) {
		jsonSerializer := &serializers.JSONSerializer{}
		ticketAuthenticator := auth.NewTicketAuthenticator(authID, "abc", map[string]any{})
		err := testAuth(t, ticketAuthenticator, jsonSerializer)
		require.EqualError(t, err, "wamp.error.authentication_failed")
	})
}

func TestCRAAuth(t *testing.T) {
	t.Run("JSONSerializer", func(t *testing.T) {
		jsonSerializer := &serializers.JSONSerializer{}
		craAuthenticator := auth.NewCRAAuthenticator(authID, secret, map[string]any{})
		err := testAuth(t, craAuthenticator, jsonSerializer)
		require.NoError(t, err)
	})

	t.Run("CBORSerializer", func(t *testing.T) {
		cborSerializer := &serializers.CBORSerializer{}
		craAuthenticator := auth.NewCRAAuthenticator(authID, secret, map[string]any{})
		err := testAuth(t, craAuthenticator, cborSerializer)
		require.NoError(t, err)
	})

	t.Run("MsgPackSerializer", func(t *testing.T) {
		msgPackSerializer := &serializers.MsgPackSerializer{}
		craAuthenticator := auth.NewCRAAuthenticator(authID, secret, map[string]any{})
		err := testAuth(t, craAuthenticator, msgPackSerializer)
		require.NoError(t, err)
	})

	t.Run("InvalidSecret", func(t *testing.T) {
		jsonSerializer := &serializers.JSONSerializer{}
		craAuthenticator := auth.NewCRAAuthenticator(authID, "abc", map[string]any{})
		err := testAuth(t, craAuthenticator, jsonSerializer)
		require.EqualError(t, err, "wamp.error.authentication_failed")
	})

	t.Run("InvalidAuthID", func(t *testing.T) {
		jsonSerializer := &serializers.JSONSerializer{}
		craAuthenticator := auth.NewCRAAuthenticator("abc", secret, map[string]any{})
		err := testAuth(t, craAuthenticator, jsonSerializer)
		require.EqualError(t, err, "wamp.error.authentication_failed")
	})
}

func TestCryptosignAuth(t *testing.T) {
	t.Run("JSONSerializer", func(t *testing.T) {
		jsonSerializer := &serializers.JSONSerializer{}
		cryptosignAuthenticator, err := auth.NewCryptoSignAuthenticator(authID, privateKey, nil)
		require.NoError(t, err)

		err = testAuth(t, cryptosignAuthenticator, jsonSerializer)
		require.NoError(t, err)
	})

	t.Run("CBORSerializer", func(t *testing.T) {
		cborSerializer := &serializers.CBORSerializer{}
		cryptosignAuthenticator, err := auth.NewCryptoSignAuthenticator(authID, privateKey, nil)
		require.NoError(t, err)

		err = testAuth(t, cryptosignAuthenticator, cborSerializer)
		require.NoError(t, err)
	})

	t.Run("MsgPackSerializer", func(t *testing.T) {
		msgPackSerializer := &serializers.MsgPackSerializer{}
		cryptosignAuthenticator, err := auth.NewCryptoSignAuthenticator(authID, privateKey, nil)
		require.NoError(t, err)

		err = testAuth(t, cryptosignAuthenticator, msgPackSerializer)
		require.NoError(t, err)
	})

	t.Run("InvalidKey", func(t *testing.T) {
		jsonSerializer := &serializers.JSONSerializer{}
		cryptosignAuthenticator, err := auth.NewCryptoSignAuthenticator(authID,
			"2e9bef98114241d2226996cf09faf87dad892643a7c5fde186783470bce21df3", nil)
		require.NoError(t, err)

		err = testAuth(t, cryptosignAuthenticator, jsonSerializer)
		require.EqualError(t, err, "wamp.error.authentication_failed")
	})
}

type testAuthenticator struct {
}

func newTestAuthenticator() *testAuthenticator {
	return &testAuthenticator{}
}

func (a *testAuthenticator) Methods() []auth.Method {
	return []auth.Method{auth.MethodAnonymous}
}

func (a *testAuthenticator) Authenticate(request auth.Request) (auth.Response, error) {
	switch request.AuthMethod() {
	case auth.MethodAnonymous:
		if request.Realm() == realm && request.AuthID() == authID {
			return auth.NewResponse(request.AuthID(), "anonymous", 0)
		}

		return nil, fmt.Errorf("invalid realm")

	default:
		return nil, fmt.Errorf("unsupported authentication method: %v", request.AuthMethod())
	}
}

func TestUnsupportedAuthMethod(t *testing.T) {
	var authenticator = newTestAuthenticator()
	ticketAuthenticator := auth.NewTicketAuthenticator(authID, "", map[string]any{})
	serializer := &serializers.JSONSerializer{}
	joiner := wampproto.NewJoiner(realm, serializer, ticketAuthenticator)
	acceptor := wampproto.NewAcceptor(serializer, authenticator)

	hello, err := joiner.SendHello()
	require.NoError(t, err)

	payload, welcomed, err := acceptor.Receive(hello)
	require.NoError(t, err)
	require.False(t, welcomed)

	abort, err := serializer.Deserialize(payload)
	require.NoError(t, err)
	require.IsType(t, &messages.Abort{}, abort)

	// test supported authmethod
	testAnonymousAuth(t, serializer)
}
