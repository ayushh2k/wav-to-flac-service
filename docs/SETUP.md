# Setup Instructions

## Prerequisites

1. **Go Programming Language**: Ensure Go is installed on your system. You can download it from [here](https://golang.org/dl/).
2. **FFmpeg**: Ensure FFmpeg is installed on your system. You can download it from [here](https://ffmpeg.org/download.html).
3. **WebSocket Client**: For testing purposes, you can use `wscat`, a command-line WebSocket client. Install it using npm:

   ```bash
   npm install -g wscat
   ```

### Running the Service Locally

##### 1. Clone the Repository

    ```bash

    git clone https://github.com/ayushh2k/wav-to-flac-service.git
    cd wav-to-flac-service
    ```

##### 2. Install Dependencies

    ```bash
    go mod download
    ```

##### 3. Run the Server

    ```bash
    go run cmd/main.go
    ```
    The server will start on localhost:8080.

### Deploying the Service

##### 1. Build the Binary

    ```bash
    go build -o wav-to-flac cmd/main.go
    ```

##### 2. Run the Binary

    ```bash
    ./wav-to-flac
    ```

##### 3. Deploy to a Server

    Copy the binary and necessary files to your server.
    Ensure FFmpeg is installed on the server.
    Run the binary on the server.

### Configuration

Port: The server runs on port 8080 by default. You can change this by modifying the r.Run(":8080") line in cmd/main.go.
