package adafruitio

//go:generate go run gen-accessors.go

const (
	// Version is the version of this library. It is sent in the User-Agent
	// header of every request.
	Version = "3.0.0"
)

// Ptr returns a pointer to v. Use it to set the pointer fields of Feed,
// Group, and Data, for example:
//
//	feed := &adafruitio.Feed{Name: adafruitio.Ptr("Temperature"), History: adafruitio.Ptr(false)}
func Ptr[T any](v T) *T {
	return &v
}
