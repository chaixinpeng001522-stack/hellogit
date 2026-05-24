package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHello(t *testing.T) {
	r := newRouter()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
	}
	if got := w.Header().Get(requestIDHeader); got == "" {
		t.Fatalf("missing %s header", requestIDHeader)
	}

	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if body["message"] != "Hello, Gin!" {
		t.Fatalf("unexpected message=%v", body["message"])
	}
	if body["status"] != "success" {
		t.Fatalf("unexpected status=%v", body["status"])
	}
	if _, ok := body["request_id"].(string); !ok {
		t.Fatalf("missing request_id in response")
	}
}

func TestHealthz(t *testing.T) {
	r := newRouter()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
	}
}

func TestNoRoute(t *testing.T) {
	r := newRouter()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/not-exists", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
	}
}