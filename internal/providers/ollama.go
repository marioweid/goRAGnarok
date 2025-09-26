package providers

import (
	"bytes"
	"encoding/json"
	"goRAGnarok/internal/interfaces"
	"io"
	"net/http"
)

type OllamaProvider struct {
	BaseURL string
}

// ```shell
// curl http://localhost:11434/api/generate -d '{
//   "model": "llama3.2",
//   "prompt":"Why is the sky blue?"
// }'
// ```

func (openAiProvider *OllamaProvider) Generate(request interfaces.GenerateRequest) (string, error) {
	payload := map[string]any{
		"model":  request.Model,
		"prompt": request.Input,
	}
	payloadBytes, _ := json.Marshal(payload)
	client := &http.Client{}
	url := openAiProvider.BaseURL + "/api/generate"
	openaiReq, err := http.NewRequest("POST", url, bytes.NewReader(payloadBytes))
	if err != nil {
		return "", err
	}
	openaiReq.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := client.Do(openaiReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close() // Ensure body is closed

	// Read the response body as a string
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(respBody), nil
}
