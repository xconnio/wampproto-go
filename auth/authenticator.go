package auth

import (
	"time"

	"github.com/xconnio/wampproto-go/messages"
)

type Method string

const (
	Anonymous  Method = "anonymous"
	Ticket     Method = "ticket"
	WAMPCRA    Method = "wampcra"
	CryptoSign Method = "cryptosign"
)

var Methods = []Method{Anonymous, Ticket, WAMPCRA, CryptoSign} // nolint:gochecknoglobals

type ClientAuthenticator interface {
	AuthMethod() string
	AuthID() string
	AuthExtra() map[string]any
	Authenticate(challenge messages.Challenge) (*messages.Authenticate, error)
}

type ServerAuthenticator interface {
	Methods() []Method
	Authenticate(request Request) (Response, error)
}

type Request interface {
	Realm() string
	AuthMethod() Method
	AuthID() string
	AuthExtra() map[string]any
}

func NewRequest(hello *messages.Hello, authMethod Method) Request {
	return &baseRequest{
		authMethod: authMethod,
		realm:      hello.Realm(),
		authID:     hello.AuthID(),
		authExtra:  hello.AuthExtra(),
	}
}

type baseRequest struct {
	authMethod Method
	realm      string
	authID     string
	authExtra  map[string]any
}

func (b *baseRequest) AuthMethod() Method {
	return b.authMethod
}

func (b *baseRequest) Realm() string {
	return b.realm
}

func (b *baseRequest) AuthID() string {
	return b.authID
}

func (b *baseRequest) AuthExtra() map[string]any {
	return b.authExtra
}

type baseResponse struct {
	authID   string
	authRole string

	ttl time.Duration
}

func NewResponse(authID, authRole string, ttl time.Duration) (Response, error) {
	return &baseResponse{
		authID:   authID,
		authRole: authRole,
		ttl:      ttl,
	}, nil
}

func (r *baseResponse) AuthID() string {
	return r.authID
}

func (r *baseResponse) AuthRole() string {
	return r.authRole
}

func (r *baseResponse) TTL() time.Duration {
	return r.ttl
}

type Response interface {
	AuthID() string
	AuthRole() string

	TTL() time.Duration
}

func NewCryptoSignRequest(hello *messages.Hello, publicKey string) Request {
	return &RequestCryptoSign{
		Request:   NewRequest(hello, CryptoSign),
		publicKey: publicKey,
	}
}

type RequestCryptoSign struct {
	Request

	publicKey string
}

func (r *RequestCryptoSign) PublicKey() string {
	return r.publicKey
}

func NewTicketRequest(hello *messages.Hello, ticket string) Request {
	return &TicketRequest{
		Request: NewRequest(hello, Ticket),
		ticket:  ticket,
	}
}

type TicketRequest struct {
	Request

	ticket string
}

func (r *TicketRequest) Ticket() string {
	return r.ticket
}

type CRAResponse struct {
	Response

	secret     string
	salt       string
	iterations int
	keyLen     int
}

func NewCRAResponse(authID, authRole, secret string, ttl time.Duration) Response {
	response, _ := NewResponse(authID, authRole, ttl)
	return &CRAResponse{
		Response: response,
		secret:   secret,
	}
}

func NewCRAResponseSalted(authID, authRole, secret, salt string, iterations, keyLen int, ttl time.Duration) Response {
	response, _ := NewResponse(authID, authRole, ttl)
	return &CRAResponse{
		Response:   response,
		secret:     secret,
		salt:       salt,
		iterations: iterations,
		keyLen:     keyLen,
	}
}

func (r *CRAResponse) Secret() string {
	return r.secret
}

func (r *CRAResponse) Salt() string {
	return r.salt
}

func (r *CRAResponse) Iterations() int {
	return r.iterations
}

func (r *CRAResponse) KeyLen() int {
	return r.keyLen
}
