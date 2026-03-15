package authuserpg

import (
	"testing"
	"time"

	"github.com/devathh/scene-ai/internal/modules/auth/domain/user"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestToDomain(t *testing.T) {
	now := time.Now().UTC()
	id := uuid.New()
	hash := user.PasswordHash("hashed_secret")

	tests := []struct {
		name          string
		model         *UserModel
		wantID        uuid.UUID
		wantFirstname string
		wantLastname  string
		wantHash      user.PasswordHash
		wantCreated   time.Time
		wantUpdated   time.Time
	}{
		{
			name: "Success_ValidModel",
			model: &UserModel{
				ID:           id,
				Firstname:    "John",
				Lastname:     "Doe",
				PasswordHash: hash.String(),
				CreatedAt:    now,
				UpdatedAt:    now,
			},
			wantID:        id,
			wantFirstname: "John",
			wantLastname:  "Doe",
			wantHash:      hash,
			wantCreated:   now,
			wantUpdated:   now,
		},
		{
			name: "Success_EmptyLastname",
			model: &UserModel{
				ID:           id,
				Firstname:    "John",
				Lastname:     "",
				PasswordHash: hash.String(),
				CreatedAt:    now,
				UpdatedAt:    now,
			},
			wantID:        id,
			wantFirstname: "John",
			wantLastname:  "",
			wantHash:      hash,
			wantCreated:   now,
			wantUpdated:   now,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToDomain(tt.model)

			assert.NotNil(t, result)
			assert.Equal(t, tt.wantID, result.ID())
			assert.Equal(t, tt.wantFirstname, result.Firstname())
			assert.Equal(t, tt.wantLastname, result.Lastname())
			assert.Equal(t, tt.wantHash, result.PasswordHash())
			assert.True(t, tt.wantCreated.Equal(result.CreatedAt()))
			assert.True(t, tt.wantUpdated.Equal(result.UpdatedAt()))
		})
	}
}

func TestToModel(t *testing.T) {
	now := time.Now().UTC()
	id := uuid.New()
	hash := user.PasswordHash("hashed_secret")

	// Create domain user using From to bypass validation and set specific ID/Times
	domainUser := user.From(
		id,
		"John",
		"Doe",
		hash,
		now,
		now,
	)

	tests := []struct {
		name          string
		domain        *user.User
		wantID        uuid.UUID
		wantFirstname string
		wantLastname  string
		wantHash      string
		wantCreated   time.Time
		wantUpdated   time.Time
	}{
		{
			name:          "Success_ValidDomain",
			domain:        domainUser,
			wantID:        id,
			wantFirstname: "John",
			wantLastname:  "Doe",
			wantHash:      hash.String(),
			wantCreated:   now,
			wantUpdated:   now,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToModel(tt.domain)

			assert.NotNil(t, result)
			assert.Equal(t, tt.wantID, result.ID)
			assert.Equal(t, tt.wantFirstname, result.Firstname)
			assert.Equal(t, tt.wantLastname, result.Lastname)
			assert.Equal(t, tt.wantHash, result.PasswordHash)
			assert.True(t, tt.wantCreated.Equal(result.CreatedAt))
			assert.True(t, tt.wantUpdated.Equal(result.UpdatedAt))
		})
	}
}
