package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"golang.org/x/crypto/pbkdf2"

	"github.com/xconnio/wampproto-go/messages"
)

const MethodCRA = "wampcra"

type craAuthenticator struct {
	authID    string
	authExtra map[string]any

	secret string
}

func NewCRAAuthenticator(authID string, authExtra map[string]any, secret string) ClientAuthenticator {
	return &craAuthenticator{
		authID:    authID,
		authExtra: authExtra,
		secret:    secret,
	}
}

func (a *craAuthenticator) AuthMethod() string {
	return MethodCRA
}

func (a *craAuthenticator) AuthID() string {
	return a.authID
}

func (a *craAuthenticator) AuthExtra() map[string]any {
	return a.authExtra
}

func (a *craAuthenticator) Authenticate(challenge messages.Challenge) (*messages.Authenticate, error) {
	ch, _ := challenge.Extra()["challenge"].(string)
	// If the client needed to look up a user's key, this would require decoding
	// the JSON-encoded challenge string and getting the authid.  For this
	// example assume that client only operates as one user and knows the key
	// to use.

	var rawSecret []byte
	saltStr, _ := challenge.Extra()["salt"].(string)
	// If no salt given, use raw password as key.
	if saltStr != "" {
		// If salting info give, then compute a derived key using PBKDF2.
		iters, _ := messages.AsInt64(challenge.Extra()["iterations"])
		keylen, _ := messages.AsInt64(challenge.Extra()["keylen"])

		rawSecret = DeriveCRAKey(saltStr, a.secret, int(iters), int(keylen))
	} else {
		rawSecret = []byte(a.secret)
	}

	challengeStr := SignCRAChallenge(ch, rawSecret)
	return messages.NewAuthenticate(challengeStr, map[string]any{}), nil
}

// SignCRAChallengeBytes computes the HMAC-SHA256, using the given key, over the
// challenge string, and returns the result.
func SignCRAChallengeBytes(ch string, key []byte) []byte {
	sig := hmac.New(sha256.New, key)
	sig.Write([]byte(ch))
	return sig.Sum(nil)
}

// SignCRAChallenge computes the HMAC-SHA256, using the given key, over the
// challenge string, and returns the result as a base64-encoded string.
func SignCRAChallenge(ch string, key []byte) string {
	return base64.StdEncoding.EncodeToString(SignCRAChallengeBytes(ch, key))
}

func DeriveCRAKey(saltStr string, secret string, iterations int, keyLength int) []byte {
	// If salting info give, then compute a derived key using PBKDF2.
	salt := []byte(saltStr)

	if iterations == 0 {
		iterations = 1000
	}
	if keyLength == 0 {
		keyLength = 32
	}

	// Compute derived key.
	dk := pbkdf2.Key([]byte(secret), salt, iterations, keyLength, sha256.New)
	// Get base64 bytes. see https://github.com/gammazero/nexus/issues/252
	derivedKey := []byte(base64.StdEncoding.EncodeToString(dk))

	return derivedKey
}

// VerifyCRASignature compares a signature to a signature that the computed over
// the given challenge string using the key.  The signature is a base64-encoded
// string, generally presented by a client, and the challenge string and key
// are used to compute the expected HMAC signature.  If these are the same,
// then true is returned.
func VerifyCRASignature(sig, chal string, key []byte) bool {
	sigBytes, err := base64.StdEncoding.DecodeString(sig)
	if err != nil {
		return false
	}

	return hmac.Equal(sigBytes, SignCRAChallengeBytes(chal, key))
}

func GenerateCRAChallenge(session int64, authid, authrole, provider string) (string, error) {
	nonce, err := makeNonce()
	if err != nil {
		return "", fmt.Errorf("failed to get nonce: %w", err)
	}

	data := map[string]any{
		"nonce":        nonce,
		"authprovider": provider,
		"authid":       authid,
		"timestamp":    NowISO8601(),
		"authrole":     authrole,
		"authmethod":   MethodCRA,
		"session":      int(session),
	}

	raw, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(raw), nil
}

// makeNonce generates 16 random bytes as a base64 encoded string.
func makeNonce() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}
