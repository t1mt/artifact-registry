package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.linka.cloud/artifact-registry/pkg/http/downloader"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: downloader <url> <output-file>")
		os.Exit(1)
	}

	url := os.Args[1]
	outputFile := os.Args[2]

	// Create progress channel
	progressChan := make(chan downloader.Progress)
	defer close(progressChan)

	// Start progress reporter
	go func() {
		for p := range progressChan {
			downloader.DisplayProgress(os.Stdout, p)
		}
		time.Sleep(500 * time.Millisecond)
	}()

	// Create downloader with:
	// - Default HTTP client
	// - Progress channel
	// - Max 3 retries
	// - Initial chunk size of 1MB
	dl := downloader.New(nil, progressChan, 3, 1024*1024)

	// Set timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	// Start download
	fmt.Printf("Downloading %s to %s...\n", url, outputFile)
	err := dl.Download(ctx, url, outputFile)
	if err != nil {
		fmt.Printf("\nDownload failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("\nDownload completed successfully!")
}
