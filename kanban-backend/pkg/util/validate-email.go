package util

import (
	"fmt"
	"net/mail"
	"strings"
)

/* Prepares and strictly validates an email address for login/registration.
 * Returns the sanitized email, the extracted domain name, and an error if invalid.
 */
func ValidateEmailInput(input string) (string, string, error) {
	// remove leading/trailing whitespaces
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return "", "", fmt.Errorf("email cannot be empty")
	}

	// parse using the native net/mail package
	addr, err := mail.ParseAddress(trimmed)
	if err != nil {
		return "", "", fmt.Errorf("invalid email format: %w", err)
	}

	// strict match enforcement: reject names like "Alice <alice@example.com>"
	if addr.Address != trimmed {
		return "", "", fmt.Errorf("email must be a raw email address without display names")
	}

	// extract domain and check for basic modern validity (must contain a dot)
	parts := strings.Split(addr.Address, "@")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("malformed email structure")
	}

	domain := parts[1]
	if !strings.Contains(domain, ".") {
		return "", "", fmt.Errorf("invalid email domain extension")
	}

	return addr.Address, domain, nil
}