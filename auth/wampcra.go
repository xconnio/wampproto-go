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
	"github.com/xconnio/wampproto-go/util"
)

const MethodCRA = "wampcra"

type wampcraAuthenticator struct {
	authID    string
	authExtra map[string]any

	secret string
}

func NewWAMPCRAAuthenticator(authID string, secret string, authExtra map[string]any) ClientAuthenticator {
	return &wampcraAuthenticator{
		authID:    authID,
		authExtra: authExtra,
		secret:    secret,
	}
}

func (a *wampcraAuthenticator) AuthMethod() string {
	return MethodCRA
}

func (a *wampcraAuthenticator) AuthID() string {
	return a.authID
}

func (a *wampcraAuthenticator) AuthExtra() map[string]any {
	return a.authExtra
}

func (a *wampcraAuthenticator) Authenticate(challenge messages.Challenge) (*messages.Authenticate, error) {
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
		iters, _ := util.AsInt(challenge.Extra()["iterations"])
		keylen, _ := util.AsInt(challenge.Extra()["keylen"])

		rawSecret = DeriveWAMPCRAKey(saltStr, a.secret, iters, keylen)
	} else {
		rawSecret = []byte(a.secret)
	}

	challengeStr := SignWAMPCRAChallenge(ch, rawSecret)
	return messages.NewAuthenticate(challengeStr, map[string]any{}), nil
}

// SignWAMPCRAChallengeBytes computes the HMAC-SHA256, using the given key, over the
// challenge string, and returns the result.
func SignWAMPCRAChallengeBytes(ch string, key []byte) []byte {
	sig := hmac.New(sha256.New, key)
	sig.Write([]byte(ch))
	return sig.Sum(nil)
}

// SignWAMPCRAChallenge computes the HMAC-SHA256, using the given key, over the
// challenge string, and returns the result as a base64-encoded string.
func SignWAMPCRAChallenge(ch string, key []byte) string {
	return base64.StdEncoding.EncodeToString(SignWAMPCRAChallengeBytes(ch, key))
}

func DeriveWAMPCRAKey(saltStr string, secret string, iterations int, keyLength int) []byte {
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

// VerifyWAMPCRASignature compares a signature to a signature that the computed over
// the given challenge string using the key.  The signature is a base64-encoded
// string, generally presented by a client, and the challenge string and key
// are used to compute the expected HMAC signature.  If these are the same,
// then true is returned.
func VerifyWAMPCRASignature(sig, chal string, key []byte) bool {
	sigBytes, err := base64.StdEncoding.DecodeString(sig)
	if err != nil {
		return false
	}

	return hmac.Equal(sigBytes, SignWAMPCRAChallengeBytes(chal, key))
}

func GenerateWAMPCRAChallenge(session uint64, authid, authrole, provider string) (string, error) {
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
		"session":      session,
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
