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
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	// DefaultMarketplaceURL is the URI to the default BIG IoT marketplace
	DefaultMarketplaceURL = "https://market.big-iot.org"

	// DefaultTimeout is the default timeout in seconds to set on on requests to
	// the marketplace
	DefaultTimeout = 10
)

// BIGIoT is our base BIGIoT client object. It contains the runtime state of our
// BIG IoT client implementations. It includes a number of unexported fields, as
// uses of this library are not permitted to modify this object directly; rather
// they should use one of the functional configuration functions when
// initializing an instance of the client.
type BIGIoT struct {
	ID          string
	Secret      string
	userAgent   string
	httpClient  *http.Client
	baseURL     *url.URL
	accessToken string
	graphqlURL  string
}

func NewBIGIoT(id, secret string, options ...Option) (*BIGIoT, error) {
	// this is a known good URL, so we can ignore the error here
	u, _ := url.Parse(DefaultMarketplaceURL)

	// set up a default http client, that enforces our default timeout. Users will
	// have to explicitly override if they want a non-timing out client.
	httpClient := &http.Client{
		Timeout: time.Second * DefaultTimeout,
	}

	base := &BIGIoT{
		ID:         id,
		Secret:     secret,
		userAgent:  fmt.Sprintf("bigiot/%s (https://github.com/thingful/bigiot)", Version),
		baseURL:    u,
		httpClient: httpClient,
	}

	var err error

	// apply all functional options
	for _, opt := range options {
		err = opt(base)
		if err != nil {
			return nil, err
		}
	}

	// now wrap our transport to add authentication header if accessToken is available
	transport := http.DefaultTransport
	if base.httpClient.Transport != nil {
		transport = base.httpClient.Transport
	}

	base.httpClient.Transport = &authTransport{
		proxied: transport,
		bigiot:  base,
	}

	// setup our graphql client pointing at the specified marketplace, and using
	// our auth enabled http client
	graphqlURL := *base.baseURL
	graphqlURL.Path = "/graphql"

	base.graphqlURL = graphqlURL.String()

	return base, nil
}

// Authenticate makes a call to the /accessToken endpoint on the marketplace to
// obtain an access token which the client will then be able to use when making
// requests to the graphql endpoint. We make a GET request passing over our
// client id and secret, and get back a token if our credentials are valid.
func (b *BIGIoT) Authenticate() (err error) {
	// deference to make sure we clone our baseURL property rather than modifying
	// the pointed to value
	authURL := *b.baseURL
	authURL.Path = "/accessToken"

	params := &url.Values{
		"clientId":     []string{b.ID},
		"clientSecret": []string{b.Secret},
	}

	authURL.RawQuery = params.Encode()

	req, err := http.NewRequest(http.MethodGet, authURL.String(), nil)
	if err != nil {
		return err
	}

	req.Header.Set(acceptHeader, textPlain)

	resp, err := b.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		if cerr := resp.Body.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return ErrUnexpectedResponse
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	b.accessToken = string(body)

	return nil
}

// Option is a type alias for our functional configuration type.
type Option func(*BIGIoT) error

// WithMarketplace is a functional configuration option allowing us to
// optionally set a custom marketplace URI when constructing a BIGIoT instance.
func WithMarketplace(marketplaceURL string) Option {
	return func(c *BIGIoT) error {
		u, err := url.Parse(marketplaceURL)
		if err != nil {
			return err
		}

		c.baseURL = u

		return nil
	}
}

// WithUserAgent allows the caller to specify the user agent that should be
// sent to the marketplace.
func WithUserAgent(userAgent string) Option {
	return func(b *BIGIoT) error {
		b.userAgent = userAgent

		return nil
	}
}

// WithHTTPClient allows a caller to pass in a custom http Client allowing them
// to customize the behaviour of our HTTP interactions.
func WithHTTPClient(client *http.Client) Option {
	return func(b *BIGIoT) error {
		b.httpClient = client

		return nil
	}
}

// Serializable is an interface for an instance that can serialize itself into
// some form that the BIG IoT Marketplace will accept as input for either query
// or mutatation.
type Serializable interface {
	Serialize() string
}
