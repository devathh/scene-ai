package dtos

type RegisterRequest struct {
	Email     string `json:"email"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Password  string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserRequest struct {
	Email     *string `json:"email,omitempty"`
	Firstname *string `json:"firstname,omitempty"`
	Lastname  *string `json:"lastname,omitempty"`
	Password  *string `json:"password,omitempty"`
}
