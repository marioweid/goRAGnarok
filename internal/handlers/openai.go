package handlers

import (
	"bytes"
	"encoding/json"
	"goRAGnarok/internal/interfaces"
	"net/http"
)

func CallOpenAIEmbeddings(s *interfaces.Server, input, model, encodingFormat string) (*http.Response, error) {
	payload := map[string]any{
		"input":           input,
		"model":           model,
		"encoding_format": encodingFormat,
	}
	payloadBytes, _ := json.Marshal(payload)
	client := &http.Client{}
	url := s.BaseURL + "/embeddings"
	req, err := http.NewRequest("POST", url, bytes.NewReader(payloadBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.APIKey)
	return client.Do(req)
}
