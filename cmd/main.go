package main

import (
	"fmt"
	"goRAGnarok/internal"
	"goRAGnarok/pkg"
	"net/http"
	"os"
)

func main() {
	// Load environment variables from .env file
	if err := pkg.LoadEnv(".env"); err != nil {
		fmt.Println("Warning: .env file not loaded:", err)
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("Warning: OPENAI_API_KEY is not set")
	} else {
		fmt.Printf("apiKey=%s\n", apiKey)
	}

	http.HandleFunc("/health", internal.HealthCheckHandler)
	http.HandleFunc("/v1/generate", internal.PostHandler)
	fmt.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Server error:", err)
	}
}
