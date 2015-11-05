package builder

import "time"

// build report
type Report struct {
	DstDir string
	Begin  time.Time
	End    time.Time
	Error  error
	Tree   *reportTree
}

type reportTree struct {
	Index      string
	OriginFile string
	DestFile   string
	BuildType  string
	Children   map[string]*reportTree
}
