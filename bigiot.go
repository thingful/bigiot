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
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/shurcooL/graphql"
)

const (
	// DefaultMarketplaceURL is the URI to the default BIG IoT marketplace
	DefaultMarketplaceURL = "https://market.big-iot.org"

	// DefaultTimeout is the default timeout in seconds to set on on requests to
	// the marketplace
	DefaultTimeout = 10
)

// Config captures the runtime state of the BIGIoT client. We've separated the
// config from the Provider type so that we can pass this config into our
// authenticating http Transport without all the other references.
type Config struct {
	ID          string
	Secret      string
	AccessToken string
	UserAgent   string
}

// Provider is our struct for containing auth credentials to interact with the
// BIGIoT marketplace.
type Provider struct {
	*Config
	BaseURL       *url.URL
	httpClient    *http.Client
	graphqlClient *graphql.Client
}

// NewProvider instantiates and returns a configured Provider instance. The
// required parameters to the function are the provider ID and secret. If you
// want to connect to a marketplace other than the offical marketplace (i.e.
// connecting to a local instance for testing), you can configure this by means
// of the variadic third parameter, which can be used for additional
// configuration.
func NewProvider(id, secret string, options ...func(*Provider) error) (*Provider, error) {
	// this is a known good URL, so we can ignore the error here
	u, _ := url.Parse(DefaultMarketplaceURL)

	// set up a default http client, that enforces our default timeout. Users will
	// have to explicitly override if they want a non-timing out client.
	httpClient := &http.Client{
		Timeout: time.Second * DefaultTimeout,
	}

	provider := &Provider{
		Config: &Config{
			ID:        id,
			Secret:    secret,
			UserAgent: fmt.Sprintf("bigiot/%s (https://github.com/thingful/bigiot)", Version),
		},
		BaseURL:    u,
		httpClient: httpClient,
	}

	var err error

	// apply all functional options
	for _, opt := range options {
		err = opt(provider)
		if err != nil {
			return nil, err
		}
	}

	// now wrap our transport to add authentication header if accessToken is available
	transport := http.DefaultTransport
	if provider.httpClient.Transport != nil {
		transport = provider.httpClient.Transport
	}

	provider.httpClient.Transport = &authTransport{
		proxied: transport,
		Config:  provider.Config,
	}

	// setup our graphql client pointing at the specified marketplace, and using
	// our auth enabled http client
	graphqlURL := *provider.BaseURL
	graphqlURL.Path = "/graphql"

	provider.graphqlClient = graphql.NewClient(graphqlURL.String(), provider.httpClient, nil)

	return provider, nil
}

// Authenticate attempts to connect to the marketplace and authenticate the
// client. Returns any error to the caller.
func (p *Provider) Authenticate() (err error) {
	// deference to make sure we clone our baseURL property rather than modifying
	// the pointed to value
	authURL := *p.BaseURL
	authURL.Path = "/accessToken"

	params := &url.Values{
		"clientId":     []string{p.ID},
		"clientSecret": []string{p.Secret},
	}

	authURL.RawQuery = params.Encode()

	req, err := http.NewRequest(http.MethodGet, authURL.String(), nil)
	if err != nil {
		return err
	}

	req.Header.Set(acceptKey, textPlain)

	resp, err := p.httpClient.Do(req)
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

	p.AccessToken = string(body)

	return nil
}

func (p *Provider) Offering(id string) (*Offering, error) {
	var q struct {
		Offering struct {
			ID   graphql.String
			Name graphql.String
		} `graphql:"offering(id: $id)"`
	}

	variables := map[string]interface{}{
		"id": graphql.String(id),
	}

	err := p.graphqlClient.Query(context.Background(), &q, variables)
	if err != nil {
		return nil, err
	}

	return &Offering{
		ID: string(q.Offering.ID),
		OfferingDescription: OfferingDescription{
			Name: string(q.Offering.Name),
		},
	}, nil
}

func (p *Provider) RegisterOffering(offeringDescription *OfferingDescription) (*Offering, error) {
	return nil, nil
}

// WithMarketplace is a functional configuration option allowing us to
// optionally set a custom marketplace URI when constructing a Provider
// instance.
func WithMarketplace(marketplaceURI string) func(*Provider) error {
	return func(p *Provider) error {
		u, err := url.Parse(marketplaceURI)
		if err != nil {
			return err
		}

		p.BaseURL = u

		return nil
	}
}

// WithUserAgent allows the caller to specify the user agent that should be sent
// to the marketplace.
func WithUserAgent(userAgent string) func(*Provider) error {
	return func(p *Provider) error {
		p.UserAgent = userAgent

		return nil
	}
}

// WithHTTPClient allows a caller to pass in a custom http Client allowing them
// to customize the behaviour of our HTTP interactions.
func WithHTTPClient(client *http.Client) func(*Provider) error {
	return func(p *Provider) error {
		p.httpClient = client

		return nil
	}
}
