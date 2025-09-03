package internal

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheckHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	HealthCheckHandler(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
	body := w.Body.String()
	if body != "OK" {
		t.Errorf("expected body 'OK', got '%s'", body)
	}
}

func TestPostHandler_ValidJSON(t *testing.T) {
	payload := map[string]interface{}{
		"model": "gpt-4.1-nano",
		"messages": []map[string]string{
			{"role": "user", "content": "Hello!"},
			{"role": "assistant", "content": "Hi, how can I help you?"},
		},
		"temperature": 0.4,
	}
	b, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/v1/post", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	PostHandler(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Errorf("failed to decode response: %v", err)
	}
	if _, ok := result["received"]; !ok {
		t.Errorf("expected 'received' field in response")
	}
}

func TestPostHandler_InvalidJSON(t *testing.T) {
	b := []byte(`{"model": "gpt-4-1106-preview", "messages": [}`) // malformed JSON
	req := httptest.NewRequest(http.MethodPost, "/v1/post", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	PostHandler(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}
