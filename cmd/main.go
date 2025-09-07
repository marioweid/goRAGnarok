package main

import (
	"log"
	"net/http"
	"os"

	"goRAGnarok/internal"
	"goRAGnarok/internal/handlers"
	"goRAGnarok/pkg"
)

func main() {
	if err := pkg.LoadEnv(".env"); err != nil {
		log.Println("Warning: .env file not loaded:", err)
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY is not set. Please set it in your environment or .env file before starting the server.")
	}

	srv := &internal.Server{APIKey: apiKey}

	http.HandleFunc("/health", handlers.HealthCheckHandler)
	http.HandleFunc("/v1/response", handlers.ResponseHandler(srv))

	log.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server error:", err)
	}
}
