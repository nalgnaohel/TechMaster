package models

type GroqClient struct {
	ApiKey string
}

type GroqMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GroqRequest struct {
	Messages    []GroqMessage `json:"messages"`
	LLMModel    string        `json:"model"`
	MaxTokens   int           `json:"max_tokens"`
	Temperature float64       `json:"temperature"`
	TopP        int           `json:"top_p"`
	Stream      bool          `json:"stream"`
	Stop        interface{}   `json:"stop"`
}
