package bigiot

const (
	// ErrUnexpectedResponse is an error returned when we receive an unexpected
	// response from the BIG IoT API.
	ErrUnexpectedResponse = Error("Unexpected HTTP response code")
)

// Error is a type alias for string, allowing us to export const error values
type Error string

// Error is the implementation of the error interface, allowing us to use our
// Error type as an error.
func (e Error) Error() string {
	return string(e)
}
