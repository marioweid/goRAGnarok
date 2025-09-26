package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"goRAGnarok/internal/handlers"
	"goRAGnarok/internal/interfaces"
	"goRAGnarok/internal/providers"
	"goRAGnarok/pkg"

	_ "github.com/lib/pq"
)

func main() {
	if err := pkg.LoadEnv(".env"); err != nil {
		log.Println("Warning: .env file not loaded:", err)
	}

	// Check OpenAI config in OpenAI Handler
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY is not set. Please set it in your environment or .env file before starting the server.")
	}

	baseURL := os.Getenv("OPENAI_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}

	// Context classes
	srv := &interfaces.Server{APIKey: apiKey, BaseURL: baseURL}
	providerLookup := make(map[string]interfaces.Provider)

	providerLookup["gpt-4.1"] = &providers.OpenAiProvider{BaseURL: baseURL, ApiKey: apiKey}
	providerLookup["gemma3:4b"] = &providers.OllamaProvider{BaseURL: "my_base_url"}

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	http.HandleFunc("/health", handlers.HealthCheckHandler)

	http.HandleFunc("/v1/response", handlers.ResponseHandler(providerLookup))
	http.HandleFunc("/v1/embeddings", handlers.EmbeddingsHandler(srv))
	http.HandleFunc("/v1/similarity-search", handlers.SimilaritySearchHandler(srv, db))

	log.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server error:", err)
	}
}
