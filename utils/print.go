package utils

import "encoding/json"

// StructPrettyPrint returns a pretty-printed string representation of a struct.
func StructPrettyPrint(s interface{}) string {
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return ""
	}

	return string(b)
}
