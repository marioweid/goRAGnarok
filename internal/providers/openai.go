package providers

import (
	"fmt"
)

type OpenAiProvider struct {
	BaseURL string
	ApiKey  string
}

func (openAiProvider *OpenAiProvider) Generate() string {
	// payload := map[string]any{
	// 		"model": req.Model,
	// 		"input": req.Input,
	// 	}
	// 	payloadBytes, _ := json.Marshal(payload)
	// 	client := &http.Client{}
	// 	url := s.BaseURL + "/responses"
	// 	openaiReq, err := http.NewRequest("POST", url, bytes.NewReader(payloadBytes))
	// 	if err != nil {
	// 		http.Error(w, "Failed to create request", http.StatusInternalServerError)
	// 		return
	// 	}
	// 	openaiReq.Header.Set("Content-Type", "application/json")
	// 	openaiReq.Header.Set("Authorization", "Bearer "+s.APIKey)

	// 	resp, err := client.Do(openaiReq)
	// 	if err != nil {
	// 		http.Error(w, "Failed to call OpenAI API", http.StatusInternalServerError)
	// 		return
	// 	}
	// 	defer resp.Body.Close()

	// 	w.Header().Set("Content-Type", "application/json")
	// 	w.WriteHeader(resp.StatusCode)
	// 	io.Copy(w, resp.Body)
	logMessage := fmt.Sprintf("apiKey: %s, baseUrl %s", openAiProvider.ApiKey, openAiProvider.BaseURL)
	fmt.Println(logMessage)
	return "OpenAiProvider"
}
