package sessionredis

import (
	"encoding/json"
	"fmt"

	"github.com/devathh/scene-ai/internal/modules/auth/domain/session"
)

func ToBytes(domain *session.Session) ([]byte, error) {
	bytes, err := json.Marshal(&SessionModel{
		UserID:      domain.UserID,
		FingerPrint: domain.FingerPrint,
		CreatedAt:   domain.CreatedAt,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal model: %w", err)
	}

	return bytes, nil
}

func ToModel(bytes []byte) (*session.Session, error) {
	var model SessionModel
	if err := json.Unmarshal(bytes, &model); err != nil {
		return nil, fmt.Errorf("failed to unmarshal bytes: %w", err)
	}

	return &session.Session{
		UserID:      model.UserID,
		FingerPrint: model.FingerPrint,
		CreatedAt:   model.CreatedAt,
	}, nil
}
