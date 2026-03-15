package user

import (
	"testing"

	"github.com/devathh/scene-ai/pkg/consts"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUser_New(t *testing.T) {
	tests := []struct {
		name          string
		firstname     string
		lastname      string
		passwordHash  PasswordHash
		wantErr       error
		wantFirstname string
		wantLastname  string
	}{
		{
			name:          "Success_ValidData",
			firstname:     "John",
			lastname:      "Doe",
			passwordHash:  "hashed_password",
			wantErr:       nil,
			wantFirstname: "John",
			wantLastname:  "Doe",
		},
		{
			name:         "Failure_EmptyFirstname",
			firstname:    "",
			lastname:     "Doe",
			passwordHash: "hashed_password",
			wantErr:      consts.ErrInvalidFirstname,
		},
		{
			name:         "Failure_ShortFirstname",
			firstname:    "Jo",
			lastname:     "Doe",
			passwordHash: "hashed_password",
			wantErr:      consts.ErrInvalidFirstname,
		},
		{
			name:         "Failure_LongFirstname",
			firstname:    string(make([]rune, 65)),
			lastname:     "Doe",
			passwordHash: "hashed_password",
			wantErr:      consts.ErrInvalidFirstname,
		},
		{
			name:          "Success_EmptyLastname",
			firstname:     "John",
			lastname:      "",
			passwordHash:  "hashed_password",
			wantErr:       nil,
			wantFirstname: "John",
			wantLastname:  "",
		},
		{
			name:         "Failure_ShortLastname",
			firstname:    "John",
			lastname:     "Do",
			passwordHash: "hashed_password",
			wantErr:      consts.ErrInvalidLastname,
		},
		{
			name:         "Failure_LongLastname",
			firstname:    "John",
			lastname:     string(make([]rune, 65)),
			passwordHash: "hashed_password",
			wantErr:      consts.ErrInvalidLastname,
		},
		{
			name:          "Success_TrimmedNames",
			firstname:     "  John  ",
			lastname:      "  Doe  ",
			passwordHash:  "hashed_password",
			wantErr:       nil,
			wantFirstname: "John",
			wantLastname:  "Doe",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := New(tt.firstname, tt.lastname, tt.passwordHash)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.Nil(t, user)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, user)
			assert.Equal(t, tt.wantFirstname, user.firstname)
			assert.Equal(t, tt.wantLastname, user.lastname)
			assert.Equal(t, tt.passwordHash, user.passwordHash)
			assert.NotEqual(t, uuid.Nil, user.id)
			assert.False(t, user.createdAt.IsZero())
			assert.False(t, user.updatedAt.IsZero())
			assert.True(t, user.createdAt.Equal(user.updatedAt))
		})
	}
}

func TestUser_Accessors(t *testing.T) {
	hash := PasswordHash("hashed_password")
	user, err := New("John", "Doe", hash)
	require.NoError(t, err)

	tests := []struct {
		name     string
		got      any
		expected any
	}{
		{"ID", user.ID(), user.id},
		{"Firstname", user.Firstname(), user.firstname},
		{"Lastname", user.Lastname(), user.lastname},
		{"PasswordHash", user.PasswordHash(), user.passwordHash},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.got)
		})
	}

	t.Run("Timestamps", func(t *testing.T) {
		assert.False(t, user.CreatedAt().IsZero())
		assert.False(t, user.UpdatedAt().IsZero())
	})
}
