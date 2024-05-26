package messages

import "fmt"

const MessageTypeError = 8
const MessageNameError = "ERROR"

var errorValidationSpec = ValidationSpec{ //nolint:gochecknoglobals
	MinLength: 5,
	MaxLength: 7,
	Message:   MessageNameError,
	Spec: Spec{
		1: ValidateMessageType,
		2: ValidateRequestID,
		3: ValidateDetails,
		4: ValidateURI,
		5: ValidateArgs,
		6: ValidateKwArgs,
	},
}

type Error interface {
	Message

	MessageType() int64
	RequestID() int64
	Details() map[string]any
	URI() string
	Args() []any
	KwArgs() map[string]any
}

type err struct {
	messageType int64
	requestID   int64
	details     map[string]any
	uri         string
	args        []any
	kwArgs      map[string]any
}

func NewEmptyError() Error {
	return &err{}
}

func NewError(messageType, requestID int64, uri string, args []any, kwArgs map[string]any) Error {
	return &err{
		messageType: messageType,
		requestID:   requestID,
		uri:         uri,
		args:        args,
		kwArgs:      kwArgs,
	}
}

func (e *err) MessageType() int64 {
	return e.messageType
}

func (e *err) RequestID() int64 {
	return e.requestID
}

func (e *err) Details() map[string]any {
	return e.details
}

func (e *err) URI() string {
	return e.uri
}

func (e *err) Args() []any {
	return e.args
}

func (e *err) KwArgs() map[string]any {
	return e.kwArgs
}

func (e *err) Type() int {
	return MessageTypeError
}

func (e *err) Parse(wampMsg []any) error {
	fields, err := ValidateMessage(wampMsg, errorValidationSpec)
	if err != nil {
		return fmt.Errorf("error: failed to validate message %s: %w", MessageNameError, err)
	}

	e.messageType = fields.MessageType
	e.requestID = fields.RequestID
	e.details = fields.Details
	e.uri = fields.URI
	e.args = fields.Args
	e.kwArgs = fields.KwArgs

	return nil
}

func (e *err) Marshal() []any {
	result := []any{MessageTypeError, e.messageType, e.requestID, e.uri}

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
