package messages

import (
	"fmt"
)

const MessageTypeAuthenticate uint64 = 5
const MessageNameAuthenticate = "AUTHENTICATE"

var authenticateValidationSpec = ValidationSpec{ //nolint:gochecknoglobals
	MinLength: 2,
	MaxLength: 3,
	Message:   MessageNameAuthenticate,
	Spec: Spec{
		1: ValidateSignature,
		2: ValidateExtra,
	},
}

type AuthenticateFields interface {
	Signature() string
	Extra() map[string]any
}

type authenticateFields struct {
	signature string
	extra     map[string]any
}

func (a *authenticateFields) Signature() string {
	return a.signature
}

func (a *authenticateFields) Extra() map[string]any {
	return a.extra
}

type Authenticate struct {
	AuthenticateFields
}

func NewAuthenticateWithFields(fields AuthenticateFields) *Authenticate {
	return &Authenticate{AuthenticateFields: fields}
}

func NewAuthenticate(signature string, extra map[string]any) *Authenticate {
	return &Authenticate{AuthenticateFields: &authenticateFields{signature: signature, extra: extra}}
}

func (a *Authenticate) Type() uint64 {
	return MessageTypeAuthenticate
}

func (a *Authenticate) Parse(wampMsg []any) error {
	fields, err := ValidateMessage(wampMsg, authenticateValidationSpec)
	if err != nil {
		return fmt.Errorf("authenticate: failed to validate message %s: %w", MessageNameAuthenticate, err)
	}

	a.AuthenticateFields = &authenticateFields{
		signature: fields.Signature,
		extra:     fields.Extra,
	}

	return nil
}

func (a *Authenticate) Marshal() []any {
	authExtra := a.Extra()
	if authExtra == nil {
		authExtra = make(map[string]any)
	}

	return []any{MessageTypeAuthenticate, a.Signature(), authExtra}
}
