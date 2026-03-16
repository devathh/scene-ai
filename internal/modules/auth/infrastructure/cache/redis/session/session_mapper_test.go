package sessionredis

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/devathh/scene-ai/internal/modules/auth/domain/session"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToBytes_Success(t *testing.T) {
	userID := uuid.New()
	fingerPrint := "test-fingerprint-12345"
	createdAt := time.Date(2023, 10, 5, 14, 30, 0, 0, time.UTC)

	domainSession := &session.Session{
		UserID:      userID,
		FingerPrint: fingerPrint,
		CreatedAt:   createdAt,
	}

	bytes, err := ToBytes(domainSession)

	require.NoError(t, err)
	require.NotNil(t, bytes)
	require.NotEmpty(t, bytes)

	var model SessionModel
	err = json.Unmarshal(bytes, &model)
	require.NoError(t, err)

	assert.Equal(t, userID, model.UserID)
	assert.Equal(t, fingerPrint, model.FingerPrint)
	assert.Equal(t, createdAt.Unix(), model.CreatedAt.Unix())
}

func TestToBytes_InvalidData(t *testing.T) {
	userID := uuid.New()
	domainSession := &session.Session{
		UserID:      userID,
		FingerPrint: "valid-fp",
		CreatedAt:   time.Now(),
	}

	bytes, err := ToBytes(domainSession)
	assert.NoError(t, err)
	assert.NotNil(t, bytes)
}

func TestToModel_Success(t *testing.T) {
	userID := uuid.New()
	fingerPrint := "browser-chrome-linux"
	createdAt := time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC)

	model := SessionModel{
		UserID:      userID,
		FingerPrint: fingerPrint,
		CreatedAt:   createdAt,
	}

	bytes, err := json.Marshal(model)
	require.NoError(t, err)

	resultSession, err := ToModel(bytes)

	require.NoError(t, err)
	require.NotNil(t, resultSession)

	assert.Equal(t, userID, resultSession.UserID)
	assert.Equal(t, fingerPrint, resultSession.FingerPrint)
	assert.Equal(t, createdAt.Unix(), resultSession.CreatedAt.Unix())
}

func TestToModel_InvalidJSON(t *testing.T) {
	invalidBytes := []byte("{ invalid json }")

	sessionObj, err := ToModel(invalidBytes)

	assert.Error(t, err)
	assert.Nil(t, sessionObj)
	assert.Contains(t, err.Error(), "failed to unmarshal bytes")
}

func TestToModel_EmptyBytes(t *testing.T) {
	emptyBytes := []byte{}

	sessionObj, err := ToModel(emptyBytes)

	assert.Error(t, err)
	assert.Nil(t, sessionObj)
}

func TestToModel_MissingFields(t *testing.T) {
	jsonBytes := []byte(`{"finger_print": "some-fp"}`)

	sessionObj, err := ToModel(jsonBytes)

	require.NoError(t, err)
	require.NotNil(t, sessionObj)

	assert.Equal(t, uuid.Nil, sessionObj.UserID)
	assert.Equal(t, "some-fp", sessionObj.FingerPrint)
	assert.True(t, sessionObj.CreatedAt.IsZero())
}

func TestRoundTrip(t *testing.T) {
	originalUserID := uuid.New()
	originalFP := "round-trip-test"
	originalTime := time.Now().UTC().Truncate(time.Second)

	original := &session.Session{
		UserID:      originalUserID,
		FingerPrint: originalFP,
		CreatedAt:   originalTime,
	}

	bytes, err := ToBytes(original)
	require.NoError(t, err)

	restored, err := ToModel(bytes)
	require.NoError(t, err)

	assert.Equal(t, original.UserID, restored.UserID)
	assert.Equal(t, original.FingerPrint, restored.FingerPrint)
	assert.Equal(t, original.CreatedAt.Unix(), restored.CreatedAt.Unix())
}
