package messages

import "fmt"

const MessageTypeAuthenticate = 5
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

type Authenticate interface {
	Message

	Signature() string
	Extra() map[string]any
}

type authenticate struct {
	signature string
	extra     map[string]any
}

func NewEmptyAuthenticate() Authenticate {
	return &authenticate{}
}

func NewAuthenticate(signature string, extra map[string]any) Authenticate {
	return &authenticate{
		signature: signature,
		extra:     extra,
	}
}

func (a *authenticate) Type() int {
	return MessageTypeAuthenticate
}

func (a *authenticate) Parse(wampMsg []any) error {
	fields, err := ValidateMessage(wampMsg, authenticateValidationSpec)
	if err != nil {
		return fmt.Errorf("authenticate: failed to validate message %s: %w", MessageNameAuthenticate, err)
	}

	a.signature = fields.Signature
	a.extra = fields.Extra

	return nil
}

func (a *authenticate) Marshal() []any {
	return []any{MessageTypeAuthenticate, a.signature, a.extra}
}

func (a *authenticate) Signature() string {
	return a.signature
}

func (a *authenticate) Extra() map[string]any {
	return a.extra
}
