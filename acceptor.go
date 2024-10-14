package wampproto

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/xconnio/wampproto-go/auth"
	"github.com/xconnio/wampproto-go/messages"
	"github.com/xconnio/wampproto-go/serializers"
)

type acceptorState uint

const (
	AcceptorStateNone acceptorState = iota
	AcceptorStateChallengeSent
	AcceptorStateWelcomeSent
)

var RouterRoles = map[string]any{ //nolint:gochecknoglobals
	"dealer": map[string]any{
		"features": map[string]any{
			FeatureProgressiveCallInvocations: true,
			FeatureCallCancelling:             true,
		},
	},
	"broker": map[string]any{
		"features": map[string]any{},
	},
}

type defaultAuthenticator struct{}

func (d *defaultAuthenticator) Methods() []auth.Method {
	return []auth.Method{auth.Anonymous}
}

func (d *defaultAuthenticator) Authenticate(request auth.Request) (auth.Response, error) {
	if request.AuthMethod() != auth.Anonymous {
		return nil, fmt.Errorf("recevied request for %s auth but only anonymous is supported", request.AuthMethod())
	}

	return auth.NewResponse("anonymous", "anonymous", 0)
}

type Acceptor struct {
	state         acceptorState
	serializer    serializers.Serializer
	authenticator auth.ServerAuthenticator
	// cached items
	authMethod auth.Method
	hello      *messages.Hello
	request    auth.Request
	response   auth.Response
	challenge  string

	sessionDetails *SessionDetails
}

func NewAcceptor(serializer serializers.Serializer, authenticator auth.ServerAuthenticator) *Acceptor {
	if authenticator == nil {
		authenticator = &defaultAuthenticator{}
	}

	if serializer == nil {
		serializer = &serializers.JSONSerializer{}
	}

	return &Acceptor{
		serializer:    serializer,
		authenticator: authenticator,
		state:         AcceptorStateNone,
	}
}

func (a *Acceptor) Receive(data []byte) (payload []byte, welcomed bool, err error) {
	msg, err := a.serializer.Deserialize(data)
	if err != nil {
		return nil, false, err
	}

	toSend, err := a.ReceiveMessage(msg)
	if err != nil {
		return nil, false, err
	}

	payload, err = a.serializer.Serialize(toSend)
	if err != nil {
		return nil, false, err
	}

	return payload, toSend.Type() == messages.MessageTypeWelcome, nil
}

