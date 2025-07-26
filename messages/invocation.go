package messages

import (
	"fmt"
)

const MessageTypeInvocation = 68
const MessageNameInvocation = "INVOCATION"

var invocationValidationSpec = ValidationSpec{ //nolint:gochecknoglobals
	MinLength: 4,
	MaxLength: 6,
	Message:   MessageNameInvocation,
	Spec: Spec{
		1: ValidateRequestID,
		2: ValidateRegistrationID,
		3: ValidateDetails,
		4: ValidateArgs,
		5: ValidateKwArgs,
	},
}

type InvocationFields interface {
	RequestID() uint64
	RegistrationID() uint64
	Details() map[string]any
	Args() []any
	KwArgs() map[string]any

	BinaryPayload
}

type invocationFields struct {
	requestID      uint64
	registrationID uint64
	details        map[string]any
	args           []any
	kwArgs         map[string]any

	serializer uint64
	payload    []byte
}

func (e *invocationFields) RequestID() uint64 {
	return e.requestID
}

func (e *invocationFields) Details() map[string]any {
	return e.details
}

func (e *invocationFields) RegistrationID() uint64 {
	return e.registrationID
}

func (e *invocationFields) Args() []any {
	return e.args
}

func (e *invocationFields) KwArgs() map[string]any {
	return e.kwArgs
}

func (e *invocationFields) PayloadIsBinary() bool {
	return e.serializer != 0
}

func (e *invocationFields) Payload() []byte {
	return e.payload
}

func (e *invocationFields) PayloadSerializer() uint64 {
	return e.serializer
}

type Invocation struct {
	InvocationFields
}

func NewInvocation(requestID, registrationID uint64, details map[string]any, args []any,
	kwArgs map[string]any) *Invocation {

	if details == nil {
		details = make(map[string]any)
	}

	return &Invocation{InvocationFields: &invocationFields{
		requestID:      requestID,
		registrationID: registrationID,
		details:        details,
		args:           args,
		kwArgs:         kwArgs,
	}}
}

func NewInvocationWithFields(fields InvocationFields) *Invocation {
	return &Invocation{InvocationFields: fields}
}

func NewInvocationBinary(requestID, registrationID uint64, details map[string]any, payload []byte,
	serializer uint64) *Invocation {

	if details == nil {
		details = make(map[string]any)
	}

	return &Invocation{InvocationFields: &invocationFields{
		requestID:      requestID,
		registrationID: registrationID,
		details:        details,
		serializer:     serializer,
		payload:        payload,
	}}
}

func (e *Invocation) Type() int {
	return MessageTypeInvocation
}

func (e *Invocation) Parse(wampMsg []any) error {
	fields, err := ValidateMessage(wampMsg, invocationValidationSpec)
	if err != nil {
		return fmt.Errorf("invocationFields: failed to validate message %s: %w", MessageNameInvocation, err)
	}

	e.InvocationFields = &invocationFields{
		requestID:      fields.RequestID,
		registrationID: fields.RegistrationID,
		details:        fields.Details,
		args:           fields.Args,
		kwArgs:         fields.KwArgs,
	}

	return nil
}

func (e *Invocation) Marshal() []any {
	result := []any{MessageTypeInvocation, e.RequestID(), e.RegistrationID(), e.Details()}

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
