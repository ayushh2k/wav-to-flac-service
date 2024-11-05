# Testing Strategy

## Unit Tests

Unit tests are used to validate the conversion logic from WAV to FLAC.

### Running Unit Tests

1. **Navigate to the Project Directory**:
   ```bash
   cd wav-to-flac-service
   ```
2. **Run the tests**:
   ```bash
   go test ./tests/conversion_test.go
   ```

## Integration Tests

Integration tests are used to validate the streaming and WebSocket functionality.

### Running Integration Tests

1. **Navigate to the Project Directory**:
   ```bash
   cd wav-to-flac-service
   ```
2. **Run the tests**:
   ```bash
   go test ./tests/integration_test.go
   ```

## Manual Testing

Manual testing involves invoking the service and sending WAV data to test the streaming and WebSocket functionality.

### Running Manual Testing

1. **Navigate to the Project Directory**:
   ```bash
   cd wav-to-flac-service
   ```
2. **Run the server**:
   ```bash
   go run cmd/main.go
   ```
3. **Connect to the WebSocket Endpoint**:        
   ```bash
   wscat -c ws://localhost:8080/stream
   ```
4. **Send WAV Data**:        
   ```bash
   cat test.wav | wscat -c ws://localhost:8080/stream
   ```
5. **Receive FLAC Data**:        
   ```bash
   cat test.wav | wscat -c ws://localhost:8080/stream > output.flac
   ```

## Error Handling

The system is designed to gracefully handle errors in streaming and conversion processes:

- **WebSocket Errors**: If there is an error in reading or writing WebSocket messages, the connection will be closed and an error message will be logged.
- **Conversion Errors**: If there is an error in the conversion process (e.g., FFmpeg error), an error message will be logged, and the client will receive an error message.

## Logging

Error messages are logged using the utils.LogError function, which logs the context and error details.
