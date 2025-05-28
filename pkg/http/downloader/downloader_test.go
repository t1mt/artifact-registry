package downloader

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestDownloader(t *testing.T) {
	tests := []struct {
		name           string
		supportsRange  bool
		fileSize       int
		expectedChunks int
	}{
		{"Small file without range", false, 1024, 1},
		{"Large file with range", true, 10 * 1024 * 1024, 10},
		{"Small file with range", true, 1024, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.supportsRange {
					w.Header().Set("Accept-Ranges", "bytes")
				}

				if r.Method == "HEAD" {
					w.Header().Set("Content-Length", fmt.Sprintf("%d", tt.fileSize))
					return
				}

				// Simulate file content
				data := make([]byte, tt.fileSize)
				for i := range data {
					data[i] = byte(i % 256)
				}

				if tt.supportsRange && r.Header.Get("Range") != "" {
					// Handle range request
					// (simplified for test - actual implementation would parse range header)
					rangeHeader := r.Header.Get("Range")
					parts := strings.Split(rangeHeader, "=")[1]
					ranges := strings.Split(parts, "-")
					start, _ := strconv.ParseInt(ranges[0], 10, 64)
					end, _ := strconv.ParseInt(ranges[1], 10, 64)

					reader := bytes.NewReader(data)
					reader.Seek(start, io.SeekStart)
					limitedReader := io.LimitReader(reader, end-start+1)
					w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, tt.fileSize))
					w.WriteHeader(http.StatusPartialContent)
					io.Copy(w, limitedReader)
				} else {
					// Full file request
					w.Write(data)
				}
			})

			server := httptest.NewServer(handler)
			defer server.Close()

			// Setup downloader
			tempDir, err := os.MkdirTemp("", "downloader-test")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tempDir)

			dest := filepath.Join(tempDir, "testfile")
			progressChan := make(chan Progress, 10)
			defer close(progressChan)

			dl := New(nil, progressChan, 3, 1024*1024) // 1MB initial chunk

			// Run download
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			err = dl.Download(ctx, server.URL, dest)
			if err != nil {
				t.Fatalf("Download failed: %v", err)
			}

			// Verify file
			fileInfo, err := os.Stat(dest)
			if err != nil {
				t.Fatalf("Failed to stat downloaded file: %v", err)
			}

			if fileInfo.Size() != int64(tt.fileSize) {
				t.Errorf("File size mismatch: got %d, want %d", fileInfo.Size(), int64(tt.fileSize))
			}

			// TODO: Add more comprehensive file content verification
			// TODO: Test progress reporting
			// TODO: Test retry mechanism
			// TODO: Test dynamic chunk adjustment
		})
	}
}

func TestRangeSupportCheck(t *testing.T) {
	// Test server that doesn't support ranges
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "HEAD" {
			w.Header().Set("Content-Length", "1024")
			return
		}
		w.Write(make([]byte, 1024))
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	dl := New(nil, nil, 3, 1024)
	ctx := context.Background()

	supportsRange, size, err := dl.CheckRangeSupport(ctx, server.URL)
	if err != nil {
		t.Fatalf("checkRangeSupport failed: %v", err)
	}

	if supportsRange {
		t.Error("Expected server to not support ranges")
	}

	if size != 1024 {
		t.Errorf("Expected size 1024, got %d", size)
	}
}
