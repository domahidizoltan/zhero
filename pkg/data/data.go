// Package data is for data helpers
package data

func Ptr[T any](v T) *T {
	return &v
}
