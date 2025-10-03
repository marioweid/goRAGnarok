package models

type OllamaMessageContent struct {
	Text string `json:"text"`
	Type string `json:"type"`
}

type OllamaMessage struct {
	Role    string                 `json:"role"`
	Content []OllamaMessageContent `json:"content"`
}

type OllamaResponse struct {
	Model    string `json:"model"`
	Response string `json:"response"`
}

type OllamaEmbeddingsResponse struct {
	Embedding []float64 `json:"embedding"`
}
