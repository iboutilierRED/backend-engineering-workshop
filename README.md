# Wikipedia Stream Stats Workshop

This Go application connects to the Wikipedia recent changes stream, collects statistics about edits, and exposes them via an HTTP API.

## Features
- Consumes live Wikipedia edit events using a streaming API
- Tracks:
  - Total messages
  - Unique users
  - Unique URIs
  - Bot vs. non-bot edits
- Exposes stats at `/stats` endpoint (JSON)
- Thread-safe, channel-based architecture
- Comprehensive unit and integration tests

## Architecture
- **StatsCollector**: Encapsulates all stats and synchronization
- **Channel-based pipeline**: Wikipedia events are parsed and sent to a buffered channel; a single goroutine updates stats
- **HTTP handler**: Safely snapshots stats for API responses

## Usage

1. **Build the wikiapp binary:**
  ```bash
  go build -o ch-1/wikiapp ./ch-1/cmd/main.go
  ```

2. **Run the wikiapp binary:**
  ```bash
  ./ch-1/wikiapp
  ```
  The server listens on port 7000 by default.

3. **Query stats:**
  ```bash
  curl http://localhost:7000/stats
  ```

...existing code...

## Testing
- Run all tests (with race detector):
  ```bash
  go test -race -v ./ch-1/cmd/
  ```
- Includes unit tests for JSON parsing, stats collection, HTTP handler, and integration tests for channel-based processing.

## Project Structure
```
workshop/
├── ch-1/
│   └── cmd/
│       ├── main.go        # Application entrypoint
│       └── main_test.go   # Tests and benchmarks
├── go.mod
└── README.md
```

## License
MIT
