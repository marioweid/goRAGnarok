package handlers

import (
	"bytes"
	"encoding/json"
	"goRAGnarok/internal"
	"io"
	"net/http"
)

func CallOpenAIEmbeddings(s *internal.Server, input, model, encodingFormat string) (*http.Response, error) {
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

func ResponseHandler(s *internal.Server) http.HandlerFunc {
	// TODO: datatypes for responses
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req struct {
			Model string `json:"model"`
			Input string `json:"input"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		if req.Model == "" {
			req.Model = "gpt-4.1"
		}

		payload := map[string]any{
			"model": req.Model,
			"input": req.Input,
		}
		payloadBytes, _ := json.Marshal(payload)
		client := &http.Client{}
		url := s.BaseURL + "/responses"
		openaiReq, err := http.NewRequest("POST", url, bytes.NewReader(payloadBytes))
		if err != nil {
			http.Error(w, "Failed to create request", http.StatusInternalServerError)
			return
		}
		openaiReq.Header.Set("Content-Type", "application/json")
		openaiReq.Header.Set("Authorization", "Bearer "+s.APIKey)

		resp, err := client.Do(openaiReq)
		if err != nil {
			http.Error(w, "Failed to call OpenAI API", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	}
}

func EmbeddingsHandler(s *internal.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req struct {
			Input          string `json:"input"`
			Model          string `json:"model"`
			EncodingFormat string `json:"encoding_format"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if req.Model == "" {
			req.Model = "text-embedding-3-small"
		}
		if req.EncodingFormat == "" {
			req.EncodingFormat = "float"
		}

		resp, err := CallOpenAIEmbeddings(s, req.Input, req.Model, req.EncodingFormat)
		if err != nil {
			http.Error(w, "Failed to call OpenAI API", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	}
}
