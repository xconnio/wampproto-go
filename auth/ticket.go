package auth

import (
	"github.com/xconnio/wampproto-go/messages"
)

const MethodTicket = "ticket"

type ticketAuthenticator struct {
	authID    string
	authExtra map[string]any

	ticket string
}

func NewTicketAuthenticator(authID string, authExtra map[string]any, ticket string) ClientAuthenticator {
	return &ticketAuthenticator{
		authID:    authID,
		authExtra: authExtra,
		ticket:    ticket,
	}
}

func (a *ticketAuthenticator) AuthMethod() string {
	return MethodTicket
}

func (a *ticketAuthenticator) AuthID() string {
	return a.authID
}

func (a *ticketAuthenticator) AuthExtra() map[string]any {
	return a.authExtra
}

func (a *ticketAuthenticator) Authenticate(_ messages.Challenge) (*messages.Authenticate, error) {
	return messages.NewAuthenticate(a.ticket, map[string]any{}), nil
}
