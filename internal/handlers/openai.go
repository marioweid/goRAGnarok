package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"goRAGnarok/internal/interfaces"
	"io"
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

func ResponseHandler(providerLookup map[string]interfaces.Provider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// check correct http method
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		// Parse request body
		var req struct {
			Model string `json:"model"`
			Input string `json:"input"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		// Gemma3 4b as default
		if req.Model == "" {
			req.Model = "gemma3:4b"
		}
		llmResponse := providerLookup[req.Model].Generate()
		fmt.Println(llmResponse)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Encode the struct into JSON
		json.NewEncoder(w).Encode(llmResponse)
	}
}

func EmbeddingsHandler(s *interfaces.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
