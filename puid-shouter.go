package main

import (
	"fmt"
	"github.com/google/uuid"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	UUIDBufferSize512    = 512
	UUIDBufferSize4096   = 4096
	UUIDBufferSize65535  = 65535
	UUIDBufferSize1MB    = 1024 * 1024
	UUIDBufferSize10MB   = 10 * 1024 * 1024
	UUIDBufferSize100MB  = 100 * 1024 * 1024
)

func main() {
	// Open log file
	logFile, err := os.OpenFile("uuids_per_second.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening log file: %v\n", err)
		os.Exit(1)
	}
	defer logFile.Close()

	// Create buffers
	buffer512 := make([]string, 0, UUIDBufferSize512)
	buffer4096 := make([]string, 0, UUIDBufferSize4096)
	buffer65535 := make([]string, 0, UUIDBufferSize65535)
	buffer1MB := make([]string, 0, UUIDBufferSize1MB)
	buffer10MB := make([]string, 0, UUIDBufferSize10MB)
	buffer100MB := make([]string, 0, UUIDBufferSize100MB)

	// Start a goroutine to dump the buffer contents to stdout and log the statistics periodically
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-time.After(time.Second):
				dumpAndClearBuffer(&buffer512, os.Stdout)
				dumpAndClearBuffer(&buffer4096, os.Stdout)
				dumpAndClearBuffer(&buffer65535, os.Stdout)
				dumpAndClearBuffer(&buffer1MB, os.Stdout)
				dumpAndClearBuffer(&buffer10MB, os.Stdout)
				dumpAndClearBuffer(&buffer100MB, os.Stdout)

				uuidsPerSecond := getUUIDCount(&buffer512) + getUUIDCount(&buffer4096) +
					getUUIDCount(&buffer65535) + getUUIDCount(&buffer1MB) +
					getUUIDCount(&buffer10MB) + getUUIDCount(&buffer100MB)
				bytesPerSecond := getByteCount(&buffer512) + getByteCount(&buffer4096) +
					getByteCount(&buffer65535) + getByteCount(&buffer1MB) +
					getByteCount(&buffer10MB) + getByteCount(&buffer100MB)
				logStats(logFile, uuidsPerSecond, bytesPerSecond)
			}
		}
	}()

	// Loop forever, generating UUIDs
	for {
		// Generate a UUID
		uuid := uuid.New().String()

		// Append the UUID to the appropriate buffer
		appendUUIDToBuffer(uuid, &buffer512)
		appendUUIDToBuffer(uuid, &buffer4096)
		appendUUIDToBuffer(uuid, &buffer65535)
		appendUUIDToBuffer(uuid, &buffer1MB)
		appendUUIDToBuffer(uuid, &buffer10MB)
		appendUUIDToBuffer(uuid, &buffer100MB)

		// Print the UUID to stdout and stderr
		fmt.Println(uuid)
		fmt.Fprintln(os.Stderr, uuid)

		// Check if any buffer is full
		if isBufferFull(&buffer512) {
			dumpBuffer(&buffer512, os.Stdout)
		}
		if isBufferFull(&buffer4096) {
			dumpBuffer(&buffer4096, os.Stdout)
		}
		if isBufferFull(&buffer65535) {
			dumpBuffer(&buffer65535, os.Stdout)
		}
		if isBufferFull(&buffer1MB) {
			dumpBuffer(&buffer1MB, os.Stdout)
		}
		if isBufferFull(&buffer10MB) {
			dumpBuffer(&buffer10MB, os.Stdout)
		}
		if isBufferFull(&buffer100MB) {
			dumpBuffer(&buffer100MB, os.Stdout)
		}
	}
}

// appendUUIDToBuffer appends the UUID to the appropriate buffer.
func appendUUIDToBuffer(uuid string, buffer *[]string) {
	*buffer = append(*buffer, uuid)
}

// isBufferFull checks if the buffer is full based on its capacity.
func isBufferFull(buffer *[]string) bool {
	return len(*buffer) == cap(*buffer)
}

// dumpBuffer dumps the buffer contents to the provided writer.
func dumpBuffer(buffer *[]string, writer io.Writer) {
	uuidBlob := strings.Join(*buffer, "")
	fmt.Fprintln(writer, uuidBlob)
}

// dumpAndClearBuffer dumps the buffer contents to the provided writer and clears the buffer.
func dumpAndClearBuffer(buffer *[]string, writer io.Writer) {
	dumpBuffer(buffer, writer)
	*buffer = (*buffer)[:0]
}

// getUUIDCount returns the number of UUIDs in the buffer.
func getUUIDCount(buffer *[]string) int {
	return len(*buffer)
}

// getByteCount returns the total number of bytes occupied by the UUIDs in the buffer.
func getByteCount(buffer *[]string) int {
	return len(*buffer) * 36
}

// logStats logs the number of UUIDs and bytes generated per second.
func logStats(logFile *os.File, uuidsPerSecond int, bytesPerSecond int) {
	fmt.Fprintf(logFile, "%s UUIDs per second, %s MB/s\n",
		formatWithCommas(uuidsPerSecond),
		formatBytes(float64(bytesPerSecond)))
}

// formatWithCommas formats an integer with commas as thousands separators.
func formatWithCommas(i int) string {
	s := fmt.Sprintf("%d", i)
	groups := make([]string, 0, len(s)/3+1)

	for len(s) > 0 {
		groupSize := len(s) % 3
		if groupSize == 0 {
			groupSize = 3
		}
		groups = append(groups, s[:groupSize])
		s = s[groupSize:]
	}

	return strings.Join(groups, ",")
}

// formatBytes formats a number of bytes in human-readable format.
func formatBytes(bytes float64) string {
	suffixes := []string{"B", "KB", "MB", "GB", "TB"}

	i := 0
	for bytes >= 1024 && i < len(suffixes)-1 {
		bytes /= 1024
		i++
	}

	return fmt.Sprintf("%.2f %s", bytes, suffixes[i])
}

