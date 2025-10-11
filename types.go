package wampproto

import "github.com/xconnio/wampproto-go/messages"

type SessionDetails struct {
	id          uint64
	realm       string
	authID      string
	authRole    string
	routerRoles map[string]any

	staticSerializer bool
}

func NewSessionDetails(id uint64, realm, authID, authRole string, staticSerializer bool,
	routerRoles map[string]any) *SessionDetails {
	return &SessionDetails{
		id:               id,
		realm:            realm,
		authID:           authID,
		authRole:         authRole,
		staticSerializer: staticSerializer,
		routerRoles:      routerRoles,
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

func (s *SessionDetails) AuthRole() string {
	return s.authRole
}

func (s *SessionDetails) StaticSerializer() bool {
	return s.staticSerializer
}

func (s *SessionDetails) RouterRoles() map[string]any {
	return s.routerRoles
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
