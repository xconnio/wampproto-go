package auth

import (
	"fmt"

	"golang.org/x/exp/slices"
)

func SelectAuthMethod(allowedMethods []Method, proposedMethods []string) (Method, error) {
	var unsupported []string
	var method Method
	for _, proposed := range proposedMethods {
		if slices.Contains(allowedMethods, Method(proposed)) {
			method = Method(proposed)
			break
		}

		unsupported = append(unsupported, proposed)
	}

	if method == "" {
		return "", fmt.Errorf("server does not have %v auth enabled", unsupported)
	}

	return method, nil
}
