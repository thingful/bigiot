package bigiot

import (
	"net/http"
)

const (
	userAgentKey = "User-Agent"

	authorizationKey = "Authorization"

	acceptKey = "Accept"

	textPlain = "text/plain"
)

// authTransport is an internal implementation of the RoundTripper interface that
// we use to wrap the transport on the http.Client used for making requests to
// the BIGIoT marketplace. This custom transport adds auth credentials if any
// are set, and also adds a user-agent string to send to the server.
type authTransport struct {
	*Config
	proxied http.RoundTripper
}

// RoundTrip is our implementation of RoundTripper, which does the job of adding
// auth credentials if any are present
func (t authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.AccessToken != "" {
		req.Header.Set(authorizationKey, "Bearer "+t.AccessToken)
	}

	// set our internal user agent if the client hasn't supplied one
	if req.Header.Get(userAgentKey) == "" {
		req.Header.Set(userAgentKey, t.UserAgent)
	}

	res, err := t.proxied.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	return res, err
}
