package handlers

import (
	"bytes"
	"goRAGnarok/internal"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheckHandler_OK(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	HealthCheckHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
	if w.Body.String() != "OK" {
		t.Errorf("expected body 'OK', got '%s'", w.Body.String())
	}
}

func TestHealthCheckHandler_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/health", nil)
	w := httptest.NewRecorder()

	HealthCheckHandler(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", w.Code)
	}
}

func TestResponseHandler_MethodNotAllowed(t *testing.T) {
	srv := &internal.Server{APIKey: "dummy"}
	req := httptest.NewRequest(http.MethodGet, "/v1/response", nil)
	w := httptest.NewRecorder()

	handler := ResponseHandler(srv)
	handler(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", w.Code)
	}
}

func TestResponseHandler_InvalidJSON(t *testing.T) {
	srv := &internal.Server{APIKey: "dummy"}
	body := bytes.NewBufferString(`invalid-json`)
	req := httptest.NewRequest(http.MethodPost, "/v1/response", body)
	w := httptest.NewRecorder()

	handler := ResponseHandler(srv)
	handler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}
