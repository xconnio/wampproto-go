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
	MessageType() uint64
	RequestID() uint64
	Details() map[string]any
	URI() string
	Args() []any
	KwArgs() map[string]any

	BinaryPayload
}

type errorFields struct {
	messageType uint64
	requestID   uint64
	details     map[string]any
	uri         string
	args        []any
	kwArgs      map[string]any
}

func (e *errorFields) MessageType() uint64 {
	return e.messageType
}

func (e *errorFields) RequestID() uint64 {
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

func (e *errorFields) PayloadIsBinary() bool {
	return false
}

func (e *errorFields) Payload() []byte {
	return nil
}

func (e *errorFields) PayloadSerializer() uint64 {
	return 0
}

type Error struct {
	ErrorFields
}

func NewErrorWithFields(fields ErrorFields) *Error { return &Error{ErrorFields: fields} }

func NewError(messageType, requestID uint64, details map[string]any, uri string, args []any,
	kwArgs map[string]any) *Error {
	if details == nil {
		details = make(map[string]any)
	}

	return &Error{ErrorFields: &errorFields{
		messageType: messageType,
		requestID:   requestID,
		details:     details,
		uri:         uri,
		args:        args,
		kwArgs:      kwArgs,
	}}
}

func (e *Error) Type() int {
	return MessageTypeError
}

func (e *Error) Parse(wampMsg []any) error {
	fields, err := ValidateMessage(wampMsg, errorValidationSpec)
	if err != nil {
		return fmt.Errorf("error: failed to validate message %s: %w", MessageNameError, err)
	}

	e.ErrorFields = &errorFields{
		messageType: fields.MessageType,
		requestID:   fields.RequestID,
		details:     fields.Details,
		uri:         fields.URI,
		args:        fields.Args,
		kwArgs:      fields.KwArgs,
	}

	return nil
}

func (e *Error) Marshal() []any {
	result := []any{MessageTypeError, e.MessageType(), e.RequestID(), e.Details(), e.URI()}

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
