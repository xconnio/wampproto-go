package auth

import (
	"errors"

	"github.com/xconnio/wampproto-go/messages"
)

const MethodAnonymous = "anonymous"

type anonymousAuthenticator struct {
	authID    string
	authExtra map[string]any
}

func NewAnonymousAuthenticator(authID string, authExtra map[string]any) ClientAuthenticator {
	if authExtra == nil {
		authExtra = map[string]any{}
	}

	return &anonymousAuthenticator{
		authID:    authID,
		authExtra: authExtra,
	}
}

func (a *anonymousAuthenticator) AuthMethod() string {
	return MethodAnonymous
}

func (a *anonymousAuthenticator) AuthID() string {
	return a.authID
}

func (a *anonymousAuthenticator) AuthExtra() map[string]any {
	return a.authExtra
}

func (a *anonymousAuthenticator) Authenticate(_ messages.Challenge) (*messages.Authenticate, error) {
	return nil, errors.New("func Authenticate() must not be called for anonymous authentication")
}
