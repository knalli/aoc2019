package dayless

import (
	"fmt"
	"time"
)

/**
Convenient method for tracking function's usage time
https://coderwall.com/p/cp5fya/measuring-execution-time-in-go
*/
func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("⏱ %s took %s\n", name, elapsed)
}
