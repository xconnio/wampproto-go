package wampproto

type SessionDetails struct {
	id       int64
	realm    string
	authID   string
	authRole string
}

func NewSessionDetails(id int64, realm, authID, authRole string) *SessionDetails {
	return &SessionDetails{
		id:       id,
		realm:    realm,
		authID:   authID,
		authRole: authRole,
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
