package model

import "time"

type Time struct {
	Raw   time.Time
	Year  int
	Month int
	Day   int
}

func NewTime(t time.Time) Time {
	ti := Time{
		Raw: t,
	}
	ti.Year = ti.Raw.Year()
	ti.Month = int(ti.Raw.Month())
	ti.Day = ti.Raw.Day()
	return ti
}
