package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestGenerateSuccess(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Error("missing API authentication")
		}
		var request struct {
			Model    string              `json:"model"`
			Messages []map[string]string `json:"messages"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatal(err)
		}
		if request.Model != "gpt-4.1-mini" || !strings.Contains(request.Messages[0]["content"], "Dog") {
			t.Error("invalid upstream request")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"Captain Paws"}}]}`))
	}))
	defer upstream.Close()
	app := newHandler(&http.Client{Timeout: time.Second}, "test-key", "gpt-4.1-mini", upstream.URL)
	r := httptest.NewRequest(http.MethodPost, "/api/generate", strings.NewReader(`{"prompt":"Dog"}`))
	w := httptest.NewRecorder()
	app.ServeHTTP(w, r)
	if w.Code != 200 || !strings.Contains(w.Body.String(), "Captain Paws") {
		t.Fatalf("got %d %s", w.Code, w.Body.String())
	}
}

func TestInvalidRequestsDoNotCallOpenAI(t *testing.T) {
	calls := 0
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls++ }))
	defer upstream.Close()
	app := newHandler(upstream.Client(), "test-key", "gpt-4.1-mini", upstream.URL)
	cases := []struct {
		method string
		body   string
		status int
	}{
		{"GET", `{}`, 405},
		{"POST", `{`, 400},
		{"POST", `{"prompt":"  "}`, 400},
		{"POST", `{"prompt":"Dog","temperature":3}`, 400},
		{"POST", `{"prompt":"Dog"} {}`, 400},
	}
	for _, tc := range cases {
		w := httptest.NewRecorder()
		app.ServeHTTP(w, httptest.NewRequest(tc.method, "/api/generate", strings.NewReader(tc.body)))
		if w.Code != tc.status {
			t.Errorf("%s %s: got %d want %d", tc.method, tc.body, w.Code, tc.status)
		}
	}
	if calls != 0 {
		t.Fatalf("invalid request triggered %d upstream calls", calls)
	}
}

func TestUpstreamError(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(429) }))
	defer upstream.Close()
	app := newHandler(upstream.Client(), "test-key", "gpt-4.1-mini", upstream.URL)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, httptest.NewRequest("POST", "/api/generate", strings.NewReader(`{"prompt":"Dog"}`)))
	if w.Code != 502 {
		t.Fatalf("got %d", w.Code)
	}
}
