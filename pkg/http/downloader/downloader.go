package downloader

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// Downloader handles file downloads with range support
type Downloader struct {
	client       *http.Client
	progressChan chan<- Progress
	maxRetries   int
	initialChunk int64
}

// Progress represents download progress
type Progress struct {
	TotalSize   int64
	Downloaded  int64
	Percentage  float64
	Speed       float64
	TimeElapsed time.Duration
}

// New creates a new Downloader instance
func New(client *http.Client, progressChan chan<- Progress, maxRetries int, initialChunk int64) *Downloader {
	if client == nil {
		client = http.DefaultClient
	}
	return &Downloader{
		client:       client,
		progressChan: progressChan,
		maxRetries:   maxRetries,
		initialChunk: initialChunk,
	}
}

// Download downloads a file from URL to destination path
func (d *Downloader) Download(ctx context.Context, url, dest string) error {
	// Check if server supports range requests
	supportsRange, size, err := d.CheckRangeSupport(ctx, url)
	if err != nil {
		return fmt.Errorf("error checking range support: %w", err)
	}

	if supportsRange {
		return d.downloadWithRanges(ctx, url, dest, size)
	}
	return d.downloadFull(ctx, url, dest)
}

// CheckRangeSupport checks if server supports range requests and returns file size
func (d *Downloader) CheckRangeSupport(ctx context.Context, url string) (bool, int64, error) {
	req, err := http.NewRequestWithContext(ctx, "HEAD", url, nil)
	if err != nil {
		return false, 0, err
	}

	resp, err := d.client.Do(req)
	if err != nil {
		return false, 0, err
	}
	defer resp.Body.Close()

	acceptsRange := resp.Header.Get("Accept-Ranges") == "bytes"
	contentLength := resp.ContentLength

	return acceptsRange, contentLength, nil
}

// downloadWithRanges downloads file using range requests
func (d *Downloader) downloadWithRanges(ctx context.Context, url, dest string, size int64) error {
	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}
	tempPath := dest + ".tmp"
	tempFile, err := os.OpenFile(tempPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer os.Remove(tempFile.Name())

	// Initialize download state
	var (
		downloaded int64
		chunkSize  = d.initialChunk
		startTime  = time.Now()
		ticker     = time.NewTicker(500 * time.Millisecond)
	)
	defer ticker.Stop()

	// Create progress reporting goroutine
	if d.progressChan != nil {
		go func() {
			for range ticker.C {
				elapsed := time.Since(startTime)
				speed := calculateSpeed(downloaded, elapsed)
				d.progressChan <- Progress{
					TotalSize:   size,
					Downloaded:  downloaded,
					Percentage:  float64(downloaded) / float64(size) * 100,
					Speed:       speed,
					TimeElapsed: elapsed,
				}
			}
		}()
	}

	// Download in chunks
	for downloaded < size {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		end := downloaded + chunkSize - 1
		if end >= size {
			end = size - 1
		}

		err := withRetry(ctx, d.maxRetries, func() error {
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				return err
			}
			req = req.WithContext(ctx)
			req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", downloaded, end))

			resp, err := d.client.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusPartialContent {
				return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
			}

			n, err := io.Copy(tempFile, resp.Body)
			if err != nil {
				return err
			}

			downloaded += n
			return nil
		})

		if err != nil {
			return fmt.Errorf("download failed at byte %d: %w", downloaded, err)
		}

		// Adjust chunk size based on download speed
		if d.progressChan != nil {
			elapsed := time.Since(startTime)
			speed := calculateSpeed(downloaded, elapsed)
			chunkSize = adjustChunkSize(chunkSize, speed)
		}
	}

	// Final progress update
	if d.progressChan != nil {
		elapsed := time.Since(startTime)
		d.progressChan <- Progress{
			TotalSize:   size,
			Downloaded:  downloaded,
			Percentage:  100,
			Speed:       calculateSpeed(downloaded, elapsed),
			TimeElapsed: elapsed,
		}
	}

	return os.Rename(tempFile.Name(), dest)
}

// downloadFull downloads entire file in one request
func (d *Downloader) downloadFull(ctx context.Context, url, dest string) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := d.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
