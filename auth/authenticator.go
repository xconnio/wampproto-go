package auth

import "github.com/xconnio/wampproto-go/messages"

type ClientAuthenticator interface {
	AuthMethod() string
	AuthID() string
	AuthExtra() map[string]any
	Authenticate(challenge messages.Challenge) (*messages.Authenticate, error)
}
