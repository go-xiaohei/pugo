package builder

import (
	"strings"
	"time"
)

// build report
// todo : finish this Report
type Report struct {
	DstDir    string
	BeginTime time.Time
	Error     error
	Tree      *reportTree
}

func newReport(dstDir string) *Report {
	return &Report{
		DstDir:    dstDir,
		BeginTime: time.Now(),
		Tree: &reportTree{
			Index: "HOME",
		},
	}
}

// build-process duration
func (r *Report) Duration() time.Duration {
	return time.Since(r.BeginTime)
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
