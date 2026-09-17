package util

import (
	"testing"
)

func TestValidateEmailInput(t *testing.T) {
	// Define the structural test cases
	tests := []struct {
		name          string
		input         string
		wantEmail     string
		wantDomain    string
		wantErrSecure bool // true if we expect an validation failure error
	}{
		{
			name:          "Valid standard email",
			input:         "user@example.com",
			wantEmail:     "user@example.com",
			wantDomain:    "example.com",
			wantErrSecure: false,
		},
		{
			name:          "Valid email with padding whitespace",
			input:         "  spaces@domain.org  ",
			wantEmail:     "spaces@domain.org",
			wantDomain:    "domain.org",
			wantErrSecure: false,
		},
		{
			name:          "Valid complex subdomains and tags",
			input:         "alex.dev+testing@mail.co.uk",
			wantEmail:     "alex.dev+testing@mail.co.uk",
			wantDomain:    "mail.co.uk",
			wantErrSecure: false,
		},
		{
			name:          "Invalid empty input",
			input:         "   ",
			wantErrSecure: true,
		},
		{
			name:          "Invalid format missing domain",
			input:         "plainaddress",
			wantErrSecure: true,
		},
		{
			name:          "Invalid format missing user",
			input:         "@missinguser.com",
			wantErrSecure: true,
		},
		{
			name:          "Reject RFC formatted display names",
			input:         "Alice <alice@example.com>",
			wantErrSecure: true,
		},
		{
			name:          "Reject local host domains without tld",
			input:         "admin@localhost",
			wantErrSecure: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEmail, gotDomain, err := ValidateEmailInput(tt.input)

			// Check if error status matches our expectations
			if (err != nil) != tt.wantErrSecure {
				t.Fatalf("ValidateAuthEmail() error = %v, wantErrSecure %v", err, tt.wantErrSecure)
			}

			// If we expected an error, stop checking further values for this case
			if tt.wantErrSecure {
				return
			}

			// Verify the sanitized outputs are correct
			if gotEmail != tt.wantEmail {
				t.Errorf("ValidateAuthEmail() gotEmail = %v, want %v", gotEmail, tt.wantEmail)
			}
			if gotDomain != tt.wantDomain {
				t.Errorf("ValidateAuthEmail() gotDomain = %v, want %v", gotDomain, tt.wantDomain)
			}
		})
	}
}
