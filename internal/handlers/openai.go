package handlers

import (
	"bytes"
	"encoding/json"
	"goRAGnarok/internal"
	"io"
	"net/http"
)

func ResponseHandler(s *internal.Server) http.HandlerFunc {
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

		payload := map[string]interface{}{
			"model": req.Model,
			"input": req.Input,
		}
		payloadBytes, _ := json.Marshal(payload)
		client := &http.Client{}
		openaiReq, err := http.NewRequest("POST", "https://api.openai.com/v1/responses", bytes.NewReader(payloadBytes))
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
