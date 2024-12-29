package util

import (
	"errors"
	"time"
)

// ParseDate parses a date string in the format "2006-01-02T15:04:05Z07:00" and returns a time.Time object.
func ParseDate(dateStr string) (time.Time, error) {
	parsedDate, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return time.Time{}, errors.New("invalid date format expected 2006-01-02T15:04:05Z07:00")
	}
	return parsedDate, nil
}

// FormatDate formats a time.Time object to a string in the format "2006-01-02T15:04:05Z07:00".
func FormatDate(date time.Time) string {
	return date.Local().Format(time.RFC3339)
}
