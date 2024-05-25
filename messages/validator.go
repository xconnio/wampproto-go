package messages

import (
	"errors"
	"fmt"
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
	RequestID int64
	URI       string
	Args      []interface{}
	KwArgs    map[string]interface{}

	SessionID int64

	Realm       string
	AuthID      string
	AuthRole    string
	AuthMethod  string
	AuthMethods []string
	AuthExtra   map[string]any
	Roles       map[string]any

	MessageType int
	Signature   string
	Reason      string
	Topic       string

	Extra   map[string]any
	Options map[string]any
	Details map[string]any

	SubscriptionID int64
	PublicationID  int64

	RegistrationID int64
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

func ValidateArguments(wampMsg []any, index int, fields *Fields) error {
	data, err := validateSlice(wampMsg, index)
	if err != nil {
		return err
	}

	fields.Args = data
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

func ValidateDetails(wampMsg []any, index int, fields *Fields) error {
	data, err := validateMap(wampMsg, index)
	if err != nil {
		return err
	}

	fields.Details = data
	return nil
}

func ValidateKwArguments(wampMsg []any, index int, fields *Fields) error {
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

func AsInt64(i interface{}) (int64, bool) {
	switch v := i.(type) {
	case int64:
		return v, true
	case uint64:
		return int64(v), true
	case uint8:
		return int64(v), true
	case int:
		return int64(v), true
	case int8:
		return int64(v), true
	case int32:
		return int64(v), true
	case uint:
		return int64(v), true
	case uint32:
		return int64(v), true
	case float64:
		return int64(v), true
	case float32:
		return int64(v), true
	}
	return 0, false
}
