package wampproto

import (
	"errors"
	"fmt"
	"strings"

	"github.com/xconnio/wampproto-go/auth"
	"github.com/xconnio/wampproto-go/messages"
	"github.com/xconnio/wampproto-go/serializers"
)

var ClientRoles = map[string]any{ //nolint:gochecknoglobals
	"caller": map[string]any{
		"features": map[string]any{
			FeatureProgressiveCallInvocations: true,
		},
	},
	"callee": map[string]any{
		"features": map[string]any{
			FeatureProgressiveCallInvocations: true,
			FeatureProgressiveCallResults:     true,
			FeatureCallCancelling:             true,
		},
	},
	"publisher": map[string]any{
		"features": map[string]any{
			FeaturePublisherExclusion: true,
		},
	},
	"subscriber": map[string]any{
		"features": map[string]any{},
	},
}

type joinerState uint

const (
	joinerStateNone joinerState = iota
	joinerStateHelloSent
	joinerStateAuthenticateSent
	joinerStateJoined
)

type Joiner struct {
	state         joinerState
	realm         string
	authenticator auth.ClientAuthenticator
	serializer    serializers.Serializer

	sessionDetails *SessionDetails
}

func NewJoiner(realm string, serializer serializers.Serializer, authenticator auth.ClientAuthenticator) *Joiner {
	if serializer == nil {
		serializer = &serializers.JSONSerializer{}
	}

	if authenticator == nil {
		authenticator = auth.NewAnonymousAuthenticator("", nil)
	}

	return &Joiner{
		state:         joinerStateNone,
		realm:         realm,
		serializer:    serializer,
		authenticator: authenticator,
	}
}

func (j *Joiner) SendHello() ([]byte, error) {
	hello := messages.NewHello(
		j.realm,
		j.authenticator.AuthID(),
		j.authenticator.AuthExtra(),
		ClientRoles,
		[]string{j.authenticator.AuthMethod()},
	)

	rawBytes, err := j.serializer.Serialize(hello)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize hello: %w", err)
	}

	j.state = joinerStateHelloSent
	return rawBytes, err
}

func (j *Joiner) Receive(data []byte) ([]byte, error) {
	msg, err := j.serializer.Deserialize(data)
	if err != nil {
		return nil, fmt.Errorf("joiner: failed to deserialize: %w", err)
	}

	msg, err = j.ReceiveMessage(msg)
	if err != nil {
		return nil, err
	}

	// when there is error AND there is nothing to send, this
	// implies that the session has been established successfully.
	// The caller may now call SessionDetails() to get details
	// about the session.
	if msg == nil {
		return nil, nil
	}

	toSend, err := j.serializer.Serialize(msg)
	if err != nil {
		return nil, err
	}

	return toSend, nil
}

func (j *Joiner) ReceiveMessage(msg messages.Message) (messages.Message, error) {
	if msg.Type() == messages.MessageTypeWelcome {
		if j.state != joinerStateHelloSent && j.state != joinerStateAuthenticateSent {
			return nil, errors.New("received WELCOME when it was not expected")
		}

		welcome := msg.(*messages.Welcome)
		roles, _ := welcome.Details()["roles"].(map[string]any)
		j.sessionDetails = NewSessionDetails(welcome.SessionID(), j.realm, welcome.Details()["authid"].(string),
			welcome.Details()["authrole"].(string), j.serializer.Static(), roles)
		j.state = joinerStateJoined

		return nil, nil
	} else if msg.Type() == messages.MessageTypeChallenge {
		if j.state != joinerStateHelloSent {
			return nil, errors.New("received CHALLENGE when it was not expected")
		}

		challenge := msg.(*messages.Challenge)
		authenticate, err := j.authenticator.Authenticate(*challenge)
		if err != nil {
			return nil, err
		}

		j.state = joinerStateAuthenticateSent
		return authenticate, nil
	} else if msg.Type() == messages.MessageTypeAbort {
		abort := msg.(*messages.Abort)

		errStr := abort.Reason()
		if len(abort.Args()) > 0 {
			args := make([]string, len(abort.Args()))
			for i, arg := range abort.Args() {
				args[i] = fmt.Sprintf("%v", arg)
			}
			errStr += ": " + strings.Join(args, ", ")
		}
		if len(abort.KwArgs()) > 0 {
			kwargs := make([]string, 0, len(abort.KwArgs()))
			for key, value := range abort.KwArgs() {
				kwargs = append(kwargs, fmt.Sprintf("%s=%v", key, value))
			}
			errStr += ": " + strings.Join(kwargs, ", ")
		}

		return nil, errors.New(errStr)
	} else {
		return nil, errors.New("received unknown message")
	}
}

func (j *Joiner) SessionDetails() (*SessionDetails, error) {
	if j.sessionDetails == nil {
		return nil, fmt.Errorf("session is not setup yet")
	}

	return j.sessionDetails, nil
}
