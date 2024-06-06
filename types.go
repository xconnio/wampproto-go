package wampproto

import "github.com/xconnio/wampproto-go/messages"

type SessionDetails struct {
	id       int64
	realm    string
	authID   string
	authRole string

	staticSerializer bool
}

func NewSessionDetails(id int64, realm, authID, authRole string, staticSerializer bool) *SessionDetails {
	return &SessionDetails{
		id:               id,
		realm:            realm,
		authID:           authID,
		authRole:         authRole,
		staticSerializer: staticSerializer,
	}
}

func (s *SessionDetails) ID() int64 {
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

type MessageWithRecipient struct {
	Message   messages.Message
	Recipient int64
}
