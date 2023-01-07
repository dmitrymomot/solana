package utils

import "github.com/mr-tron/base58"

// Base58ToBytes converts base58 string to bytes.
func Base58ToBytes(s string) ([]byte, error) {
	return base58.Decode(s)
}
