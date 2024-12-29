package util

import (
	"errors"
	"time"
)

// ParseDate parses a date string in the format "2006-01-02" and returns a time.Time object.
func ParseDate(dateStr string) (time.Time, error) {
	const layout = "2006-01-02"
	parsedDate, err := time.Parse(layout, dateStr)
	if err != nil {
		return time.Time{}, errors.New("invalid date format, expected YYYY-MM-DD")
	}
	return parsedDate, nil
}

// FormatDate formats a time.Time object into a string in the format "2006-01-02".
func FormatDate(date time.Time) string {
	const layout = "2006-01-02"
	return date.Format(layout)
}
