package consts

type contextKey string

const (
	ClientIP  contextKey = "x-client-ip"
	UserAgent contextKey = "x-user-agent"
)
