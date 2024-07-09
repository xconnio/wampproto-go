package auth_test

import (
	"crypto/ed25519"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/xconnio/wampproto-go/auth"
	"github.com/xconnio/wampproto-go/messages"
)

const (
	testAuthID     = "foo"
	testPublicKey  = "2b7ec216daa877c7f4c9439db8a722ea2340eacad506988db2564e258284f895"
	testPrivateKey = "022b089bed5ab78808365e82dd12c796c835aeb98b4a5a9e099d3e72cb719516"
	testChallenge  = "5c7b195948a15a9f94bdc31cf6d88294e380985e2c40f5f912fd707e080cb5ff"
)

func TestNewCryptoSignAuthenticator(t *testing.T) {
	authenticator, err := auth.NewCryptoSignAuthenticator(testAuthID, nil, testPrivateKey)
	require.NoError(t, err)

	require.Equal(t, testAuthID, authenticator.AuthID())
	require.Equal(t, auth.MethodCryptoSign, authenticator.AuthMethod())
	require.Equal(t, testPublicKey, authenticator.AuthExtra()["pubkey"])

	t.Run("InvalidPrivateKey", func(t *testing.T) {
		_, err = auth.NewCryptoSignAuthenticator(testAuthID, nil, "invalidkey")
		require.Error(t, err)
	})

	t.Run("MismatchedPublicKey", func(t *testing.T) {
		authExtra := map[string]any{"pubkey": "d057db19aa21f20419dfd385c3dae0cc39ecfc94f8acc921d5f9fc76443098f0"}
		_, err = auth.NewCryptoSignAuthenticator("authID", authExtra, testPrivateKey)
		require.Error(t, err)
	})

	t.Run("Authenticate", func(t *testing.T) {
		challenge := messages.NewChallenge(auth.MethodCryptoSign, map[string]any{"challenge": testChallenge})

		authenticate, err := authenticator.Authenticate(*challenge)
		require.NoError(t, err)

		signed, err := hex.DecodeString(authenticate.Signature())
		require.NoError(t, err)
		require.Equal(t, 96, len(signed))
	})
}

func TestSignCryptoSignChallenge(t *testing.T) {
	privateKey, err := hex.DecodeString(testPrivateKey + testPublicKey)
	require.NoError(t, err)

	signature, err := auth.SignCryptoSignChallenge(testChallenge, privateKey)
	require.NoError(t, err)

	signatureBytes, err := hex.DecodeString(signature)
	require.NoError(t, err)

	require.Equal(t, 96, len(signatureBytes))
}

func TestVerifyCryptoSignSignature(t *testing.T) {
	privateKey, err := hex.DecodeString(testPrivateKey + testPublicKey)
	require.NoError(t, err)

	signature, err := auth.SignCryptoSignChallenge(testChallenge, privateKey)
	require.NoError(t, err)

	publicKeyBytes, err := hex.DecodeString(testPublicKey)
	require.NoError(t, err)

	isVerified, err := auth.VerifyCryptoSignSignature(signature, publicKeyBytes)
	require.NoError(t, err)
	require.True(t, isVerified)
}

func TestGenerateCryptoSignChallenge(t *testing.T) {
	challenge, err := auth.GenerateCryptoSignChallenge()
	require.NoError(t, err)
	require.Equal(t, 64, len(challenge))
}

func TestGenerateCryptoSignKeyPair(t *testing.T) {
	publicKey, privateKey, err := auth.GenerateCryptoSignKeyPair()
	require.NoError(t, err)

	pubKeyBytes, err := hex.DecodeString(publicKey)
	require.NoError(t, err)
	require.Equal(t, ed25519.PublicKeySize, len(pubKeyBytes))

	priKeyBytes, err := hex.DecodeString(privateKey)
	require.NoError(t, err)
	require.Equal(t, ed25519.SeedSize, len(priKeyBytes))
}
