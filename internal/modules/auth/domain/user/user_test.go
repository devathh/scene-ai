package user

import (
	"testing"
	"time"

	"github.com/devathh/scene-ai/pkg/consts"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUser_New(t *testing.T) {
	validEmail := Email("user@example.com")
	invalidEmail := Email("invalid-email")

	tests := []struct {
		name          string
		email         Email
		firstname     string
		lastname      string
		passwordHash  PasswordHash
		wantErr       error
		wantFirstname string
		wantLastname  string
		wantEmail     Email
	}{
		{
			name:          "Success_ValidData",
			email:         validEmail,
			firstname:     "John",
			lastname:      "Doe",
			passwordHash:  "hashed_password",
			wantErr:       nil,
			wantFirstname: "John",
			wantLastname:  "Doe",
			wantEmail:     validEmail,
		},
		{
			name:         "Failure_InvalidEmail",
			email:        invalidEmail,
			firstname:    "John",
			lastname:     "Doe",
			passwordHash: "hashed_password",
			wantErr:      consts.ErrInvalidEmail,
		},
		{
			name:         "Failure_EmptyFirstname",
			email:        validEmail,
			firstname:    "",
			lastname:     "Doe",
			passwordHash: "hashed_password",
			wantErr:      consts.ErrInvalidFirstname,
		},
		{
			name:         "Failure_ShortFirstname",
			email:        validEmail,
			firstname:    "Jo",
			lastname:     "Doe",
			passwordHash: "hashed_password",
			wantErr:      consts.ErrInvalidFirstname,
		},
		{
			name:         "Failure_LongFirstname",
			email:        validEmail,
			firstname:    string(make([]rune, 65)),
			lastname:     "Doe",
			passwordHash: "hashed_password",
			wantErr:      consts.ErrInvalidFirstname,
		},
		{
			name:          "Success_EmptyLastname",
			email:         validEmail,
			firstname:     "John",
			lastname:      "",
			passwordHash:  "hashed_password",
			wantErr:       nil,
			wantFirstname: "John",
			wantLastname:  "",
			wantEmail:     validEmail,
		},
		{
			name:         "Failure_ShortLastname",
			email:        validEmail,
			firstname:    "John",
			lastname:     "Do",
			passwordHash: "hashed_password",
			wantErr:      consts.ErrInvalidLastname,
		},
		{
			name:         "Failure_LongLastname",
			email:        validEmail,
			firstname:    "John",
			lastname:     string(make([]rune, 65)),
			passwordHash: "hashed_password",
			wantErr:      consts.ErrInvalidLastname,
		},
		{
			name:          "Success_TrimmedNames",
			email:         validEmail,
			firstname:     "  John  ",
			lastname:      "  Doe  ",
			passwordHash:  "hashed_password",
			wantErr:       nil,
			wantFirstname: "John",
			wantLastname:  "Doe",
			wantEmail:     validEmail,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Сохраняем время до создания пользователя для проверки временных меток
			beforeCreate := time.Now().UTC()
			user, err := New(tt.email, tt.firstname, tt.lastname, tt.passwordHash)
			afterCreate := time.Now().UTC()

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.Nil(t, user)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, user)
			assert.Equal(t, tt.wantFirstname, user.firstname)
			assert.Equal(t, tt.wantLastname, user.lastname)
			assert.Equal(t, tt.wantEmail, user.email)
			assert.Equal(t, tt.passwordHash, user.passwordHash)
			assert.NotEqual(t, uuid.Nil, user.id)
			assert.False(t, user.createdAt.IsZero())
			assert.False(t, user.updatedAt.IsZero())
			// Проверяем, что время создания находится в разумном диапазоне (между before и after)
			assert.True(t, !user.createdAt.Before(beforeCreate) && !user.createdAt.After(afterCreate), "CreatedAt should be within the creation window")
			assert.True(t, !user.updatedAt.Before(beforeCreate) && !user.updatedAt.After(afterCreate), "UpdatedAt should be within the creation window")
			// В конструкторе New createdAt и updatedAt должны быть равны (или очень близки)
			assert.True(t, user.createdAt.Equal(user.updatedAt) || user.createdAt.Sub(user.updatedAt).Abs() < time.Millisecond)
		})
	}
}

func TestUser_From(t *testing.T) {
	id := uuid.New()
	now := time.Now().UTC()
	past := now.Add(-1 * time.Hour)
	email := Email("test@example.com")
	hash := PasswordHash("hashed_password")

	user := From(id, email, "John", "Doe", hash, past, now)

	require.NotNil(t, user)
	assert.Equal(t, id, user.id)
	assert.Equal(t, email, user.email)
	assert.Equal(t, "John", user.firstname)
	assert.Equal(t, "Doe", user.lastname)
	assert.Equal(t, hash, user.passwordHash)
	assert.Equal(t, past, user.createdAt)
	assert.Equal(t, now, user.updatedAt)
}

func TestUser_Accessors(t *testing.T) {
	email := Email("accessor@test.com")
	hash := PasswordHash("hashed_password")
	user, err := New(email, "John", "Doe", hash)
	require.NoError(t, err)

	tests := []struct {
		name     string
		got      any
		expected any
	}{
		{"ID", user.ID(), user.id},
		{"Email", user.Email(), user.email},
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
