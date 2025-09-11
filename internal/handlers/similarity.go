package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"goRAGnarok/internal"
	"goRAGnarok/internal/database"
	"io"
	"net/http"
)

type SimilaritySearchRequest struct {
	Input string `json:"input"`
	TopN  int    `json:"top_n"`
	Model string `json:"model,omitempty"`
}

type SimilaritySearchResponse struct {
	Results []database.SearchResult `json:"results"`
}

func SimilaritySearchHandler(s *internal.Server, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req SimilaritySearchRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		if req.TopN <= 0 {
			req.TopN = 5
		}
		if req.Model == "" {
			req.Model = "text-embedding-3-small"
		}

		// Get embedding from OpenAI
		resp, err := CallOpenAIEmbeddings(s, req.Input, req.Model, "float")
		if err != nil {
			http.Error(w, "Failed to get embedding", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		var embedResp struct {
			Data []struct {
				Embedding []float32 `json:"embedding"`
			} `json:"data"`
		}
		body, _ := io.ReadAll(resp.Body)
		if err := json.Unmarshal(body, &embedResp); err != nil || len(embedResp.Data) == 0 {
			http.Error(w, "Invalid embedding response", http.StatusInternalServerError)
			return
		}
		vector := embedResp.Data[0].Embedding

		// Search database
		results, err := database.SimilaritySearch(context.Background(), db, vector, req.TopN)
		if err != nil {
			http.Error(w, "Database search error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(SimilaritySearchResponse{Results: results})
	}
}
