package messages

import (
	"errors"
	"fmt"

	"github.com/xconnio/wampproto-go/util"
)

const errString = "item at index %d must be of type %s but was %T"

type Validator func(wampMsg []any, index int, fields *Fields) error

type Spec map[int]Validator

type ValidationSpec struct {
	MinLength int
	MaxLength int
	Message   string
	Spec      Spec
}

type Fields struct {
	RequestID uint64
	URI       string
	Args      []any
	KwArgs    map[string]any

	SessionID uint64

	Realm       string
	AuthID      string
	AuthRole    string
	AuthMethod  string
	AuthMethods []string
	AuthExtra   map[string]any
	Roles       map[string]any

	MessageType uint64
	Signature   string
	Reason      string

	Extra   map[string]any
	Options map[string]any
	Details map[string]any

	SubscriptionID uint64
	PublicationID  uint64

	RegistrationID uint64
}

func sanityCheck(wampMsg []any, minLength, maxLength int) error {
	if len(wampMsg) < minLength {
		return fmt.Errorf("unexpected message length, must be atleast %d, was %d", minLength, len(wampMsg))
	}

	if len(wampMsg) > maxLength {
		return fmt.Errorf("unexpected message length, must be atmost %d, was %d", maxLength, len(wampMsg))
	}

	return nil
}

func validateID(wampMsg []any, index int) (uint64, error) {
	item, ok := util.AsUInt64(wampMsg[index])
	if !ok {
		return 0, fmt.Errorf(errString, index, "uint64", wampMsg[index])
	}

	return item, nil
}

func validateString(wampMsg []any, index int) (string, error) {
	item, ok := wampMsg[index].(string)
	if !ok {
		return "", fmt.Errorf(errString, index, "string", wampMsg[index])
	}

	return item, nil
}

func validateSlice(wampMsg []any, index int) ([]any, error) {
	item, ok := wampMsg[index].([]any)
	if !ok {
		return nil, fmt.Errorf(errString, index, "[]any", wampMsg[index])
	}

	return item, nil
}

func validateMap(wampMsg []any, index int) (map[string]any, error) {
	item, ok := wampMsg[index].(map[string]any)
	if !ok {
		return nil, fmt.Errorf(errString, index, "map[string]any", wampMsg[index])
	}

	return item, nil
}

func ValidateArgs(wampMsg []any, index int, fields *Fields) error {
	// Ignore if Args item doesn't exist
	if index >= len(wampMsg) {
		return nil
	}

	data, err := validateSlice(wampMsg, index)
	if err != nil {
		return err
	}

	fields.Args = data
	return nil
}

func ValidateSessionID(wampMsg []any, index int, fields *Fields) error {
	data, err := validateID(wampMsg, index)
	if err != nil {
		return err
	}

	fields.SessionID = data
	return nil
}

func ValidateMessageType(wampMsg []any, index int, fields *Fields) error {
	data, err := validateID(wampMsg, index)
	if err != nil {
		return err
	}

	fields.MessageType = data
	return nil
}

func ValidateRequestID(wampMsg []any, index int, fields *Fields) error {
	data, err := validateID(wampMsg, index)
	if err != nil {
		return err
	}

	fields.RequestID = data
	return nil
}

func ValidateRegistrationID(wampMsg []any, index int, fields *Fields) error {
	data, err := validateID(wampMsg, index)
	if err != nil {
		return err
	}

	fields.RegistrationID = data
	return nil
}

func ValidatePublicationID(wampMsg []any, index int, fields *Fields) error {
	data, err := validateID(wampMsg, index)
	if err != nil {
		return err
	}

	fields.PublicationID = data
	return nil
}

func ValidateSubscriptionID(wampMsg []any, index int, fields *Fields) error {
	data, err := validateID(wampMsg, index)
	if err != nil {
		return err
	}

	fields.SubscriptionID = data
	return nil
}

func ValidateSignature(wampMsg []any, index int, fields *Fields) error {
	data, err := validateString(wampMsg, index)
	if err != nil {
		return err
	}

	fields.Signature = data
	return nil
}

func ValidateURI(wampMsg []any, index int, fields *Fields) error {
	data, err := validateString(wampMsg, index)
	if err != nil {
		return err
	}

	fields.URI = data
	return nil
}

func ValidateRealm(wampMsg []any, index int, fields *Fields) error {
	data, err := validateString(wampMsg, index)
	if err != nil {
		return err
	}

	fields.Realm = data
	return nil
}

func ValidateAuthMethod(wampMsg []any, index int, fields *Fields) error {
	data, err := validateString(wampMsg, index)
	if err != nil {
		return err
	}

	fields.AuthMethod = data
	return nil
}

func ValidateReason(wampMsg []any, index int, fields *Fields) error {
	data, err := validateString(wampMsg, index)
	if err != nil {
		return err
	}

	fields.Reason = data
	return nil
}

func ValidateExtra(wampMsg []any, index int, fields *Fields) error {
	data, err := validateMap(wampMsg, index)
	if err != nil {
		return err
	}

	fields.Extra = data
	return nil
}

func ValidateOptions(wampMsg []any, index int, fields *Fields) error {
	data, err := validateMap(wampMsg, index)
	if err != nil {
		return err
	}

	fields.Options = data
	return nil
}

func ValidateDetails(wampMsg []any, index int, fields *Fields) error {
	data, err := validateMap(wampMsg, index)
	if err != nil {
		return err
	}

	fields.Details = data
	return nil
}

func ValidateKwArgs(wampMsg []any, index int, fields *Fields) error {
	// Ignore if KwArgs item doesn't exist
	if index >= len(wampMsg) {
		return nil
	}

	data, err := validateMap(wampMsg, index)
	if err != nil {
		return err
	}

	fields.KwArgs = data
	return nil
}

func ValidateMessage(wampMsg []any, spec ValidationSpec) (*Fields, error) {
	if err := sanityCheck(wampMsg, spec.MinLength, spec.MaxLength); err != nil {
		return nil, err
	}

	f := &Fields{}
	var errs []error
	for index, validator := range spec.Spec {
		err := validator(wampMsg, index, f)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) != 0 {
		return nil, errors.Join(errs...)
	}

	return f, nil
}
