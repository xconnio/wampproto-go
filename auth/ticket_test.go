package auth_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/xconnio/wampproto-go/auth"
	"github.com/xconnio/wampproto-go/messages"
)

func TestNewTicketAuthenticator(t *testing.T) {
	authenticator := auth.NewTicketAuthenticator(testAuthID, nil, "ticket")

	require.Equal(t, testAuthID, authenticator.AuthID())
	require.Equal(t, auth.MethodTicket, authenticator.AuthMethod())
	require.Nil(t, authenticator.AuthExtra())

	t.Run("Authenticate", func(t *testing.T) {
		challenge := messages.NewChallenge(auth.MethodTicket, nil)

		authenticate, err := authenticator.Authenticate(*challenge)
		require.NoError(t, err)
		require.Equal(t, "ticket", authenticate.Signature())
	})
}
