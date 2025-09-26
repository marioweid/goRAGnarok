package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"goRAGnarok/internal/interfaces"
	"net/http"
)

type OpenAiProvider struct {
	BaseURL string
	ApiKey  string
}

func (openAiProvider *OpenAiProvider) Generate(request interfaces.GenerateRequest) (string, error) {
	payload := map[string]any{
		"model": request.Model,
		"input": request.Input,
	}
	payloadBytes, _ := json.Marshal(payload)
	client := &http.Client{}
	fmt.Println(client) // TODO USELESS
	url := openAiProvider.BaseURL + "/responses"
	openaiReq, err := http.NewRequest("POST", url, bytes.NewReader(payloadBytes))
	if err != nil {
		return "", err
	}
	openaiReq.Header.Set("Content-Type", "application/json")
	openaiReq.Header.Set("Authorization", "Bearer "+openAiProvider.ApiKey)

	// TODO COMMENT IN WITH INTERNET
	// resp, err := client.Do(openaiReq)
	// if err != nil {
	// 	return "", err
	// }
	// defer resp.Body.Close()

	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(resp.StatusCode)
	// io.Copy(w, resp.Body)
	logMessage := fmt.Sprintf("apiKey: %s, baseUrl %s", openAiProvider.ApiKey, openAiProvider.BaseURL)
	fmt.Println(logMessage)
	return "OpenAiProvider", nil
}
