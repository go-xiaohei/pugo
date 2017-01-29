package helper

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMarkdown(t *testing.T) {
	Convey("Markdown", t, func() {
		h1 := []byte("#h1")
		So(string(Markdown(h1)), ShouldEqual, `<h1 id="h1">h1</h1>`+"\n")

		code := []byte("```go\npackage main\n```")
		So(string(Markdown(code)), ShouldEqual, `<pre><code class="language-go">package main
</code></pre>
`)
	})
}
