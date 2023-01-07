package utils

import "encoding/json"

// MapToStruct converts a map to a struct.
func MapToStruct(m map[string]interface{}, s interface{}) error {
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, s)
}
