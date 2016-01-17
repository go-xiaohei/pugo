package helper

import "strings"

var (
	// TimeLayoutReplacer replace YYYY-MM-D layout to go time layout
	TimeLayoutReplacer = strings.NewReplacer(
		"H", "15",
		"hh", "15",
		"h", "03",
		"mm", "04",
		"ss", "05",
		"MMMM", "January",
		"MMM", "Jan",
		"MM", "01",
		"M", "1",
		"pm", "PM",
		"ZZZZ", "-0700",
		"ZZZ", "MST",
		"ZZ", "Z07:00",
		"YYYY", "2006",
		"YY", "06",
		"DDDD", "Monday",
		"DDD", "Mon",
		"DD", "02",
		"D", "2",
	)
)
