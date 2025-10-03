package interfaces

import (
	"context"
	"goRAGnarok/internal/models"
)

type Server struct {
	BaseURL string
	APIKey  string
}

type Provider interface {
	Generate(ctx context.Context, request models.GenerateRequest) (models.AiResponse, error)
	Embeddings(ctx context.Context, request models.EmbeddingsRequest) (models.EmbeddingsResponse, error)
}
