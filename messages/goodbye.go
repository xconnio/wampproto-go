package messages

import "fmt"

const MessageTypeGoodbye = 6
const MessageNameGoodbye = "GOODBYE"

var goodByeValidationSpec = ValidationSpec{ //nolint:gochecknoglobals
	MinLength: 3,
	MaxLength: 3,
	Message:   MessageNameGoodbye,
	Spec: Spec{
		1: ValidateDetails,
		2: ValidateReason,
	},
}

type GoodByeFields interface {
	Details() map[string]any
	Reason() string
}

type goodByeFields struct {
	details map[string]any
	reason  string
}

func NewGoodByeFields(reason string, details map[string]any) GoodByeFields {
	return &goodByeFields{
		reason:  reason,
		details: details,
	}
}

func (g *goodByeFields) Reason() string {
	return g.reason
}

func (g *goodByeFields) Details() map[string]any {
	return g.details
}

type GoodBye struct {
	GoodByeFields
}

func NewGoodByeWithFields(fields GoodByeFields) *GoodBye { return &GoodBye{GoodByeFields: fields} }

func NewGoodBye(reason string, details map[string]any) *GoodBye {
	return &GoodBye{GoodByeFields: NewGoodByeFields(reason, details)}
}

func (g *GoodBye) Type() int {
	return MessageTypeGoodbye
}

func (g *GoodBye) Parse(wampMsg []any) error {
	fields, err := ValidateMessage(wampMsg, goodByeValidationSpec)
	if err != nil {
		return fmt.Errorf("goodbye: failed to validate message %s: %w", MessageNameGoodbye, err)
	}

	g.GoodByeFields = NewGoodByeFields(fields.Reason, fields.Details)

	return nil
}

func (g *GoodBye) Marshal() []any {
	return []any{MessageTypeGoodbye, g.Details(), g.Reason()}
}
