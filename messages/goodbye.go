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

type GoodBye interface {
	Message

	Details() map[string]any
	Reason() string
}

type goodBye struct {
	details map[string]any
	reason  string
}

func NewEmptyGoodBye() GoodBye {
	return &goodBye{}
}

func NewGoodBye(reason string, details map[string]any) GoodBye {
	return &goodBye{
		reason:  reason,
		details: details,
	}
}

func (g *goodBye) Reason() string {
	return g.reason
}

func (g *goodBye) Details() map[string]any {
	return g.details
}

func (g *goodBye) Type() int {
	return MessageTypeGoodbye
}

func (g *goodBye) Parse(wampMsg []any) error {
	fields, err := ValidateMessage(wampMsg, goodByeValidationSpec)
	if err != nil {
		return fmt.Errorf("goodbye: failed to validate message %s: %w", MessageNameGoodbye, err)
	}

	g.details = fields.Details
	g.reason = fields.Reason

	return nil
}

func (g *goodBye) Marshal() []any {
	return []any{MessageTypeGoodbye, g.details, g.reason}
}
