package downloader

import (
	"fmt"
	"io"
	"strings"
	"time"
)

const (
	progressWidth = 50
)

// DisplayProgress renders a progress bar to the given writer
func DisplayProgress(w io.Writer, p Progress) {
	if p.TotalSize <= 0 {
		fmt.Fprintf(w, "\rDownloading... [unknown size]")
		return
	}

	// Calculate progress bar
	completed := int(float64(progressWidth) * p.Percentage / 100)
	remaining := progressWidth - completed

	// Format speed
	speed := formatBytes(int64(p.Speed)) + "/s"

	// Format time remaining
	var timeRemaining string
	if p.Speed > 0 {
		remainingBytes := p.TotalSize - p.Downloaded
		secs := float64(remainingBytes) / p.Speed
		timeRemaining = formatDuration(time.Duration(secs) * time.Second)
	} else {
		timeRemaining = "unknown"
	}

	// Build progress string
	progressStr := fmt.Sprintf(
		"\r%s [%s%s] %s %s %s remaining",
		formatBytes(p.Downloaded),
		strings.Repeat("=", completed),
		strings.Repeat(" ", remaining),
		fmt.Sprintf("%.1f%%", p.Percentage),
		speed,
		timeRemaining,
	)

	fmt.Fprint(w, progressStr)
}

// formatBytes formats bytes into human-readable string
func formatBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}

// formatDuration formats duration into human-readable string
func formatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second

	if h > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
	}
	return fmt.Sprintf("%02d:%02d", m, s)
}
