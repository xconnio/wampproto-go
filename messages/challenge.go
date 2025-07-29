package messages

import "fmt"

const MessageTypeChallenge uint64 = 4
const MessageNameChallenge = "CHALLENGE"

var challengeValidationSpec = ValidationSpec{ //nolint:gochecknoglobals
	MinLength: 3,
	MaxLength: 3,
	Message:   MessageNameChallenge,
	Spec: Spec{
		1: ValidateAuthMethod,
		2: ValidateExtra,
	},
}

type ChallengeFields interface {
	AuthMethod() string
	Extra() map[string]any
}

type challengeFields struct {
	authMethod string
	extra      map[string]any
}

func (c *challengeFields) AuthMethod() string {
	return c.authMethod
}

func (c *challengeFields) Extra() map[string]any {
	return c.extra
}

type Challenge struct {
	ChallengeFields
}

func NewChallengeWithFields(fields ChallengeFields) *Challenge {
	return &Challenge{ChallengeFields: fields}
}

func NewChallenge(authMethod string, extra map[string]any) *Challenge {
	return &Challenge{ChallengeFields: &challengeFields{authMethod: authMethod, extra: extra}}
}

func (c *Challenge) Type() uint64 {
	return MessageTypeChallenge
}

func (c *Challenge) Parse(wampMsg []any) error {
	fields, err := ValidateMessage(wampMsg, challengeValidationSpec)
	if err != nil {
		return fmt.Errorf("challenge: failed to validate message %s: %w", MessageNameChallenge, err)
	}

	c.ChallengeFields = &challengeFields{fields.AuthMethod, fields.Extra}

	return nil
}

func (c *Challenge) Marshal() []any {
	return []any{MessageTypeChallenge, c.AuthMethod(), c.Extra()}
}
