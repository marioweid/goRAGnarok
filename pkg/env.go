package pkg

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

// LoadEnv loads environment variables from a .env file
func LoadEnv(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || len(strings.TrimSpace(line)) == 0 {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			os.Setenv(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
		}
	}
	return scanner.Err()
}

// Config holds all the configuration for the application
type Config struct {
	OpenAIAPIKey         string
	OpenAIBaseURL        string
	OpenAIEmbeddingModel string
	OllamaBaseURL        string
	OllamaEmbeddingModel string
	DatabaseURL          string
}

// NewConfig creates a new Config struct and populates it with environment variables.
func NewConfig() (*Config, error) {
	if err := LoadEnv(".env"); err != nil {
		log.Println("Warning: .env file not loaded:", err)
	}

	cfg := &Config{
		OpenAIAPIKey:         os.Getenv("OPENAI_API_KEY"),
		OpenAIBaseURL:        os.Getenv("OPENAI_BASE_URL"),
		OpenAIEmbeddingModel: os.Getenv("OPENAI_EMBEDDING_MODEL"),
		OllamaBaseURL:        os.Getenv("OLLAMA_BASE_URL"),
		OllamaEmbeddingModel: os.Getenv("OLLAMA_EMBEDDING_MODEL"),
		DatabaseURL:          os.Getenv("DATABASE_URL"),
	}

	// Set defaults for optional values
	if cfg.OpenAIBaseURL == "" {
		cfg.OpenAIBaseURL = "https://api.openai.com/v1"
	}
	if cfg.OpenAIEmbeddingModel == "" {
		cfg.OpenAIEmbeddingModel = "text-embedding-3-small"
	}
	if cfg.OllamaBaseURL == "" {
		cfg.OllamaBaseURL = "http://localhost:11434"
	}
	if cfg.OllamaEmbeddingModel == "" {
		cfg.OllamaEmbeddingModel = "all-minilm"
	}

	// Validate required fields
	if cfg.OpenAIAPIKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY is not set")
	}
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set")
	}

	return cfg, nil
}
