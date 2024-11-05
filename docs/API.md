# API Documentation

## WebSocket Endpoint

### `/stream`

**Method:** GET

**Description:**
This endpoint establishes a WebSocket connection for real-time streaming of WAV audio data to be converted to FLAC format.

**Usage:**
1. Connect to the WebSocket endpoint using a WebSocket client.
2. Send WAV audio data in binary format.
3. Receive the converted FLAC audio data in binary format.

**Example:**
```bash
# Using wscat (WebSocket client)
wscat -c ws://localhost:8080/stream
```
Request:

    Binary Data: WAV audio data.

Response:

    Binary Data: FLAC audio data.