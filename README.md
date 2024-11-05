# API Usage

## Endpoints

### POST /convert
Converts a WAV audio stream to FLAC format.

### WebSocket /ws
Streams FLAC audio data back to the client.

# Setup Instructions

## Running Locally
1. Clone the repository.
2. Run `go mod tidy` to install dependencies.
3. Run `go run main.go` to start the server.

# Testing Strategy

## Unit Tests
Run unit tests using `go test ./tests/...`.

## Integration Tests
Run integration tests using `go test ./tests/... -tags=integration`.