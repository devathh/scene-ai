package sessionredis

import (
	"time"

	"github.com/google/uuid"
)

type SessionModel struct {
	UserID      uuid.UUID `json:"user_id"`
	FingerPrint string    `json:"finger_print"`
	CreatedAt   time.Time `json:"created_at"`
}
