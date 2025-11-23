package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"goRAGnarok/internal/database"
	"goRAGnarok/internal/interfaces"
	"goRAGnarok/internal/models"
	"io"
	"net/http"
)

type RAGRequest struct {
	Input string `json:"input"`
	Model string `json:"model"`
	TopN  int    `json:"top_n"`
}

type RAGResponse struct {
	Answer  string                  `json:"answer"`
	Context []database.SearchResult `json:"context"`
}

func RAGHandler(s *interfaces.Server, db *sql.DB, providers map[string]interfaces.Provider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req RAGRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if req.TopN <= 0 {
			req.TopN = 3
		}
		if req.Model == "" {
			req.Model = "gpt-4.1" // Default model
		}

		// 1. Get embedding for the query
		// We use the existing CallOpenAIEmbeddings helper or the provider directly.
		// Since SimilaritySearchHandler uses CallOpenAIEmbeddings, let's stick to that for consistency
		// or use the provider if available. The plan mentioned using OpenAI embeddings as per existing similarity.go.
		// Let's reuse the logic from SimilaritySearchHandler for embedding generation.
		
		// Note: In a more robust design, we might abstract embedding generation further.
		// For now, we'll use the helper from openai.go in handlers package.
		embedResp, err := CallOpenAIEmbeddings(s, req.Input, "text-embedding-3-small", "float")
		if err != nil {
			http.Error(w, "Failed to get embedding", http.StatusInternalServerError)
			return
		}
		defer embedResp.Body.Close()

		var embedData struct {
			Data []struct {
				Embedding []float32 `json:"embedding"`
			} `json:"data"`
		}
		body, _ := io.ReadAll(embedResp.Body)
		if err := json.Unmarshal(body, &embedData); err != nil || len(embedData.Data) == 0 {
			http.Error(w, "Invalid embedding response", http.StatusInternalServerError)
			return
		}
		vector := embedData.Data[0].Embedding

		// 2. Perform similarity search
		searchResults, err := database.SimilaritySearch(context.Background(), db, vector, req.TopN)
		if err != nil {
			http.Error(w, "Database search error", http.StatusInternalServerError)
			return
		}

		// 3. Construct prompt
		contextStr := ""
		for _, res := range searchResults {
			contextStr += fmt.Sprintf("Title: %s\nContent: %s\n\n", res.Title, res.Content)
		}

		prompt := fmt.Sprintf(`You are a helpful assistant. Use the following context to answer the user's question.
If the answer is not in the context, say you don't know.

Context:
%s

User Question: %s`, contextStr, req.Input)

		// 4. Generate answer
		provider, ok := providers[req.Model]
		if !ok {
			http.Error(w, "Model not supported", http.StatusBadRequest)
			return
		}

		genReq := models.GenerateRequest{
			Model: req.Model,
			Input: prompt,
		}

		aiResp, err := provider.Generate(context.Background(), genReq)
		if err != nil {
			http.Error(w, "Generation error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// 5. Return response
		resp := RAGResponse{
			Answer:  aiResp.Response,
			Context: searchResults,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
