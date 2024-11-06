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
	wavData, err := os.ReadFile("../assets/harvard.wav")
	if err != nil {
		t.Fatalf("Failed to read test WAV file: %v", err)
	}

	converter, err := services.NewStreamingConverter()
	if err != nil {
		t.Fatalf("Failed to create streaming converter: %v", err)
	}
	defer converter.Close()

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

	if flacData.Len() == 0 {
		t.Error("FLAC data is empty")
	}

	flacHeader := flacData.Bytes()[:4]
	expectedHeader := []byte("fLaC")
	if !bytes.Equal(flacHeader, expectedHeader) {
		t.Error("Output data does not appear to be in FLAC format")
	}
}

func TestConverterLifecycle(t *testing.T) {
	converter, err := services.NewStreamingConverter()
	if err != nil {
		t.Fatalf("Failed to create converter: %v", err)
	}

	testData := []byte("RIFF....WAVEfmt ")
	err = converter.WriteInput(testData)
	if err != nil {
		t.Errorf("Failed to write test data: %v", err)
	}

	buffer := make([]byte, 1024)
	_, err = converter.ReadConverted(buffer)
	if err != nil && err != io.EOF {
		t.Errorf("Failed to read converted data: %v", err)
	}

	err = converter.Close()
	if err != nil {
		t.Errorf("Failed to close converter: %v", err)
	}
}

func TestConverterErrors(t *testing.T) {
	converter, err := services.NewStreamingConverter()
	if err != nil {
		t.Fatalf("Failed to create converter: %v", err)
	}
	defer converter.Close()

	converter.Close()
	err = converter.WriteInput([]byte("test"))
	if err == nil {
		t.Error("Expected error when writing to closed converter")
	}
}
