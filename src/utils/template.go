package utils

import (
	"fmt"
	"html/template"
	"strings"
	"time"
)

func TimeUnixFormat(unix int64, layout string) string {
	t := time.Unix(unix, 0)
	return t.Format(layout)
}

// format time unixstamp friendly,
// like xxx seconds ago
func TimeUnixFriend(unixStamp int64) string {
	t := time.Unix(unixStamp, 0)
	seconds := int64(time.Since(t).Seconds())
	if seconds < 0 {
		return "FUTURE"
	}
	if seconds < 60 {
		return fmt.Sprintf("%d Seconds Ago", seconds)
	}
	if seconds < 3600 {
		return fmt.Sprintf("%d Minutes Ago", seconds/60)
	}
	if seconds < 86400 {
		return fmt.Sprintf("%d Hours Ago", seconds/3600)
	}
	if seconds < 86400*30 {
		return fmt.Sprintf("%d Days Ago", seconds/86400)
	}
	return t.Format("01.02")
}

// newline 2 break
func Nl2Br(str string) template.HTML {
	return template.HTML(Nl2BrString(str))
}

// newline 2 break, return string
func Nl2BrString(str string) string {
	return strings.Replace(str, "\n", "<br/>", -1)
}

func FriendBytesSize(size int64) string {
	sFloat := float64(size)
	if sFloat >= 1024*1024 {
		return fmt.Sprintf("%.1f MB", sFloat/1024/1024)
	}
	if sFloat > 1024 {
		return fmt.Sprintf("%.1f KB", sFloat/1024)
	}
	return fmt.Sprintf("%.1f B", sFloat)
}
