package metadata

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// MetadataFromJSON parses the metadata from JSON
// The JSON must be a valid JSON string that can be unmarshalled into a Metadata struct
func MetadataFromJSON(data []byte) (*Metadata, error) {
	m := &Metadata{}
	if err := json.Unmarshal(data, m); err != nil {
		return nil, fmt.Errorf("failed to decode metadata from json string: %w", err)
	}

	return m, nil
}

// MetadataFromURI parses the metadata from a URI
// The URI must be a valid HTTP(S) URL
func MetadataFromURI(uri string) (*Metadata, error) {
	if uri == "" {
		return nil, nil
	}

	resp, err := http.Get(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to download metadata from uri: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata from uri: %w", err)
	}

	return MetadataFromJSON(body)
}
