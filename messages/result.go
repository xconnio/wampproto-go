package messages

import "fmt"

const MessageTypeResult = 50
const MessageNameResult = "RESULT"

var resultValidationSpec = ValidationSpec{ //nolint:gochecknoglobals
	MinLength: 4,
	MaxLength: 6,
	Message:   MessageNameResult,
	Spec: Spec{
		1: ValidateRequestID,
		2: ValidateDetails,
		4: ValidateArgs,
		5: ValidateKwArgs,
	},
}

type Result interface {
	Message

	RequestID() int64
	Details() map[string]any
	Args() []any
	KwArgs() map[string]any
}

type resultMsg struct {
	requestID int64
	details   map[string]any
	args      []any
	kwArgs    map[string]any
}

func NewEmptyResult() Result {
	return &resultMsg{}
}

func NewResult(requestID int64, details map[string]any, args []any, kwArgs map[string]any) Result {
	return &resultMsg{
		requestID: requestID,
		details:   details,
		args:      args,
		kwArgs:    kwArgs,
	}
}

func (e *resultMsg) RequestID() int64 {
	return e.requestID
}

func (e *resultMsg) Details() map[string]any {
	return e.details
}

func (e *resultMsg) Args() []any {
	return e.args
}

func (e *resultMsg) KwArgs() map[string]any {
	return e.kwArgs
}

func (e *resultMsg) Type() int {
	return MessageTypeResult
}

func (e *resultMsg) Parse(wampMsg []any) error {
	fields, err := ValidateMessage(wampMsg, resultValidationSpec)
	if err != nil {
		return fmt.Errorf("result: failed to validate message %s: %w", MessageNameResult, err)
	}

	e.requestID = fields.RequestID
	e.details = fields.Details
	e.args = fields.Args
	e.kwArgs = fields.KwArgs

	return nil
}

func (e *resultMsg) Marshal() []any {
	result := []any{MessageTypeResult, e.requestID, e.details}

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
