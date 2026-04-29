package presenter

// Type is a generic presenter interface.
// I is the use-case output type; R is the HTTP response object type.
type Type[I any, R any] interface {
	Output(output I) R
	Error(err error) R
}
