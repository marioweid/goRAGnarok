package models

type OpenAiMessageContent struct {
	Text string `json:"text"`
	Type string `json:"type"`
}

type OpenAiMessage struct {
	Role    string                 `json:"role"`
	Content []OpenAiMessageContent `json:"content"`
}

type OpenAiResponse struct {
	Model  string          `json:"model"`
	Output []OpenAiMessage `json:"output"`
}

type OpenAiUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type OpenAiEmbeddingsData struct {
	Object    string    `json:"object"`
	Index     int       `json:"index"`
	Embedding []float64 `json:"embedding"`
}

type OpenAiEmbeddingsResponse struct {
	Object string                 `json:"object"`
	Data   []OpenAiEmbeddingsData `json:"data"`
	Model  string                 `json:"model"`
	Usage  OpenAiUsage            `json:"usage"`
}
