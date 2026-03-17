package dtos

type Token struct {
	Access     string `json:"access"`
	AccessTTL  int64  `json:"access-ttl"`
	Refresh    string `json:"refresh"`
	RefreshTTL int64  `json:"refresh-ttl"`
}

type ErrMsg struct {
	Error string `json:"error"`
}

type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname,omitempty"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}
