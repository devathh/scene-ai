package scenario

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatusInt(t *testing.T) {
	tests := []struct {
		name     string
		status   status
		expected int
	}{
		{
			name:     "UnknownStatus",
			status:   STATUS_UNKNOWN,
			expected: 0,
		},
		{
			name:     "GeneratedStatus",
			status:   STATUS_GENERATED,
			expected: 1,
		},
		{
			name:     "ModifiedStatus",
			status:   STATUS_MODIFIED,
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.status.Int()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestStatusConstants(t *testing.T) {
	assert.Equal(t, 0, int(STATUS_UNKNOWN))
	assert.Equal(t, 1, int(STATUS_GENERATED))
	assert.Equal(t, 2, int(STATUS_MODIFIED))

	var s status
	assert.Equal(t, STATUS_UNKNOWN, s)
}
