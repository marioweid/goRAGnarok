package providers

import (
	"fmt"
)

type OllamaProvider struct {
	BaseURL string
}

func (openAiProvider *OllamaProvider) Generate() string {
	logMessage := fmt.Sprintf("baseUrl %s", openAiProvider.BaseURL)
	fmt.Println(logMessage)
	return "OllamaProvider"
}
