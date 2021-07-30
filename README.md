# processmessage

Go library for displaying progress messages.

## Usage example

```go
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/nxshock/progressmessage"
)

func main() {
	// Create new message
	pm := progressmessage.New("Progress: %d%%...")

	// Start message display
	pm.Start()

	// Let's do some job
	for i := 0; i < 100; i++ {
		// Simulate some work
		time.Sleep(time.Second / 5)

		// Update progress variables in same order as specified on creating the message
		pm.Update(i + 1)
	}

	// Stop message display
	pm.Stop()

	// Cursor stays in progress message position so you can display result message manually
	fmt.Fprintln(os.Stderr, "\rProcessing finished.")
}
```
