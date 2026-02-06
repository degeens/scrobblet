package csv

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/degeens/scrobblet/internal/common"
)

type Client struct {
	filePath string
	mu       sync.Mutex
}

func NewClient(filePath string) *Client {
	return &Client{
		filePath: filePath,
	}
}

func (c *Client) WriteScrobble(trackedTrack *common.TrackedTrack) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if file exists and if it's empty
	needsHeader := false
	fileInfo, err := os.Stat(c.filePath)
	if os.IsNotExist(err) {
		needsHeader = true
		// Create parent directories if they don't exist
		dir := filepath.Dir(c.filePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to stat CSV file: %w", err)
	} else if fileInfo.Size() == 0 {
		needsHeader = true
	}

	// Open file in append mode, create if doesn't exist
	file, err := os.OpenFile(c.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	writer := csv.NewWriter(file)

	// Write header if needed
	if needsHeader {
		header := []string{"Artists", "Title", "Album", "Track Number", "Duration", "Started At", "Ended At"}
		if err := writer.Write(header); err != nil {
			return fmt.Errorf("failed to write CSV header: %w", err)
		}
	}

	// Prepare record
	artists := strings.Join(trackedTrack.Track.Artists, ", ")
	title := trackedTrack.Track.Title
	album := trackedTrack.Track.Album
	trackNumber := strconv.Itoa(trackedTrack.Track.TrackNumber)
	duration := strconv.Itoa(int(trackedTrack.Track.Duration / time.Second))
	startedAt := trackedTrack.StartedAt.Format(time.RFC3339)
	endedAt := trackedTrack.LastUpdatedAt.Format(time.RFC3339)

	record := []string{
		artists,
		title,
		album,
		trackNumber,
		duration,
		startedAt,
		endedAt,
	}

	// Write record
	if err := writer.Write(record); err != nil {
		return fmt.Errorf("failed to write CSV record: %w", err)
	}

	// Flush and check for errors
	writer.Flush()
	if err := writer.Error(); err != nil {
		return fmt.Errorf("failed to flush CSV writer: %w", err)
	}

	return nil
}
