package utils

import "time"

func TimeUnixFormat(unix int64, layout string) string {
	t := time.Unix(unix, 0)
	return t.Format(layout)
}
