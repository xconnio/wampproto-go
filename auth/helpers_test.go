package auth_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/xconnio/wampproto-go/auth"
)

func TestSelectAuthMethod(t *testing.T) {
	// case where proposed methods include allowed methods
	allowedMethods := []auth.Method{auth.Ticket, auth.Anonymous}
	proposedMethods := []string{string(auth.Ticket), testChallengeKey, string(auth.Anonymous)}

	method, err := auth.SelectAuthMethod(allowedMethods, proposedMethods)
	require.NoError(t, err)
	require.Equal(t, auth.Ticket, method)

	// case where no proposed methods match allowed methods
	allowedMethods = []auth.Method{auth.Ticket, auth.Anonymous}
	proposedMethods = []string{testChallengeKey}

	method, err = auth.SelectAuthMethod(allowedMethods, proposedMethods)
	require.EqualError(t, err, "server does not have [challenge] auth enabled")
	require.Equal(t, auth.Method(""), method)

	// case where proposed methods is empty
	allowedMethods = []auth.Method{auth.Ticket, auth.Anonymous}
	proposedMethods = []string{}

	method, err = auth.SelectAuthMethod(allowedMethods, proposedMethods)
	require.NoError(t, err)
	require.Equal(t, auth.Anonymous, method)
}
