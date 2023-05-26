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
	UUIDBufferSize512   = 512
	UUIDBufferSize4096  = 4096
	UUIDBufferSize65535 = 65535
	UUIDBufferSize1MB   = 1024 * 1024
	UUIDBufferSize10MB  = 10 * 1024 * 1024
	UUIDBufferSize100MB = 100 * 1024 * 1024
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

	// Start a goroutine to dump the buffer contents to stdout periodically
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-time.After(time.Second):
				dumpAndClearBuffer(&buffer512, os.Stdout, "512-byte buffer")
				dumpAndClearBuffer(&buffer4096, os.Stdout, "4096-byte buffer")
				dumpAndClearBuffer(&buffer65535, os.Stdout, "65535-byte buffer")
				dumpAndClearBuffer(&buffer1MB, os.Stdout, "1MB buffer")
				dumpAndClearBuffer(&buffer10MB, os.Stdout, "10MB buffer")
				dumpAndClearBuffer(&buffer100MB, os.Stdout, "100MB buffer")
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
			dumpBuffer(&buffer512, os.Stdout, "512-byte buffer")
		} else if isBufferFull(&buffer4096) {
			dumpBuffer(&buffer4096, os.Stdout, "4096-byte buffer")
		} else if isBufferFull(&buffer65535) {
			dumpBuffer(&buffer65535, os.Stdout, "65535-byte buffer")
		} else if isBufferFull(&buffer1MB) {
			dumpBuffer(&buffer1MB, os.Stdout, "1MB buffer")
		} else if isBufferFull(&buffer10MB) {
			dumpBuffer(&buffer10MB, os.Stdout, "10MB buffer")
		} else if isBufferFull(&buffer100MB) {
			dumpBuffer(&buffer100MB, os.Stdout, "100MB buffer")
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
func dumpBuffer(buffer *[]string, writer io.Writer, bufferName string) {
	fmt.Printf("Dumping %s...\n", bufferName)
	uuidBlob := strings.Join(*buffer, "")
	fmt.Fprintln(writer, uuidBlob)
}

// dumpAndClearBuffer dumps the buffer contents to the provided writer and clears the buffer.
func dumpAndClearBuffer(buffer *[]string, writer io.Writer, bufferName string) {
	dumpBuffer(buffer, writer, bufferName)
	*buffer = (*buffer)[:0]
}

