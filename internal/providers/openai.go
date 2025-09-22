package providers

import (
	"fmt"
)

type OpenAiProvider struct {
	BaseURL string
	ApiKey  string
}

func (openAiProvider *OpenAiProvider) Generate() string {
	logMessage := fmt.Sprintf("apiKey: %s, baseUrl %s", openAiProvider.ApiKey, openAiProvider.BaseURL)
	fmt.Println(logMessage)
	return "OpenAiProvider"
}
