# OpenAI API quickstart (Go)

Minimal pet superhero-name generator, intentionally implemented with the Go standard library only. This eliminates legacy third-party dependency vulnerabilities in the original example while keeping the GET `/` page and POST `/api/generate` response shape (`{"result":"..."}`).

## Setup

Requires Go 1.19+ and an OpenAI API key.

Copy `.env.example` into your preferred local environment setup. The server reads environment variables directly; it does **not** automatically read a `.env` file.

On Linux/macOS:

```sh
export OPENAI_API_KEY="your-key"
export OPENAI_MODEL="gpt-4.1-mini" # optional
go test ./...
go run .
```

On Windows PowerShell:

```powershell
$env:OPENAI_API_KEY="your-key"
$env:OPENAI_MODEL="gpt-4.1-mini"
go test ./...
go run .
```

Open http://localhost:8080.

## Security notes

- No external Go modules. There is no `go.sum` because all imports are from the Go standard library.
- HTTP requests are bounded by size and deadlines; upstream requests have timeouts.
- The model is chosen by the server-side `OPENAI_MODEL` environment variable, not by the browser client, avoiding unauthenticated model-cost escalation.
- Model output is inserted into the page with `textContent` to avoid HTML/script injection.
- The old `text-davinci-003` completions model was replaced with the Chat Completions API using `gpt-4.1-mini` by default. API output style can change accordingly.
- This is a demo, not a hardened public deployment. Add authentication and rate limits before exposing it to the internet.

## API

POST `/api/generate` with JSON body `{"prompt":"Dog"}`. Responds `{"result":"..."}` on success. Optional `temperature` remains supported between 0 and 2. The legacy `model` and `max_tokens` client fields are accepted but ignored to prevent untrusted cost escalation.
