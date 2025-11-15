package utils

import (
	"strings"
	"time"
)

// TrimAndConvertToLowerCase - Helps in trimming whitespaces and converting any string to lowercase
func TrimAndConvertToLowerCase(content string) string {
	return strings.TrimSpace(strings.ToLower(content))
}

var getStartOfToday = func() time.Time {
	now := time.Now()
	year, month, day := now.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, now.Location())
}
