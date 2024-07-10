package auth_test

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/xconnio/wampproto-go/auth"
	"github.com/xconnio/wampproto-go/tests"
)

const testPublicKey = "2b7ec216daa877c7f4c9439db8a722ea2340eacad506988db2564e258284f895"
const testPrivateKey = "022b089bed5ab78808365e82dd12c796c835aeb98b4a5a9e099d3e72cb719516"

func TestGenerateCryptosignChallenge(t *testing.T) {
	challenge, err := auth.GenerateCryptoSignChallenge()
	require.NoError(t, err)

	var signChallengeCommand = fmt.Sprintf("auth cryptosign sign-challenge %s %s", challenge, testPrivateKey)
	signature, err := tests.RunCommand(signChallengeCommand)
	require.NoError(t, err)

	var verifySignatureCommand = fmt.Sprintf("auth cryptosign verify-signature %s %s", signature, testPublicKey)
	_, err = tests.RunCommand(verifySignatureCommand)
	require.NoError(t, err)
}

func TestSignCryptosignChallenge(t *testing.T) {
	var challengeCommand = "auth cryptosign generate-challenge"
	challenge, err := tests.RunCommand(challengeCommand)
	require.NoError(t, err)

	privateKey, err := hex.DecodeString(testPrivateKey + testPublicKey)
	require.NoError(t, err)
	signature, err := auth.SignCryptoSignChallenge(challenge, privateKey)
	require.NoError(t, err)

	var verifySignatureCommand = fmt.Sprintf("auth cryptosign verify-signature %s %s", signature, testPublicKey)
	_, err = tests.RunCommand(verifySignatureCommand)
	require.NoError(t, err)
}

func TestVerifyCryptosignSignature(t *testing.T) {
	var challengeCommand = "auth cryptosign generate-challenge"
	challenge, err := tests.RunCommand(challengeCommand)
	require.NoError(t, err)

	var signChallengeCommand = fmt.Sprintf("auth cryptosign sign-challenge %s %s", challenge, testPrivateKey)
	signature, err := tests.RunCommand(signChallengeCommand)
	require.NoError(t, err)

	publicKey, err := hex.DecodeString(testPublicKey)
	require.NoError(t, err)
	isVerified, err := auth.VerifyCryptoSignSignature(signature, publicKey)
	require.NoError(t, err)
	require.True(t, isVerified)
}
