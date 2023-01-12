package utils

import "encoding/json"

// PrettyPrint returns a pretty-printed string representation of a struct or a map.
func PrettyPrint(s interface{}) string {
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return ""
	}

	return string(b)
}
