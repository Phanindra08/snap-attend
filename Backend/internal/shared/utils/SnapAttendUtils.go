package utils

import "strings"

// ConvertToLowerCaseString normalizes email addresses
func ConvertToLowerCaseString(content string) string {
	return strings.TrimSpace(strings.ToLower(content))
}
