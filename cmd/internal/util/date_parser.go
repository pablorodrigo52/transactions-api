package util

import (
	"errors"
	"time"
)

// ParseDateWithFormat parses a date string with a given format and returns a time.Time object.
func ParseDateWithFormat(dateStr, format string) (time.Time, error) {
	if format == "" {
		format = time.RFC3339
	}

	parsedDate, err := time.Parse(format, dateStr)
	if err != nil {
		return time.Time{}, errors.New("invalid date format expected " + format)
	}
	return parsedDate, nil
}

// ParseDate parses a date string in the format "2006-01-02T15:04:05Z07:00" and returns a time.Time object.
func ParseDate(dateStr string) (time.Time, error) {
	return ParseDateWithFormat(dateStr, "")
}

// FormatDate formats a time.Time object to a string in the format "2006-01-02T15:04:05Z07:00".
func FormatDate(date time.Time) string {
	return date.Local().Format(time.RFC3339)
}
