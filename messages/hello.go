package messages

import (
	"fmt"
)

const MessageTypeHello = 1
const MessageNameHello = "HELLO"

var helloValidationSpec = ValidationSpec{ //nolint:gochecknoglobals
	MinLength: 2,
	MaxLength: 3,
	Message:   MessageNameHello,
	Spec: Spec{
		1: ValidateRealm,
		2: ValidateDetails,
	},
}

type Hello interface {
	Message

	Realm() string
	AuthID() string
	AuthMethods() []any
	AuthExtra() map[string]any
	Roles() map[string]any
}

type hello struct {
	realm       string
	authID      string
	authMethods []any
	authExtra   map[string]any
	roles       map[string]any
}

func NewEmptyHello() Hello {
	return &hello{}
}

func NewHello(realm, authID string, authExtra, roles map[string]any, authMethods []any) Hello {
	return &hello{
		realm:       realm,
		authID:      authID,
		authMethods: authMethods,
		authExtra:   authExtra,
		roles:       roles,
	}
}

func (h *hello) Realm() string {
	return h.realm
}

func (h *hello) AuthID() string {
	return h.authID
}

func (h *hello) AuthMethods() []any {
	return h.authMethods
}

func (h *hello) AuthExtra() map[string]any {
	return h.authExtra
}

func (h *hello) Roles() map[string]any {
	return h.roles
}

func (h *hello) Type() int {
	return MessageTypeHello
}

func (h *hello) Parse(wampMsg []any) error {
	fields, err := ValidateMessage(wampMsg, helloValidationSpec)
	if err != nil {
		return fmt.Errorf("hello: failed to validate message %s: %w", MessageNameHello, err)
	}

	h.realm = fields.Realm
	h.authID = fields.Details["authid"].(string)
	h.authExtra = fields.Details["authextra"].(map[string]any)
	h.roles = fields.Details["roles"].(map[string]any)
	h.authMethods = fields.Details["authmethods"].([]any)

	return nil
}

func (h *hello) Marshal() []any {
	if h.authExtra == nil {
		h.authExtra = map[string]any{}
	}

	details := map[string]any{
		"authid":      h.authID,
		"authmethods": h.authMethods,
		"authextra":   h.authExtra,
		"roles":       h.roles,
	}

	return []any{MessageTypeHello, h.realm, details}
}