func (a *Acceptor) ReceiveMessage(msg messages.Message) (messages.Message, error) {
	if a.state == AcceptorStateWelcomeSent {
		return nil, errors.New("session was established, not expecting any new messages")
	}

	if msg.Type() == messages.MessageTypeHello {
		if a.state != AcceptorStateNone {
			return nil, fmt.Errorf("state must be %d when processing HELLO but was %d", AcceptorStateNone, a.state)
		}

		hello := msg.(*messages.Hello)
		authMethod, err := auth.SelectAuthMethod(a.authenticator.Methods(), hello.AuthMethods())
		if err != nil {
			return nil, err
		}

		a.authMethod = authMethod
		a.hello = hello

		switch authMethod {
		case auth.Anonymous:
			request := auth.NewRequest(hello, authMethod)
			response, err := a.authenticator.Authenticate(request)
			if err != nil {
				abort := messages.NewAbort(map[string]any{}, "wamp.error.authentication_failed", []any{err.Error()}, nil)
				return abort, nil
			}

			return a.sendWelcome(GenerateID(), response), nil
		case auth.Ticket:
			a.state = AcceptorStateChallengeSent
			return messages.NewChallenge(string(authMethod), map[string]any{}), nil
		case auth.WAMPCRA:
			request := auth.NewRequest(hello, authMethod)
			response, err := a.authenticator.Authenticate(request)
			if err != nil {
				abort := messages.NewAbort(map[string]any{}, "wamp.error.authentication_failed", []any{err.Error()}, nil)
				return abort, nil
			}

			craResponse, ok := response.(*auth.CRAResponse)
			if !ok {
				return nil, errors.New("internal response for WAMPCRA auth was of invalid type")
			}

			chStr, err := auth.GenerateCRAChallenge(GenerateID(), response.AuthID(), response.AuthRole(), "dynamic")
			if err != nil {
				return nil, err
			}

			extra := map[string]any{
				"challenge": chStr,
			}

			if craResponse.Salt() != "" {
				extra["salt"] = craResponse.Salt()
				extra["keylen"] = craResponse.KeyLen()
				extra["iterations"] = craResponse.Iterations()
			}

			a.challenge = chStr
			a.response = response
			a.state = AcceptorStateChallengeSent

			return messages.NewChallenge(string(authMethod), extra), nil
		case auth.CryptoSign:
			pKey, ok := hello.AuthExtra()["pubkey"]
			if !ok {
				return nil, fmt.Errorf("pubkey missing in authextra")
			}

			publicKey, ok := pKey.(string)
			if !ok {
				return nil, fmt.Errorf("pubkey must be a string in authextra, was %T", publicKey)
			}

			if publicKey == "" {
				return nil, fmt.Errorf("pubkey empty in authextra")
			}

			request := auth.NewCryptoSignRequest(hello, publicKey)
			response, err := a.authenticator.Authenticate(request)
			if err != nil {
				abort := messages.NewAbort(map[string]any{}, "wamp.error.authentication_failed", []any{err.Error()}, nil)
				return abort, nil
			}

			a.request = request
			a.response = response

			chStr, err := auth.GenerateCryptoSignChallenge()
			if err != nil {
				return nil, err
			}

			a.challenge = chStr
			a.state = AcceptorStateChallengeSent

			return messages.NewChallenge(string(authMethod), map[string]any{"challenge": chStr}), nil
		default:
			return nil, fmt.Errorf("received HELLO for unexpected authmethod %s", authMethod)
		}
	} else if msg.Type() == messages.MessageTypeAuthenticate {
		if a.state != AcceptorStateChallengeSent {
			return nil, fmt.Errorf("received AUTHENTICATE while state was %d", a.state)
		}

		switch a.authMethod {
		case auth.Ticket:
			authenticate := msg.(*messages.Authenticate)
			request := auth.NewTicketRequest(a.hello, authenticate.Signature())
			response, err := a.authenticator.Authenticate(request)
			if err != nil {
				abort := messages.NewAbort(map[string]any{}, "wamp.error.authentication_failed", []any{err.Error()}, nil)
				return abort, nil
			}

			return a.sendWelcome(GenerateID(), response), nil
		case auth.WAMPCRA:
			authenticate := msg.(*messages.Authenticate)
			response := a.response.(*auth.CRAResponse)

			var secret []byte
			if response.Salt() == "" {
				secret = []byte(response.Secret())
			} else {
				secret = auth.DeriveCRAKey(response.Salt(), response.Secret(), response.Iterations(), response.KeyLen())
			}

			if !auth.VerifyCRASignature(authenticate.Signature(), a.challenge, secret) {
				abort := messages.NewAbort(map[string]any{}, "wamp.error.authentication_failed", nil, nil)
				return abort, nil
			}

			return a.sendWelcome(GenerateID(), a.response), nil
		case auth.CryptoSign:
			authenticate := msg.(*messages.Authenticate)
			request := a.request.(*auth.RequestCryptoSign)
			key, err := hex.DecodeString(request.PublicKey())
			if err != nil {
				return nil, fmt.Errorf("failed to decode public key")
			}

			verified, _ := auth.VerifyCryptoSignSignature(authenticate.Signature(), key)
			if !verified {
				abort := messages.NewAbort(map[string]any{}, "wamp.error.authentication_failed", nil, nil)
				return abort, nil
			}

			return a.sendWelcome(GenerateID(), a.response), nil
		default:
			return nil, fmt.Errorf("received AUTHENTICATE for unexpected authmethod %s", a.authMethod)
		}
	} else {
		return nil, fmt.Errorf("received unexpected message %T", msg)
	}
}

func (a *Acceptor) sendWelcome(sessionID int64, response auth.Response) *messages.Welcome {
	welcome := messages.NewWelcome(sessionID, map[string]any{
		"realm":      a.hello.Realm(),
		"roles":      RouterRoles,
		"authid":     response.AuthID(),
		"authrole":   response.AuthRole(),
		"authmethod": a.authMethod,
	})

	a.sessionDetails = NewSessionDetails(sessionID, a.hello.Realm(), response.AuthID(), response.AuthRole(),
		a.serializer.Static())
	a.state = AcceptorStateWelcomeSent

	return welcome
}

func (a *Acceptor) SessionDetails() (*SessionDetails, error) {
	if a.sessionDetails == nil {
		return nil, fmt.Errorf("session is not setup yet")
	}

	return a.sessionDetails, nil
}
