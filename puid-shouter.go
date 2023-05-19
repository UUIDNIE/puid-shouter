package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/google/uuid"
)

const (
	UUIDStringSize = 36  // Size of a UUID string
	MB             = 1e6 // Number of bytes in a megabyte
)

func generateUUIDs(wg *sync.WaitGroup, writer *bufio.Writer, genBigMsg bool) {
	defer wg.Done()

	// Buffer for holding the big message
	var bigMsgBuffer []string
	bigMsgTargetSize := 100 * MB // 100MB

	for {
		// Generate a big message once per second
		if genBigMsg && time.Now().Second()%10 == 0 {
			// Calculate how many UUIDs we need to reach the target size
			numUUIDs := int(bigMsgTargetSize / UUIDStringSize)

			for i := 0; i < numUUIDs; i++ {
				bigMsgBuffer = append(bigMsgBuffer, uuid.New().String())
			}

			// Write the big message to the buffer
			_, _ = fmt.Fprint(writer, bigMsgBuffer)

			// Clear the big message buffer
			bigMsgBuffer = nil
		} else {
			// Regular UUID generation
			_, _ = fmt.Fprintln(writer, uuid.New())
		}
	}
}

func main() {
	var wg sync.WaitGroup
	numCPU := runtime.NumCPU()

	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	genBigMsg := true

	for i := 0; i < numCPU; i++ {
		wg.Add(1)
		go generateUUIDs(&wg, writer, genBigMsg)

		// Only generate big messages in the first goroutine
		genBigMsg = false
	}

	wg.Wait()
}

