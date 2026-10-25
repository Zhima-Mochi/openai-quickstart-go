package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const openAIURL = "https://api.openai.com/v1/chat/completions"

type generateRequest struct {
	Prompt      string   `json:"prompt"`
	Model       string   `json:"model,omitempty"`
	Temperature *float64 `json:"temperature,omitempty"`
	MaxTokens   int      `json:"max_tokens,omitempty"`
}

type chatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func newHandler(client *http.Client, apiKey, model, endpoint string) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		http.ServeFile(w, r, "index.html")
	})
	mux.HandleFunc("/api/generate", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if r.Method != http.MethodPost {
			jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, 8192)
		var input generateRequest
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&input); err != nil {
			jsonError(w, http.StatusBadRequest, "invalid request")
			return
		}
		var extra any
		if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
			jsonError(w, http.StatusBadRequest, "invalid request")
			return
		}
		animal := strings.TrimSpace(input.Prompt)
		if animal == "" || len(animal) > 200 {
			jsonError(w, http.StatusBadRequest, "please enter a valid animal (1–200 bytes)")
			return
		}
		temperature := 0.6
		if input.Temperature != nil {
			temperature = *input.Temperature
			if temperature < 0 || temperature > 2 {
				jsonError(w, http.StatusBadRequest, "temperature must be between 0 and 2")
				return
			}
		}
		// Keep model selection server-side: this public demo must not let visitors choose costly models.
		requestBody := struct {
			Model       string              `json:"model"`
			Messages    []map[string]string `json:"messages"`
			Temperature float64             `json:"temperature"`
			MaxTokens   int                 `json:"max_tokens"`
		}{
			Model:       model,
			Messages:    []map[string]string{{"role": "user", "content": fmt.Sprintf("Suggest three names for an animal that is a superhero.\nAnimal: Cat\nNames: Captain Sharpclaw, Agent Fluffball, The Incredible Feline\nAnimal: Dog\nNames: Ruff the Protector, Wonder Canine, Sir Barks-a-Lot\nAnimal: %s\nNames:", animal)}},
			Temperature: temperature,
			MaxTokens:   100,
		}
		payload, err := json.Marshal(requestBody)
		if err != nil {
			jsonError(w, http.StatusInternalServerError, "internal error")
			return
		}
		upstream, err := http.NewRequestWithContext(r.Context(), http.MethodPost, endpoint, bytes.NewReader(payload))
		if err != nil {
			jsonError(w, http.StatusInternalServerError, "internal error")
			return
		}
		upstream.Header.Set("Authorization", "Bearer "+apiKey)
		upstream.Header.Set("Content-Type", "application/json")
		response, err := client.Do(upstream)
		if err != nil {
			log.Printf("OpenAI request failed: %v", err)
			jsonError(w, http.StatusBadGateway, "text generation unavailable")
			return
		}
		defer response.Body.Close()
		if response.StatusCode != http.StatusOK {
			log.Printf("OpenAI request failed with HTTP %d", response.StatusCode)
			jsonError(w, http.StatusBadGateway, "text generation unavailable")
			return
		}
		var result chatResponse
		if err := json.NewDecoder(io.LimitReader(response.Body, 1<<20)).Decode(&result); err != nil || len(result.Choices) == 0 || result.Choices[0].Message.Content == "" {
			jsonError(w, http.StatusBadGateway, "invalid generation response")
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]string{"result": result.Choices[0].Message.Content})
	})
	return mux
}

func jsonError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func main() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY must be set")
	}
	model := os.Getenv("OPENAI_MODEL")
	if model == "" {
		model = "gpt-4.1-mini"
	}
	server := &http.Server{
		Addr:              ":8080",
		Handler:           newHandler(&http.Client{Timeout: 25 * time.Second}, apiKey, model, openAIURL),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	log.Printf("Listening on http://localhost%s", server.Addr)
	log.Fatal(server.ListenAndServe())
}
