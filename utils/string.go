package utils

import "strings"

// Trim string between two substrings and return the string without it and substrings.
func TrimStringBetween(str, start, end string) string {
	indx1 := strings.Index(str, start)
	indx2 := strings.Index(str, end)
	if indx1 == -1 || indx2 == -1 {
		return strings.TrimSpace(str)
	}
	return strings.TrimSpace(str[:indx1] + str[indx2+len(end):])
}

// TrimRightZeros trims trailing zeros from string.
func TrimRightZeros(str string) string {
	return strings.TrimRight(str, "0")
}
