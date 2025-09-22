package main

import (
	"fmt"
	"goRAGnarok/internal/interfaces"
	"goRAGnarok/internal/providers"

	_ "github.com/lib/pq"
)

func providerRouter(provider interfaces.Provider) string {
	providerType := provider.Generate()
	fmt.Println(providerType)
	return providerType
}

// todo implement lookup table
// get key from body
// depending on the key, lookup the handler
// call generate of the given handler

// type LookupTable struct {
// 	providers map[string]interface
// }

// var providers = []interfaces.Provider{
// 	&providers.OpenAiProvider{BaseURL: "my_base_url", ApiKey: "my_api_key"},
// 	&providers.OllamaProvider{BaseURL: "my_base_url"},
// }

func main() {
	providerLookup := make(map[string]interfaces.Provider)
	// openAi := providers.OpenAiProvider{BaseURL: "my_base_url", ApiKey: "my_api_key"}
	// p = Person{Name: "Ada", Age: 28}

	providerLookup["openai"] = &providers.OpenAiProvider{BaseURL: "my_base_url", ApiKey: "my_api_key"}
	providerLookup["ollama"] = &providers.OllamaProvider{BaseURL: "my_base_url"}
	providerLookup["openai"].Generate()
	providerLookup["ollama"].Generate()

	// use lookup table
	selectedProvider, ok := providerLookup["openai"]
	fmt.Printf("Value for 'grape': %v, OK: %v\n", selectedProvider, ok)
	if ok {
		providerRouter(selectedProvider)
	} else {
		fmt.Printf("Provider not found!")
	}
	providerRouter(providerLookup["ollama"])
}
