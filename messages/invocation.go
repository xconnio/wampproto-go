package messages

import "fmt"

const MessageTypeInvocation = 68
const MessageNameInvocation = "INVOCATION"

var invocationValidationSpec = ValidationSpec{ //nolint:gochecknoglobals
	MinLength: 4,
	MaxLength: 6,
	Message:   MessageNameInvocation,
	Spec: Spec{
		1: ValidateRequestID,
		2: ValidateDetails,
		3: ValidateURI,
		4: ValidateArgs,
		5: ValidateKwArgs,
	},
}

type Invocation interface {
	Message

	RequestID() int64
	Details() map[string]any
	Procedure() string
	Args() []any
	KwArgs() map[string]any
}

type invocation struct {
	requestID int64
	details   map[string]any
	procedure string
	args      []any
	kwArgs    map[string]any
}

func NewEmptyInvocation() Invocation {
	return &invocation{}
}

func NewInvocation(requestID int64, details map[string]any, procedure string, args []any,
	kwArgs map[string]any) Invocation {
	return &invocation{
		requestID: requestID,
		details:   details,
		procedure: procedure,
		args:      args,
		kwArgs:    kwArgs,
	}
}

func (e *invocation) RequestID() int64 {
	return e.requestID
}

func (e *invocation) Details() map[string]any {
	return e.details
}

func (e *invocation) Procedure() string {
	return e.procedure
}

func (e *invocation) Args() []any {
	return e.args
}

func (e *invocation) KwArgs() map[string]any {
	return e.kwArgs
}

func (e *invocation) Type() int {
	return MessageTypeInvocation
}

func (e *invocation) Parse(wampMsg []any) error {
	fields, err := ValidateMessage(wampMsg, invocationValidationSpec)
	if err != nil {
		return fmt.Errorf("invocation: failed to validate message %s: %w", MessageNameInvocation, err)
	}

	e.requestID = fields.RequestID
	e.details = fields.Details
	e.procedure = fields.URI
	e.args = fields.Args
	e.kwArgs = fields.KwArgs

	return nil
}

func (e *invocation) Marshal() []any {
	result := []any{MessageTypeInvocation, e.requestID, e.details, e.procedure}

	if e.args != nil {
		result = append(result, e.args)
	}

	if e.kwArgs != nil {
		if e.args == nil {
			result = append(result, []any{})
		}

		result = append(result, e.kwArgs)
	}

	return result
}
