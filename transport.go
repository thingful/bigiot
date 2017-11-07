// Copyright 2017 Thingful Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bigiot

import (
	"bytes"
	"net/http"
)

// authTransport is an internal implementation of the RoundTripper interface that
// we use to wrap the transport on the http.Client used for making requests to
// the BIGIoT marketplace. This custom transport adds auth credentials if any
// are set, and also adds a user-agent string to send to the server.
type authTransport struct {
	bigiot  *BIGIoT
	proxied http.RoundTripper
}

// RoundTrip is our implementation of RoundTripper, which does the job of adding
// auth credentials if any are present
func (t authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.bigiot.accessToken != "" {
		var buf bytes.Buffer
		buf.WriteString("Bearer ")
		buf.WriteString(t.bigiot.accessToken)
		req.Header.Set(authorizationHeader, buf.String())
	}

	// set our internal user agent if the client hasn't supplied one
	if req.Header.Get(userAgentHeader) == "" {
		req.Header.Set(userAgentHeader, t.bigiot.userAgent)
	}

	res, err := t.proxied.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	return res, err
}
