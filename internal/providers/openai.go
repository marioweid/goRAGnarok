package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"goRAGnarok/internal/models"
	"io"
	"net/http"
	"time"
)

type OpenAiProvider struct {
	BaseURL        string
	ApiKey         string
	EmbeddingModel string
	client         *http.Client
}

func NewOpenAiProvider(baseURL, apiKey, embeddingModel string) *OpenAiProvider {
	return &OpenAiProvider{
		BaseURL:        baseURL,
		ApiKey:         apiKey,
		EmbeddingModel: embeddingModel,
		client: &http.Client{
			Timeout: time.Second * 30,
		},
	}
}

func (openAiProvider *OpenAiProvider) Generate(request models.GenerateRequest) (models.AiResponse, error) {
	payload := map[string]any{
		"model": request.Model,
		"input": request.Input,
	}
	payloadBytes, _ := json.Marshal(payload)
	url := openAiProvider.BaseURL + "/responses"
	openaiReq, err := http.NewRequest("POST", url, bytes.NewReader(payloadBytes))
	if err != nil {
		return models.AiResponse{}, err
	}
	openaiReq.Header.Set("Content-Type", "application/json")
	openaiReq.Header.Set("Authorization", "Bearer "+openAiProvider.ApiKey)

	resp, err := openAiProvider.client.Do(openaiReq)
	if err != nil {
		return models.AiResponse{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return models.AiResponse{}, fmt.Errorf("ollama API error: status %d", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.AiResponse{}, err
	}

	var openAiResponse models.OpenAiResponse
	err = json.Unmarshal(respBody, &openAiResponse)
	if err != nil {
		return models.AiResponse{}, err
	}
	defer resp.Body.Close()
	return models.AiResponse{
		Response: openAiResponse.Output[0].Content[0].Text,
		Model:    openAiResponse.Model,
		Role:     openAiResponse.Output[0].Role,
	}, nil
}

func (openAiProvider *OpenAiProvider) Embeddings(request models.EmbeddingsRequest) (models.EmbeddingsResponse, error) {
	// Call OpenAi for Embeddings
	payload := map[string]any{
		"input": request.Input,
		"model": openAiProvider.EmbeddingModel,
	}
	payloadBytes, _ := json.Marshal(payload)
	url := openAiProvider.BaseURL + "/embeddings"
	openaiReq, err := http.NewRequest("POST", url, bytes.NewReader(payloadBytes))
	if err != nil {
		return models.EmbeddingsResponse{}, err
	}
	openaiReq.Header.Set("Content-Type", "application/json")
	openaiReq.Header.Set("Authorization", "Bearer "+openAiProvider.ApiKey)
	resp, err := openAiProvider.client.Do(openaiReq)
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

	var openAiEmbeddings models.OpenAiEmbeddingsResponse
	err = json.Unmarshal(respBody, &openAiEmbeddings)
	if err != nil {
		return models.EmbeddingsResponse{}, err
	}
	defer resp.Body.Close()
	return models.EmbeddingsResponse{Model: openAiEmbeddings.Model, Embeddings: openAiEmbeddings.Data[0].Embedding}, nil
}
