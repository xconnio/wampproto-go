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

type HelloFields interface {
	Realm() string
	AuthID() string
	AuthMethods() []string
	AuthExtra() map[string]any
	Roles() map[string]any
}

type helloFields struct {
	realm       string
	authID      string
	authMethods []string
	authExtra   map[string]any
	roles       map[string]any
}

func (h *helloFields) Realm() string {
	return h.realm
}

func (h *helloFields) AuthID() string {
	return h.authID
}

func (h *helloFields) AuthMethods() []string {
	return h.authMethods
}

func (h *helloFields) AuthExtra() map[string]any {
	return h.authExtra
}

func (h *helloFields) Roles() map[string]any {
	return h.roles
}

type Hello struct {
	HelloFields
}

func NewHelloWithFields(fields HelloFields) *Hello { return &Hello{HelloFields: fields} }

func NewHello(realm, authID string, authExtra, roles map[string]any, authMethods []string) *Hello {
	return &Hello{HelloFields: &helloFields{
		realm:       realm,
		authID:      authID,
		authMethods: authMethods,
		authExtra:   authExtra,
		roles:       roles,
	}}
}

func (h *Hello) Type() int {
	return MessageTypeHello
}

func (h *Hello) Parse(wampMsg []any) error {
	fields, err := ValidateMessage(wampMsg, helloValidationSpec)
	if err != nil {
		return fmt.Errorf("hello: failed to validate message %s: %w", MessageNameHello, err)
	}

	authID, _ := fields.Details["authid"].(string)

	authMethods, ok := fields.Details["authmethods"].([]any)
	if !ok {
		authMethods = []any{}
	}
	authExtra, ok := fields.Details["authextra"].(map[string]any)
	if !ok {
		authExtra = map[string]any{}
	}
	roles, ok := fields.Details["roles"].(map[string]any)
	if !ok {
		roles = map[string]any{}
	}

	authMethodsStr, err := AnysToStrings(authMethods)
	if err != nil {
		return fmt.Errorf("hello: failed to parse auth methods: %w", err)
	}

	h.HelloFields = &helloFields{
		realm:       fields.Realm,
		authID:      authID,
		authMethods: authMethodsStr,
		authExtra:   authExtra,
		roles:       roles,
	}

	return nil
}

func (h *Hello) Marshal() []any {
	authExtra := map[string]any{}
	if h.AuthExtra() != nil {
		authExtra = h.AuthExtra()
	}

	details := map[string]any{
		"authid":      h.AuthID(),
		"authmethods": h.AuthMethods(),
		"authextra":   authExtra,
		"roles":       h.Roles(),
	}

	return []any{MessageTypeHello, h.Realm(), details}
}
