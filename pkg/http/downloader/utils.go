package downloader

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

// adjustChunkSize dynamically adjusts the chunk size based on download speed
func adjustChunkSize(currentChunk int64, speed float64) int64 {
	const (
		minChunkSize = 64 * 1024        // 64KB
		maxChunkSize = 16 * 1024 * 1024 // 16MB
		fastSpeed    = 5 * 1024 * 1024  // 5MB/s
	)

	if speed > fastSpeed {
		// Increase chunk size for fast connections
		newChunk := currentChunk * 2
		if newChunk > maxChunkSize {
			return maxChunkSize
		}
		return newChunk
	} else if speed < fastSpeed/2 {
		// Decrease chunk size for slow connections
		newChunk := currentChunk / 2
		if newChunk < minChunkSize {
			return minChunkSize
		}
		return newChunk
	}
	return currentChunk
}

// withRetry executes a function with retry logic
func withRetry(ctx context.Context, maxRetries int, fn func() error) error {
	var lastErr error
	baseDelay := time.Second

	for i := 0; i < maxRetries; i++ {
		err := fn()
		if err == nil {
			return nil
		}

		// Check if context is done
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		lastErr = err
		delay := time.Duration(i*i) * baseDelay
		time.Sleep(delay)
	}

	return fmt.Errorf("after %d retries: %w", maxRetries, lastErr)
}

// verifyFile checks if the downloaded file matches the expected checksum
func verifyFile(filePath string, expectedChecksum string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return err
	}

	actualChecksum := fmt.Sprintf("%x", hash.Sum(nil))
	if actualChecksum != expectedChecksum {
		return errors.New("checksum mismatch")
	}

	return nil
}

// calculateSpeed calculates download speed in bytes per second
func calculateSpeed(downloaded int64, elapsed time.Duration) float64 {
	if elapsed <= 0 {
		return 0
	}
	return float64(downloaded) / elapsed.Seconds()
}

// getFileSize returns the size of a file
func getFileSize(filePath string) (int64, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return 0, err
	}

	return stat.Size(), nil
}
