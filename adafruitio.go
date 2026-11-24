// Package adafruit provides bindings for the io.adafruit.com APIs

package adafruitio

const (
	Version = "3.0.0"
)

// Bool returns a pointer to v. Use it to set optional boolean fields, such as
// Feed.History.
func Bool(v bool) *bool {
	return &v
}
