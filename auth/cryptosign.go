package auth

import (
	"crypto/ed25519"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/xconnio/wampproto-go/messages"
)

const MethodCryptoSign = "cryptosign"

type cryptoSignAuthenticator struct {
	authID    string
	authExtra map[string]any

	privateKey ed25519.PrivateKey
}

func NewCryptoSignAuthenticator(authID string, authExtra map[string]any,
	privateKeyHex string) (ClientAuthenticator, error) {

	privateKeyRaw, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return nil, errors.New("invalid private key")
	}

	privateKey := ed25519.NewKeyFromSeed(privateKeyRaw)
	publicKey := privateKey.Public().(ed25519.PublicKey)
	publicKeyHex := hex.EncodeToString(publicKey)

	if authExtra == nil {
		authExtra = map[string]any{"pubkey": publicKeyHex}
	} else if val, ok := authExtra["pubkey"].(string); ok {
		if val != publicKeyHex {
			return nil, errors.New("provided pubkey does not correspond to the private key")
		}
	}

	return &cryptoSignAuthenticator{
		authID:     authID,
		authExtra:  authExtra,
		privateKey: privateKey,
	}, nil
}

func (a *cryptoSignAuthenticator) AuthMethod() string {
	return MethodCryptoSign
}

func (a *cryptoSignAuthenticator) AuthID() string {
	return a.authID
}

func (a *cryptoSignAuthenticator) AuthExtra() map[string]any {
	return a.authExtra
}

func (a *cryptoSignAuthenticator) Authenticate(challenge messages.Challenge) (messages.Authenticate, error) {
	challengeHex, _ := challenge.Extra()["challenge"].(string)
	result, err := SignCryptoSignChallenge(challengeHex, a.privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign challenge")
	}

	return messages.NewAuthenticate(result, map[string]any{}), nil
}

func SignCryptoSignChallenge(challenge string, privateKey ed25519.PrivateKey) (string, error) {
	challengeRaw, err := hex.DecodeString(challenge)
	if err != nil {
		return "", fmt.Errorf("failed to decode challenge: %w", err)
	}

	signedRaw := ed25519.Sign(privateKey, challengeRaw)
	signedHex := hex.EncodeToString(signedRaw)

	return signedHex + challenge, nil
}
