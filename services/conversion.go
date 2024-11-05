// services/conversion.go

package services

import (
	"bytes"
	"fmt"
	"os/exec"
)

func WavToFlac(wavData []byte) ([]byte, error) {
	wavBuf := bytes.NewBuffer(wavData)

	// Use ffmpeg to convert WAV to FLAC
	cmd := exec.Command("ffmpeg", "-i", "pipe:0", "-f", "flac", "pipe:1")
	cmd.Stdin = wavBuf
	var flacBuf bytes.Buffer
	cmd.Stdout = &flacBuf

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("ffmpeg error: %v", err)
	}

	return flacBuf.Bytes(), nil
}
