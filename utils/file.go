package utils

import (
	"net/url"
	"path/filepath"
	"strings"
)

// GetFileTypeByURI returns the file type of the given URI
func GetFileTypeByURI(uri string) string {
	ext := filepath.Ext(uri)
	if ext == "" {
		parsedUri, err := url.Parse(uri)
		if err != nil {
			return ""
		}
		ext = parsedUri.Query().Get("ext")
		if ext != "" {
			ext = parsedUri.Query().Get("format")
		}
	}
	ext = strings.Trim(ext, ".")

	switch ext {
	case "png":
		return "image/png"
	case "gif":
		return "image/gif"
	case "jpg", "jpeg":
		return "image/jpeg"
	case "mp4":
		return "video/mp4"
	case "mov":
		return "video/quicktime"
	case "mp3":
		return "audio/mpeg"
	case "flac":
		return "audio/flac"
	case "wav":
		return "audio/wav"
	case "glb":
		return "model/gltf-binary"
	case "gltf":
		return "model/gltf+json"
	case "html":
		return "text/html"
	case "js":
		return "application/javascript"
	case "css":
		return "text/css"
	case "json":
		return "application/json"
	case "xml":
		return "application/xml"
	case "svg":
		return "image/svg+xml"
	case "ico":
		return "image/x-icon"
	case "zip":
		return "application/zip"
	case "pdf":
		return "application/pdf"
	case "txt":
		return "text/plain"
	case "md":
		return "text/markdown"
	case "csv":
		return "text/csv"
	}

	return ""
}
