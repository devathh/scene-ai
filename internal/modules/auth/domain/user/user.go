package user

import (
	"strings"
	"time"

	"github.com/devathh/scene-ai/pkg/consts"
	"github.com/google/uuid"
)

type User struct {
	id           uuid.UUID
	firstname    string
	lastname     string
	passwordHash PasswordHash
	createdAt    time.Time
	updatedAt    time.Time
}

func New(
	firstname, lastname string,
	passwordHash PasswordHash,
) (*User, error) {
	firstname = strings.TrimSpace(firstname)
	if len([]rune(firstname)) < 3 || len([]rune(firstname)) > 64 {
		return nil, consts.ErrInvalidFirstname
	}

	lastname = strings.TrimSpace(lastname)
	if lastname != "" &&
		(len([]rune(lastname)) < 3 || len([]rune(lastname)) > 64) {
		return nil, consts.ErrInvalidLastname
	}

	return &User{
		id:           uuid.New(),
		firstname:    firstname,
		lastname:     lastname,
		passwordHash: passwordHash,
		createdAt:    time.Now().UTC(),
		updatedAt:    time.Now().UTC(),
	}, nil
}

func From(
	id uuid.UUID,
	firstname, lastname string,
	passwordHash PasswordHash,
	createdAt, updatedAt time.Time,
) *User {
	return &User{
		id:           id,
		firstname:    firstname,
		lastname:     lastname,
		passwordHash: passwordHash,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}
}

func (u *User) ID() uuid.UUID {
	return u.id
}

func (u *User) Firstname() string {
	return u.firstname
}

func (u *User) Lastname() string {
	return u.lastname
}

func (u *User) PasswordHash() PasswordHash {
	return u.passwordHash
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}
