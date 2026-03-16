package session

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	UserID      uuid.UUID
	FingerPrint string
	CreatedAt   time.Time
}
