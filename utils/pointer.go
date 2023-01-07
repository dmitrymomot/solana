package utils

// Pointer converts any type to a pointer.
// This is useful for passing a value to a function that expects a pointer.
func Pointer[T any](i T) *T {
	return &i
}
