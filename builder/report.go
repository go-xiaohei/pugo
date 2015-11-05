package builder

import "time"

// build report
type Report struct {
	Time  time.Time
	Error error
}
