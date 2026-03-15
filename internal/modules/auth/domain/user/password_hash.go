package user

import (
	"strings"

	"github.com/devathh/scene-ai/pkg/consts"
	"golang.org/x/crypto/bcrypt"
)

type PasswordHash string

func (ph PasswordHash) Compare(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(ph), []byte(password)) == nil
}

func (ph PasswordHash) String() string {
	return string(ph)
}

func NewPasswordHash(raw string) (PasswordHash, error) {
	raw = strings.TrimSpace(raw)
	if len([]rune(raw)) < 6 {
		return "", consts.ErrInvalidPassword
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
	if err != nil {
		return "", consts.ErrInvalidPassword
	}

	return PasswordHash(hash), nil
}
