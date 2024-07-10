package auth_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/xconnio/wampproto-go/auth"
	"github.com/xconnio/wampproto-go/tests"
)

const testSecret = "private"

func TestGenerateCRAChallenge(t *testing.T) {
	challenge, err := auth.GenerateCRAChallenge(1, "anonymous", "anonymous", "static")
	require.NoError(t, err)

	var signChallengeCommand = fmt.Sprintf("auth cra sign-challenge %s %s", challenge, testSecret)
	signature, err := tests.RunCommand(signChallengeCommand)
	require.NoError(t, err)

	var verifySignatureCommand = fmt.Sprintf("auth cra verify-signature %s %s %s", challenge, signature, testSecret)
	_, err = tests.RunCommand(verifySignatureCommand)
	require.NoError(t, err)
}

func TestSignCRAChallenge(t *testing.T) {
	var challengeCommand = "auth cra generate-challenge 1 anonymous anonymous static"
	challenge, err := tests.RunCommand(challengeCommand)
	require.NoError(t, err)

	signature := auth.SignCRAChallenge(challenge, []byte(testSecret))
	require.NoError(t, err)

	var verifySignatureCommand = fmt.Sprintf("auth cra verify-signature %s %s %s", challenge, signature, testSecret)
	_, err = tests.RunCommand(verifySignatureCommand)
	require.NoError(t, err)
}

func TestVerifyCRAChallenge(t *testing.T) {
	var challengeCommand = "auth cra generate-challenge 1 anonymous anonymous static"
	challenge, err := tests.RunCommand(challengeCommand)
	require.NoError(t, err)

	var signChallengeCommand = fmt.Sprintf("auth cra sign-challenge %s %s", challenge, testSecret)
	signature, err := tests.RunCommand(signChallengeCommand)
	require.NoError(t, err)

	isVerified := auth.VerifyCRASignature(signature, challenge, []byte(testSecret))
	require.True(t, isVerified)
}
