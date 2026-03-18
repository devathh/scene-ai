package dtos

import "time"

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

type CreateScenarioRequest struct {
	Title             string               `json:"title"`
	ScenarioPrompt    string               `json:"scenario_prompt"`
	GlobalStylePrompt string               `json:"global_style_prompt"`
	Scenes            []CreateSceneRequest `json:"scenes"`
}

type CreateSceneRequest struct {
	Order       int           `json:"order"`
	Title       string        `json:"title"`
	Duration    time.Duration `json:"duration"`
	VideoPrompt string        `json:"video_prompt"`
}

type UpdateScenarioRequest struct {
	Title             *string              `json:"title,omitempty"`
	ScenarioPrompt    *string              `json:"scenario_prompt,omitempty"`
	GlobalStylePrompt *string              `json:"global_style_prompt,omitempty"`
	Scenes            []UpdateSceneRequest `json:"scenes,omitempty"`
}

type UpdateSceneRequest struct {
	ID          *string        `json:"id,omitempty"`
	Order       *int           `json:"order,omitempty"`
	Title       *string        `json:"title,omitempty"`
	Duration    *time.Duration `json:"duration,omitempty"`
	VideoPrompt *string        `json:"video_prompt,omitempty"`
}
