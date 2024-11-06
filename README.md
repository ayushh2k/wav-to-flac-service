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


sequenceDiagram
    participant Client
    participant WebSocket
    participant Server
    participant Converter
    
    Client->>WebSocket: Connect
    WebSocket->>Server: Establish Connection
    
    loop Streaming Process
        Client->>WebSocket: Send WAV Chunk
        WebSocket->>Server: Receive WAV Chunk
        Server->>Converter: Process WAV Chunk
        Converter->>Server: Return FLAC Chunk
        Server->>WebSocket: Send FLAC Chunk
        WebSocket->>Client: Receive FLAC Chunk
    end
    
    Client->>WebSocket: Close Connection
    WebSocket->>Server: End Stream