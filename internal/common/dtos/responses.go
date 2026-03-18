package dtos

import "time"

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

type Scenario struct {
	ID                string  `json:"id"`
	AuthorID          string  `json:"author_id"`
	Title             string  `json:"title"`
	ScenarioPrompt    string  `json:"scenario_prompt"`
	GlobalStylePrompt string  `json:"global_style_prompt"`
	Status            int     `json:"status"`
	Scenes            []Scene `json:"scenes"`
	CreatedAt         int64   `json:"created_at"`
	UpdatedAt         int64   `json:"updated_at"`
}

type Scene struct {
	ID          string        `json:"id"`
	Order       int           `json:"order"`
	Title       string        `json:"title"`
	Duration    time.Duration `json:"duration"`
	VideoPrompt string        `json:"video_prompt"`
}

type Scenarios struct {
	Scenarios []*Scenario `json:"scenarios"`
}
