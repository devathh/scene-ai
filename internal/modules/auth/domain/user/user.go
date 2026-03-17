package user

import (
	"strings"
	"time"

	"github.com/devathh/scene-ai/pkg/consts"
	"github.com/google/uuid"
)

type User struct {
	id           uuid.UUID
	email        Email
	firstname    string
	lastname     string
	passwordHash PasswordHash
	createdAt    time.Time
	updatedAt    time.Time
}

func New(
	email Email,
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

	if !email.IsValid() {
		return nil, consts.ErrInvalidEmail
	}

	return &User{
		id:           uuid.New(),
		email:        email,
		firstname:    firstname,
		lastname:     lastname,
		passwordHash: passwordHash,
		createdAt:    time.Now().UTC(),
		updatedAt:    time.Now().UTC(),
	}, nil
}

func From(
	id uuid.UUID,
	email Email,
	firstname, lastname string,
	passwordHash PasswordHash,
	createdAt, updatedAt time.Time,
) *User {
	return &User{
		id:           id,
		email:        email,
		firstname:    firstname,
		lastname:     lastname,
		passwordHash: passwordHash,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}
}

type UpdateMask struct {
	Email     *Email
	Firstname *string
	Lastname  *string
	Password  *string
}

func (u *User) Apply(mask UpdateMask) error {
	updated := false

	if mask.Email != nil {
		if !mask.Email.IsValid() {
			return consts.ErrInvalidEmail
		}
		u.email = *mask.Email
		updated = true
	}

	if mask.Firstname != nil {
		name := strings.TrimSpace(*mask.Firstname)
		if len([]rune(name)) < 3 || len([]rune(name)) > 64 {
			return consts.ErrInvalidFirstname
		}
		u.firstname = name
		updated = true
	}

	if mask.Lastname != nil {
		name := strings.TrimSpace(*mask.Lastname)
		if *mask.Lastname != "" && (len([]rune(name)) < 3 || len([]rune(name)) > 64) {
			return consts.ErrInvalidLastname
		}
		u.lastname = name
		updated = true
	}

	if mask.Password != nil {
		newHash, err := NewPasswordHash(*mask.Password)
		if err != nil {
			return err
		}
		u.passwordHash = newHash
		updated = true
	}

	if updated {
		u.updatedAt = time.Now().UTC()
	}

	return nil
}

func (u *User) ID() uuid.UUID {
	return u.id
}

func (u *User) Email() Email {
	return u.email
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
