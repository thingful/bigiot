package bigiot

import (
	"net/http"
)

// authTransport is an internal implementation of the RoundTripper interface that
// we use to wrap the transport on the http.Client used for making requests to
// the BIGIoT marketplace. This custom transport adds auth credentials if any
// are set, and also adds a user-agent string to send to the server.
type authTransport struct {
	accessToken string
	proxied     http.RoundTripper
}

// RoundTrip is our implementation of RoundTripper, which does the job of adding
// auth credentials if any are present
func (t authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+t.accessToken)
	}

	req.Header.Set("User-Agent", "bigiot-go")

	res, err := t.proxied.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	return res, err
}
