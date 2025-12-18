// Package middlewares have focus in middlewares for validation or parse information
package middlewares

// MwKey is the keys used internal for context
type MwKey string

// HeaderKey is the keys used in headers by middlewares
type HeaderKey string

// String convert HeaderKey in string
func (h HeaderKey) String() string {
	return string(h)
}
