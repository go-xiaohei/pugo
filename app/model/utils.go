package model

import "time"

// Time is custom time struct, use to display more time data without function calling
type Time struct {
	Raw   time.Time
	Year  int
	Month int
	Day   int
}

// Format returns layout time string
func (t Time) Format(layout string) string {
	return t.Raw.Format(layout)
}

// NewTime parses str as time,
// if error, use t2
func NewTime(str string, t2 time.Time) Time {
	var (
		t   time.Time
		err error
	)
	if len(str) <= 10 {
		t, err = time.Parse("2006-01-02", str)
		if err != nil {
			t = t2
		}
	} else if len(str) <= 16 {
		t, err = time.Parse("2006-01-02 15:04", str)
		if err != nil {
			t = t2
		}
	} else {
		t, err = time.Parse("2006-01-02 15:04:05", str)
		if err != nil {
			t = t2
		}
	}
	ti := Time{
		Raw: t,
	}
	ti.Year = ti.Raw.Year()
	ti.Month = int(ti.Raw.Month())
	ti.Day = ti.Raw.Day()
	return ti
}
