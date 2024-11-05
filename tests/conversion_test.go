// tests/conversion_test.go

package tests

import (
	"os"
	"testing"

	"wav-to-flac-service/services"
)

func TestWavToFlac(t *testing.T) {
	wavData, err := os.ReadFile("../harvard.wav")
	if err != nil {
		t.Fatalf("Failed to read test WAV file: %v", err)
	}

	flacData, err := services.WavToFlac(wavData)
	if err != nil {
		t.Fatalf("WavToFlac failed: %v", err)
	}

	if len(flacData) == 0 {
		t.Error("FLAC data is empty")
	}
}
