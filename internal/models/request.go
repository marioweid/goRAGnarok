package models

type GenerateRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

type AiResponse struct {
	Response string
	Model    string
	Role     string
}

type EmbeddingsRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

type EmbeddingsResponse struct {
	Model      string    `json:"model"`
	Embeddings []float64 `json:"embeddings"`
}
