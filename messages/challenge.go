package messages

import "fmt"

const MessageTypeChallenge = 4
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

type Challenge interface {
	Message

	AuthMethod() string
	Extra() map[string]any
}

type challenge struct {
	authMethod string
	extra      map[string]any
}

func NewEmptyChallenge() Challenge {
	return &challenge{}
}

func NewChallenge(authMethod string, extra map[string]any) Challenge {
	return &challenge{
		authMethod: authMethod,
		extra:      extra,
	}
}

func (c *challenge) Type() int {
	return MessageTypeChallenge
}

func (c *challenge) Parse(wampMsg []any) error {
	fields, err := ValidateMessage(wampMsg, challengeValidationSpec)
	if err != nil {
		return fmt.Errorf("challenge: failed to validate message %s: %w", MessageNameChallenge, err)
	}

	c.authMethod = fields.AuthMethod
	c.extra = fields.Extra

	return nil
}

func (c *challenge) Marshal() []any {
	return []any{MessageTypeChallenge, c.authMethod, c.extra}
}

func (c *challenge) AuthMethod() string {
	return c.authMethod
}

func (c *challenge) Extra() map[string]any {
	return c.extra
}
