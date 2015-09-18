package utils

import (
	"fmt"
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
