package scenario

import "github.com/devathh/scene-ai/pkg/consts"

type status int

func (s status) Int() int {
	return int(s)
}

const (
	STATUS_UNKNOWN status = iota
	STATUS_GENERATING
	STATUS_GENERATED
	STATUS_MODIFIED
)

func NewStatus(raw int) (status, error) {
	if raw <= STATUS_UNKNOWN.Int() || raw > STATUS_MODIFIED.Int() {
		return STATUS_UNKNOWN, consts.ErrInvalidStatus
	}

	return status(raw), nil
}
