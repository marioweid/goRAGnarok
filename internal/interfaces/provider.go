package interfaces

import (
	"goRAGnarok/internal/models"
)

type Server struct {
	BaseURL string
	APIKey  string
}

type Provider interface {
	Generate(request models.GenerateRequest) (models.AiResponse, error)
	Embeddings(request models.EmbeddingsRequest) (models.EmbeddingsResponse, error)
}
