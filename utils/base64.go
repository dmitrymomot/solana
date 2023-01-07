package utils

import "encoding/base64"

// Base64ToBytes converts base64 string to bytes.
func Base64ToBytes(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

// BytesToBase64 converts bytes to base64 string.
func BytesToBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}
