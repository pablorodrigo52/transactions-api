package util

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseDateWithFormat(t *testing.T) {
	tests := []struct {
		testName      string
		input         string
		format        string
		expected      time.Time
		expectedError string
	}{
		{
			testName:      "Parse date with success using default format",
			input:         "2023-10-10T10:10:10Z",
			format:        "",
			expected:      time.Date(2023, 10, 10, 10, 10, 10, 0, time.UTC),
			expectedError: "",
		},
		{
			testName:      "Parse date with success using custom format",
			input:         "10-10-2023",
			format:        "02-01-2006",
			expected:      time.Date(2023, 10, 10, 0, 0, 0, 0, time.UTC),
			expectedError: "",
		},
		{
			testName:      "Parse date error, invalid date",
			input:         "invalid-date",
			format:        "",
			expected:      time.Time{},
			expectedError: "invalid date format expected 2006-01-02T15:04:05Z07:00",
		},
		{
			testName:      "Parse date error, invalid format date",
			input:         "2023-10-10",
			format:        "2006-01-02T15:04:05Z07:00",
			expected:      time.Time{},
			expectedError: "invalid date format expected 2006-01-02T15:04:05Z07:00",
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			result, err := ParseDateWithFormat(test.input, test.format)

			if err != nil {
				assert.Error(t, err)
				assert.Equal(t, test.expectedError, err.Error())
			} else {
				assert.Nil(t, err)
				assert.Equal(t, test.expected, result)
			}
		})
	}
}

func TestParseDate(t *testing.T) {
	tests := []struct {
		testName      string
		input         string
		expected      time.Time
		expectedError string
	}{
		{
			testName:      "Parse date with success",
			input:         "2023-10-10T10:10:10Z",
			expected:      time.Date(2023, 10, 10, 10, 10, 10, 0, time.UTC),
			expectedError: "",
		},
		{
			testName:      "Parse date error, invalid date",
			input:         "invalid-date",
			expected:      time.Time{},
			expectedError: "invalid date format expected 2006-01-02T15:04:05Z07:00",
		},
		{
			testName:      "Parse date error, invalid format date",
			input:         "2023-10-10",
			expected:      time.Time{},
			expectedError: "invalid date format expected 2006-01-02T15:04:05Z07:00",
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			result, err := ParseDate(test.input)

			if err != nil {
				assert.Error(t, err)
				assert.Equal(t, test.expectedError, err.Error())
			} else {
				assert.Nil(t, err)
				assert.Equal(t, test.expected, result)
			}
		})
	}
}

func TestFormatDate(t *testing.T) {
	tests := []struct {
		testName string
		input    time.Time
		expected string
	}{
		{
			testName: "Format date with success",
			input:    time.Date(2023, 10, 10, 10, 10, 10, 0, time.UTC),
			expected: "2023-10-10T07:10:10-03:00",
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			result := FormatDate(test.input)
			assert.Equal(t, test.expected, result)
		})
	}
}
