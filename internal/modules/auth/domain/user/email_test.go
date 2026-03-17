package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmail_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		expected bool
	}{
		{
			name:     "Valid_StandardEmail",
			email:    "user@example.com",
			expected: true,
		},
		{
			name:     "Valid_EmailWithPlus",
			email:    "user+tag@example.co.uk",
			expected: true,
		},
		{
			name:     "Valid_EmailWithSubdomain",
			email:    "test@mail.sub.example.com",
			expected: true,
		},
		{
			name:     "Invalid_EmptyString",
			email:    "",
			expected: false,
		},
		{
			name:     "Invalid_NoAtSymbol",
			email:    "userexample.com",
			expected: false,
		},
		{
			name:     "Invalid_NoDomain",
			email:    "user@",
			expected: false,
		},
		{
			name:     "Invalid_NoLocalPart",
			email:    "@example.com",
			expected: false,
		},
		{
			name:     "Invalid_SpaceInEmail",
			email:    "user @example.com",
			expected: false,
		},
		{
			name:     "Invalid_MultipleAtSymbols",
			email:    "user@@example.com",
			expected: false,
		},
		{
			name:     "Valid_NoTLD_AllowedByStdLib",
			email:    "user@example",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Email(tt.email)
			result := e.IsValid()
			assert.Equal(t, tt.expected, result, "Email validation failed for: %s", tt.email)
		})
	}
}
