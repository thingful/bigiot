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
	"fmt"
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
	config  *Base
	proxied http.RoundTripper
}

// RoundTrip is our implementation of RoundTripper, which does the job of adding
// auth credentials if any are present
func (t authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	fmt.Println(req.Body)
	if t.config.accessToken != "" {
		req.Header.Set(authorizationKey, "Bearer "+t.config.accessToken)
	}

	// set our internal user agent if the client hasn't supplied one
	if req.Header.Get(userAgentKey) == "" {
		req.Header.Set(userAgentKey, t.config.userAgent)
	}

	res, err := t.proxied.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	fmt.Println(res)

	return res, err
}
