package builder

import (
	"strings"
	"time"
)

// build report
type Report struct {
	DstDir string
	Begin  time.Time
	Error  error
	Tree   *reportTree
}

func newReport(dstDir string) *Report {
	return &Report{
		DstDir: dstDir,
		Begin:  time.Now(),
		Tree: &reportTree{
			Index: "HOME",
		},
	}
}

// build-process duration
func (r *Report) Duration() time.Duration {
	return time.Since(r.Begin)
}

// record route tree
func (r *Report) addTree(url, urlType, originFile, destFile string) {
	if !strings.HasPrefix(url, "/") {
		url = "/" + url
	}
	url = "HOME" + url
	println(url)

}

type reportTree struct {
	Index      string
	OriginFile string
	DestFile   string
	BuildType  string
	Children   map[string]*reportTree
}
