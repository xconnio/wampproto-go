package auth_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/xconnio/wampproto-go/auth"
	"github.com/xconnio/wampproto-go/messages"
)

func TestNewAnonymousAuthenticator(t *testing.T) {
	authenticator := auth.NewAnonymousAuthenticator(testAuthID, nil)

	require.Equal(t, testAuthID, authenticator.AuthID())
	require.Equal(t, auth.MethodAnonymous, authenticator.AuthMethod())
	require.NotNil(t, authenticator.AuthExtra())
	require.Empty(t, authenticator.AuthExtra()) // AuthExtra should be an empty map

	t.Run("Authenticate", func(t *testing.T) {
		challenge := messages.NewChallenge(auth.MethodAnonymous, nil)

		authenticate, err := authenticator.Authenticate(*challenge)
		require.EqualError(t, err, "func Authenticate() must not be called for anonymous authentication")
		require.Nil(t, authenticate)
	})
}
