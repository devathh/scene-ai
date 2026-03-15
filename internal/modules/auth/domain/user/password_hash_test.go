package user

import (
	"testing"

	"github.com/devathh/scene-ai/pkg/consts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPasswordHash(t *testing.T) {
	tests := []struct {
		name     string
		raw      string
		wantErr  error
		wantEmpty bool
	}{
		{
			name:    "Success_ValidPassword",
			raw:     "securepassword123",
			wantErr: nil,
		},
		{
			name:    "Failure_EmptyPassword",
			raw:     "",
			wantErr: consts.ErrInvalidPassword,
		},
		{
			name:    "Failure_ShortPassword",
			raw:     "short",
			wantErr: consts.ErrInvalidPassword,
		},
		{
			name:    "Failure_MinusOneLengthPassword",
			raw:     "12345",
			wantErr: consts.ErrInvalidPassword,
		},
		{
			name:    "Success_ExactlyMinLengthPassword",
			raw:     "123456",
			wantErr: nil,
		},
		{
			name:    "Success_TrimmedPassword",
			raw:     "  securepassword123  ",
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := NewPasswordHash(tt.raw)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.Empty(t, hash)
				return
			}

			require.NoError(t, err)
			require.NotEmpty(t, hash)
		})
	}
}

func TestPasswordHash_Compare(t *testing.T) {
	raw := "testpassword123"
	hash, err := NewPasswordHash(raw)
	require.NoError(t, err)

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Success_MatchingPassword",
			input:    raw,
			expected: true,
		},
		{
			name:     "Failure_WrongPassword",
			input:    "wrongpassword",
			expected: false,
		},
		{
			name:     "Failure_EmptyInput",
			input:    "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hash.Compare(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPasswordHash_String(t *testing.T) {
	raw := "testpassword123"
	hash, err := NewPasswordHash(raw)
	require.NoError(t, err)

	t.Run("NotRawAndNotEmpty", func(t *testing.T) {
		assert.NotEqual(t, raw, hash.String())
		assert.NotEmpty(t, hash.String())
	})
}
