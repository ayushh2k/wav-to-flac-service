// services/conversion.go
package services

import (
	"io"
	"os/exec"
	"sync"
)

type StreamingConverter struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser
	mutex  sync.Mutex
}

func NewStreamingConverter() (*StreamingConverter, error) {
	// Configure ffmpeg for streaming conversion
	cmd := exec.Command("ffmpeg",
		"-f", "wav", // Input format
		"-i", "pipe:0", // Read from stdin
		"-f", "flac", // Output format
		"-compression_level", "0", // Fast compression
		"-flush_packets", "1", // Flush packets immediately
		"-fflags", "+nobuffer", // Disable buffering
		"pipe:1", // Write to stdout
	)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		stdin.Close()
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		stdin.Close()
		stdout.Close()
		return nil, err
	}

	return &StreamingConverter{
		cmd:    cmd,
		stdin:  stdin,
		stdout: stdout,
	}, nil
}

func (c *StreamingConverter) WriteInput(data []byte) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	_, err := c.stdin.Write(data)
	return err
}

func (c *StreamingConverter) ReadConverted(buffer []byte) (int, error) {
	return c.stdout.Read(buffer)
}

func (c *StreamingConverter) Close() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.stdin.Close()
	c.stdout.Close()
	return c.cmd.Wait()
}
