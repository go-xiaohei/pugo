package helper

import (
	"bytes"
	"fmt"
	"strings"

	"gopkg.in/inconshreveable/log15.v2"
)

// LogfmtFormat format logs as fmt.Sprintf
// Example:
//  Debug("Debug|%s|%d","a",1)
func LogfmtFormat() log15.Format {
	return log15.FormatFunc(func(r *log15.Record) []byte {
		var color = 0
		switch r.Lvl {
		case log15.LvlCrit:
			color = 35
		case log15.LvlError:
			color = 31
		case log15.LvlWarn:
			color = 33
		case log15.LvlInfo:
			color = 32
		case log15.LvlDebug:
			color = 36
		}
		t := r.Time.Format("01-02 15:04:05.999")
		b := &bytes.Buffer{}
		lvl := strings.ToUpper(r.Lvl.String())
		format := ""
		if color > 0 {
			format = fmt.Sprintf("\x1b[%dm%s\x1b[0m|%s|%s ", color, lvl, friendTime(t), r.Msg)
		} else {
			format = fmt.Sprintf("[%s] [%s] %s ", lvl, friendTime(t), r.Msg)
		}
		b.WriteString(fmt.Sprintf(format, r.Ctx...))
		b.WriteString("\n")
		return b.Bytes()
	})
}

func friendTime(t string) string {
	if len(t) < 18 {
		return t + strings.Repeat(" ", 18-len(t))
	}
	return t
}
