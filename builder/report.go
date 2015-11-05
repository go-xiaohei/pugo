package builder

import "time"

// build report
type Report struct {
	DstDir string
	Begin  time.Time
	End    time.Time
	Error  error
}
