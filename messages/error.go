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

type ErrorFields interface {
	MessageType() int64
	RequestID() int64
	Details() map[string]any
	URI() string
	Args() []any
	KwArgs() map[string]any
}

type errorFields struct {
	messageType int64
	requestID   int64
	details     map[string]any
	uri         string
	args        []any
	kwArgs      map[string]any
}

func NewErrorFields(messageType, requestID int64, uri string, args []any, kwArgs map[string]any) ErrorFields {
	return &errorFields{
		messageType: messageType,
		requestID:   requestID,
		uri:         uri,
		args:        args,
		kwArgs:      kwArgs,
	}
}

func (e *errorFields) MessageType() int64 {
	return e.messageType
}

func (e *errorFields) RequestID() int64 {
	return e.requestID
}

func (e *errorFields) Details() map[string]any {
	return e.details
}

func (e *errorFields) URI() string {
	return e.uri
}

func (e *errorFields) Args() []any {
	return e.args
}

func (e *errorFields) KwArgs() map[string]any {
	return e.kwArgs
}

type Error struct {
	ErrorFields
}

func NewError(fields ErrorFields) *Error {
	return &Error{
		ErrorFields: fields,
	}
}

func (e *Error) Type() int {
	return MessageTypeError
}

func (e *Error) Parse(wampMsg []any) error {
	fields, err := ValidateMessage(wampMsg, errorValidationSpec)
	if err != nil {
		return fmt.Errorf("error: failed to validate message %s: %w", MessageNameError, err)
	}

	e.ErrorFields = NewErrorFields(fields.MessageType, fields.RequestID, fields.URI, fields.Args, fields.KwArgs)

	return nil
}

func (e *Error) Marshal() []any {
	result := []any{MessageTypeError, e.MessageType(), e.RequestID(), e.URI()}

	if e.Args() != nil {
		result = append(result, e.Args())
	}

	if e.KwArgs() != nil {
		if e.Args() == nil {
			result = append(result, []any{})
		}

		result = append(result, e.KwArgs())
	}

	return result
}
