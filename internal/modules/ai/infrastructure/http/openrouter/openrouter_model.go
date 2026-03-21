package openrouterhttp

type request struct {
	Model    string    `json:"model"`
	Messages []message `json:"messages"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type response struct {
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type responseScenario struct {
	Title             string `json:"title"`
	GlobalStylePrompt string `json:"global_style_prompt"`
}

type responseScene struct {
	Order       int    `json:"order"`
	Title       string `json:"title"`
	DurationSec int    `json:"duration"`
	VideoPrompt string `json:"video_prompt"`
	Status      string `json:"status"`
}
