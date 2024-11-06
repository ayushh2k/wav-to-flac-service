// tests/conversion_test.go
package tests

import (
	"bytes"
	"io"
	"os"
	"testing"
	"wav-to-flac-service/services"
)

func TestWavToFlacStreaming(t *testing.T) {
	// Read test WAV file
	wavData, err := os.ReadFile("../assets/harvard.wav")
	if err != nil {
		t.Fatalf("Failed to read test WAV file: %v", err)
	}

	// Create new streaming converter
	converter, err := services.NewStreamingConverter()
	if err != nil {
		t.Fatalf("Failed to create streaming converter: %v", err)
	}
	defer converter.Close()

	// Write WAV data in chunks
	chunkSize := 4096
	for i := 0; i < len(wavData); i += chunkSize {
		end := i + chunkSize
		if end > len(wavData) {
			end = len(wavData)
		}

		err := converter.WriteInput(wavData[i:end])
		if err != nil {
			t.Fatalf("Failed to write WAV chunk: %v", err)
		}
	}

	// Read converted FLAC data
	var flacData bytes.Buffer
	buffer := make([]byte, 4096)
	for {
		n, err := converter.ReadConverted(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("Failed to read converted data: %v", err)
		}
		if n > 0 {
			flacData.Write(buffer[:n])
		}
	}

	// Verify conversion result
	if flacData.Len() == 0 {
		t.Error("FLAC data is empty")
	}

	// Optional: Basic validation of FLAC format
	flacHeader := flacData.Bytes()[:4]
	expectedHeader := []byte("fLaC") // FLAC magic number
	if !bytes.Equal(flacHeader, expectedHeader) {
		t.Error("Output data does not appear to be in FLAC format")
	}
}

// Additional test for converter lifecycle
func TestConverterLifecycle(t *testing.T) {
	// Test creation
	converter, err := services.NewStreamingConverter()
	if err != nil {
		t.Fatalf("Failed to create converter: %v", err)
	}

	// Test basic I/O
	testData := []byte("RIFF....WAVEfmt ") // Basic WAV header simulation
	err = converter.WriteInput(testData)
	if err != nil {
		t.Errorf("Failed to write test data: %v", err)
	}

	// Test reading
	buffer := make([]byte, 1024)
	_, err = converter.ReadConverted(buffer)
	if err != nil && err != io.EOF {
		t.Errorf("Failed to read converted data: %v", err)
	}

	// Test cleanup
	err = converter.Close()
	if err != nil {
		t.Errorf("Failed to close converter: %v", err)
	}
}

// Test error handling
func TestConverterErrors(t *testing.T) {
	converter, err := services.NewStreamingConverter()
	if err != nil {
		t.Fatalf("Failed to create converter: %v", err)
	}
	defer converter.Close()

	// Test writing after close
	converter.Close()
	err = converter.WriteInput([]byte("test"))
	if err == nil {
		t.Error("Expected error when writing to closed converter")
	}
}
