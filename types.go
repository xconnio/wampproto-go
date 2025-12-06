package wampproto

import (
	"github.com/xconnio/wampproto-go/auth"
	"github.com/xconnio/wampproto-go/messages"
)

type SessionDetails struct {
	id          uint64
	realm       string
	authID      string
	authRoles   []string
	routerRoles map[string]any
	createdAt   string
	authExtra   map[string]any
	authMethod  string

	staticSerializer bool
}

func NewSessionDetails(id uint64, realm, authID, authMethod string, authRoles []string, staticSerializer bool,
	routerRoles, authExtra map[string]any) *SessionDetails {
	if routerRoles == nil {
		routerRoles = make(map[string]any)
	}
	if authExtra == nil {
		authExtra = make(map[string]any)
	}
	return &SessionDetails{
		id:               id,
		realm:            realm,
		authID:           authID,
		authRoles:        authRoles,
		staticSerializer: staticSerializer,
		routerRoles:      routerRoles,
		createdAt:        auth.NowISO8601(),
		authMethod:       authMethod,
		authExtra:        authExtra,
	}
}

func (s *SessionDetails) ID() uint64 {
	return s.id
}

func (s *SessionDetails) Realm() string {
	return s.realm
}

func (s *SessionDetails) AuthID() string {
	return s.authID
}

func (s *SessionDetails) AuthRoles() []string {
	return s.authRoles
}

func (s *SessionDetails) AuthMethod() string {
	return s.authMethod
}

func (s *SessionDetails) StaticSerializer() bool {
	return s.staticSerializer
}

func (s *SessionDetails) RouterRoles() map[string]any {
	return s.routerRoles
}

func (s *SessionDetails) CreatedAt() string {
	return s.createdAt
}

func (s *SessionDetails) AuthExtra() map[string]any {
	return s.authExtra
}

type MessageWithRecipient struct {
	Message   messages.Message
	Recipient uint64
}

type Subscription struct {
	ID          uint64
	Topic       string
	Subscribers map[uint64]uint64
	Match       string
}

type Publication struct {
	Event      *messages.Event
	Recipients []uint64
	Ack        *MessageWithRecipient
}
