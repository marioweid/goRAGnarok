package providers

import (
	"bytes"
	"encoding/json"
	"goRAGnarok/internal/models"
	"io"
	"net/http"
)

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

type OllamaProvider struct {
	BaseURL string
}

func (openAiProvider *OllamaProvider) Generate(request models.GenerateRequest) (models.AiResponse, error) {
	payload := map[string]any{
		"model":  request.Model,
		"prompt": request.Input,
		"stream": false,
	}
	payloadBytes, _ := json.Marshal(payload)
	client := &http.Client{}
	url := openAiProvider.BaseURL + "/api/generate"
	openaiReq, err := http.NewRequest("POST", url, bytes.NewReader(payloadBytes))
	if err != nil {
		return models.AiResponse{}, err
	}
	openaiReq.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := client.Do(openaiReq)
	if err != nil {
		return models.AiResponse{}, err
	}

	// Read the response body as a string
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.AiResponse{}, err
	}

	var ollamaResponse OllamaResponse
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

func (openAiProvider *OllamaProvider) Embeddings(request models.EmbeddingsRequest) (models.EmbeddingsResponse, error) {
	// Call OpenAi for Embeddings
	// TODO
	payload := map[string]any{
		"input": request.Input,
		"model": request.Model,
	}
	payloadBytes, _ := json.Marshal(payload)
	client := &http.Client{}
	url := openAiProvider.BaseURL + "/embeddings"
	openaiReq, err := http.NewRequest("POST", url, bytes.NewReader(payloadBytes))
	if err != nil {
		return models.EmbeddingsResponse{}, err
	}
	openaiReq.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(openaiReq)
	if err != nil {
		return models.EmbeddingsResponse{}, err
	}

	// Prepare response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.EmbeddingsResponse{}, err
	}

	var openAiResponse models.OpenAiResponse
	err = json.Unmarshal(respBody, &openAiResponse)
	if err != nil {
		return models.EmbeddingsResponse{}, err
	}
	defer resp.Body.Close()
	return models.EmbeddingsResponse{}, nil
}
