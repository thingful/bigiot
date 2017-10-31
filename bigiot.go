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
	"context"
	"net/http"
	"net/url"

	"github.com/shurcooL/graphql"
)

const (
	// DefaultMarketplaceURL is the URI to the default BIG IoT marketplace
	DefaultMarketplaceURL = "https://market.big-iot.org"

	// DefaultTimeout is the default timeout in seconds to set on on requests to
	// the marketplace
	DefaultTimeout = 10
)

// Base is our base BIGIoT client object. It contains the runtime state of our
// BIG IoT client implementations. It includes a number of unexported fields, as
// uses of this library are not permitted to modify this object directly; rather
// they should use one of the functional configuration functions when
// initializing an instance of the client.
type Base struct {
	userAgent     string
	httpClient    *http.Client
	baseURL       *url.URL
	id            string
	secret        string
	accessToken   string
	graphqlClient *graphql.Client
}

func (b *Base) Query(ctx context.Context, q interface{}, variables map[string]interface{}) error {
	return b.graphqlClient.Query(ctx, q, variables)
}

func (b *Base) Mutate(ctx context.Context, m interface{}, input Input, variables map[string]interface{}) error {
	if variables == nil {
		variables = map[string]interface{}{"input": input}
	} else {
		variables["input"] = input
	}

	return b.graphqlClient.Mutate(ctx, m, variables)
}

// Option is a type alias for our functional configuration type.
type Option func(*Base) error

// WithMarketplace is a functional configuration option allowing us to
// optionally set a custom marketplace URI when constructing a BIGIoT instance.
func WithMarketplace(marketplaceURI string) Option {
	return func(c *Base) error {
		u, err := url.Parse(marketplaceURI)
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
	return func(b *Base) error {
		b.userAgent = userAgent

		return nil
	}
}

// WithHTTPClient allows a caller to pass in a custom http Client allowing them
// to customize the behaviour of our HTTP interactions.
func WithHTTPClient(client *http.Client) Option {
	return func(b *Base) error {
		b.httpClient = client

		return nil
	}
}
