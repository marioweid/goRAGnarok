package handlers

import (
	"encoding/json"
	"goRAGnarok/internal/interfaces"
	"goRAGnarok/internal/models"
	"net/http"
)

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func ResponseHandler(providerLookup map[string]interfaces.Provider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// check correct http method
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		// Parse request body
		var req models.GenerateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		llmResponse, err := providerLookup[req.Model].Generate(req)
		if err != nil {
			http.Error(w, "Failed to create request", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Encode the struct into JSON
		json.NewEncoder(w).Encode(llmResponse)
	}
}

func EmbeddingsHandler(providerLookup map[string]interfaces.Provider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.EmbeddingsRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		llmEmbeddings, err := providerLookup[req.Model].Embeddings(req)
		if err != nil {
			http.Error(w, "Failed to create embeddings", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(llmEmbeddings)
	}
}
