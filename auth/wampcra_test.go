package auth_test

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/xconnio/wampproto-go/auth"
	"github.com/xconnio/wampproto-go/messages"
)

const (
	testSecret       = "secret"
	testCRAChallenge = `{
  "authid": "authid",
  "authmethod": "wampcra",
  "authprovider": "provider",
  "authrole": "authrole",
  "nonce": "VJ/iO7bpl5rCiRGJ7IGuQg==",
  "session": 12345,
  "timestamp": "2024-07-09T14:32:29+0500"
}`
)

func TestNewCRAAuthenticator(t *testing.T) {
	authenticator := auth.NewWAMPCRAAuthenticator(testAuthID, testSecret, nil)

	require.Equal(t, testAuthID, authenticator.AuthID())
	require.Equal(t, auth.MethodCRA, authenticator.AuthMethod())
	require.Nil(t, authenticator.AuthExtra())

	t.Run("Authenticate", func(t *testing.T) {
		challengeData := map[string]any{"challenge": testCRAChallenge}
		challenge := messages.NewChallenge(auth.MethodCRA, challengeData)

		authenticate, err := authenticator.Authenticate(*challenge)
		require.NoError(t, err)

		expectedSig := auth.SignWAMPCRAChallenge(testCRAChallenge, []byte(testSecret))
		require.Equal(t, expectedSig, authenticate.Signature())
	})
}

func signCRAChallenge() []byte {
	key := []byte(testSecret)

	sig := hmac.New(sha256.New, key)
	sig.Write([]byte(testCRAChallenge))
	sigBytes := sig.Sum(nil)

	return sigBytes
}

func TestSignCRAChallengeBytes(t *testing.T) {
	expectedBytes := signCRAChallenge()

	signedBytes := auth.SignWAMPCRAChallengeBytes(testCRAChallenge, []byte(testSecret))
	require.Equal(t, expectedBytes, signedBytes)
}

func TestSignCRAChallenge(t *testing.T) {
	expectedBytes := signCRAChallenge()

	expectedSig := base64.StdEncoding.EncodeToString(expectedBytes)
	signed := auth.SignWAMPCRAChallenge(testCRAChallenge, []byte(testSecret))
	require.Equal(t, expectedSig, signed)
}

func TestDeriveCRAKey(t *testing.T) {
	salt := "somesalt"
	secret := testSecret
	iterations := 1000
	keyLength := 32

	derivedKey := auth.DeriveWAMPCRAKey(salt, secret, iterations, keyLength)
	require.NotNil(t, derivedKey)
	require.Equal(t, base64.StdEncoding.EncodedLen(keyLength), len(derivedKey))
}

func TestVerifyCRASignature(t *testing.T) {
	key := []byte(testSecret)
	sig := auth.SignWAMPCRAChallenge(testCRAChallenge, key)

	valid := auth.VerifyWAMPCRASignature(sig, testCRAChallenge, key)
	require.True(t, valid)

	invalid := auth.VerifyWAMPCRASignature("invalidsig", testCRAChallenge, key)
	require.False(t, invalid)
}

func TestGenerateCRAChallenge(t *testing.T) {
	session := uint64(12345)
	authid := "authid"
	authrole := "authrole"
	provider := "provider"

	challenge, err := auth.GenerateWAMPCRAChallenge(session, authid, authrole, provider)
	require.NoError(t, err)

	var data map[string]any
	err = json.Unmarshal([]byte(challenge), &data)
	require.NoError(t, err)

	require.Equal(t, session, uint64(data["session"].(float64)))
	require.Equal(t, authid, data["authid"])
	require.Equal(t, authrole, data["authrole"])
	require.Equal(t, provider, data["authprovider"])
	require.Equal(t, auth.MethodCRA, data["authmethod"])
	require.NotEmpty(t, data["nonce"])
	require.NotEmpty(t, data["timestamp"])
}
