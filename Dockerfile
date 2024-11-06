# Use an official Go runtime as a parent image
FROM golang:1.23-alpine

# Install ffmpeg
RUN apk add --no-cache ffmpeg

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o wav-to-flac-service ./cmd/main.go

EXPOSE 8080

# Command to run the executable
CMD ["./wav-to-flac-service"]