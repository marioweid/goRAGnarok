package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"goRAGnarok/internal/models"
	"io"
	"net/http"
	"time"
)

type OllamaProvider struct {
	BaseURL        string
	EmbeddingModel string
	client         *http.Client
}

func NewOllamaProvider(baseURL, embeddingModel string) *OllamaProvider {
	return &OllamaProvider{
		BaseURL:        baseURL,
		EmbeddingModel: embeddingModel,
		client: &http.Client{
			Timeout: time.Second * 120, // Because my MacBook 2020 is a potato
		},
	}
}

func (ollamaProvider *OllamaProvider) Generate(ctx context.Context, request models.GenerateRequest) (models.AiResponse, error) {
	payload := map[string]any{
		"model":  request.Model,
		"prompt": request.Input,
		"stream": false,
	}
	payloadBytes, _ := json.Marshal(payload)
	url := ollamaProvider.BaseURL + "/api/generate"
	ollamaReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(payloadBytes))
	if err != nil {
		return models.AiResponse{}, err
	}
	ollamaReq.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := ollamaProvider.client.Do(ollamaReq)
	if err != nil {
		return models.AiResponse{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return models.AiResponse{}, fmt.Errorf("ollama API error: status %d", resp.StatusCode)
	}

	// Read the response body as a string
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.AiResponse{}, err
	}

	var ollamaResponse models.OllamaResponse
	err = json.Unmarshal(respBody, &ollamaResponse)
	if err != nil {
		return models.AiResponse{}, err
	}

	defer resp.Body.Close() // Ensure body is closed
	return models.AiResponse{
		Response: ollamaResponse.Response,
		Model:    ollamaResponse.Model,
		Role:     "assistant",
	}, nil
}

func (ollamaProvider *OllamaProvider) Embeddings(ctx context.Context, request models.EmbeddingsRequest) (models.EmbeddingsResponse, error) {
	// Call Ollama for Embeddings
	payload := map[string]any{
		"prompt": request.Input,
		"model":  ollamaProvider.EmbeddingModel,
	}
	payloadBytes, _ := json.Marshal(payload)
	url := ollamaProvider.BaseURL + "/api/embeddings"
	ollamaReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(payloadBytes))
	if err != nil {
		return models.EmbeddingsResponse{}, err
	}
	ollamaReq.Header.Set("Content-Type", "application/json")
	resp, err := ollamaProvider.client.Do(ollamaReq)
	if err != nil {
		return models.EmbeddingsResponse{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return models.EmbeddingsResponse{}, fmt.Errorf("ollama API error: status %d", resp.StatusCode)
	}

	// Prepare response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.EmbeddingsResponse{}, err
	}

	var ollamaEmbeddings models.OllamaEmbeddingsResponse
	err = json.Unmarshal(respBody, &ollamaEmbeddings)
	if err != nil {
		return models.EmbeddingsResponse{}, err
	}
	defer resp.Body.Close()
	return models.EmbeddingsResponse{Model: ollamaProvider.EmbeddingModel, Embeddings: ollamaEmbeddings.Embedding}, nil
}
