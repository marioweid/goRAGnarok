package main

import (
	"database/sql"
	"log"
	"net/http"

	"goRAGnarok/internal/handlers"
	"goRAGnarok/internal/interfaces"
	"goRAGnarok/internal/providers"
	"goRAGnarok/pkg"

	_ "github.com/lib/pq"
)

func main() {
	// Load configuration
	cfg, err := pkg.NewConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Context classes
	srv := &interfaces.Server{APIKey: cfg.OpenAIAPIKey, BaseURL: cfg.OpenAIBaseURL}
	providerLookup := make(map[string]interfaces.Provider)

	ollamaProvider := providers.NewOllamaProvider(cfg.OllamaBaseURL, cfg.OllamaEmbeddingModel)
	openAiProvider := providers.NewOpenAiProvider(cfg.OpenAIBaseURL, cfg.OpenAIAPIKey, cfg.OpenAIEmbeddingModel)

	providerLookup["gemma3:4b"] = ollamaProvider
	providerLookup["gpt-4.1"] = openAiProvider

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	http.HandleFunc("/health", handlers.HealthCheckHandler)

	http.HandleFunc("/v1/response", handlers.ResponseHandler(providerLookup))
	http.HandleFunc("/v1/embeddings", handlers.EmbeddingsHandler(providerLookup))
	http.HandleFunc("/v1/similarity-search", handlers.SimilaritySearchHandler(srv, db))

	log.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server error:", err)
	}
}
