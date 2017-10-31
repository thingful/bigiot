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

// Provider is our type for interacting with the marketplace from the
// perspective of a data provider. We embed the base Config type which stores
// our runtime configuration (auth credentials, base url etc.).
type Provider struct {
	*Base
}

// NewProvider instantiates and returns a configured Provider instance. The
// required parameters to the function are the provider ID and secret. If you
// want to connect to a marketplace other than the offical marketplace (i.e.
// connecting to a local instance for testing), you can configure this by means
// of the variadic third parameter, which can be used for additional
// configuration.
func NewProvider(id, secret string, options ...Option) (*Provider, error) {
	// this is a known good URL, so we can ignore the error here
	u, _ := url.Parse(DefaultMarketplaceURL)

	// set up a default http client, that enforces our default timeout. Users will
	// have to explicitly override if they want a non-timing out client.
	httpClient := &http.Client{
		Timeout: time.Second * DefaultTimeout,
	}

	config := &Base{
		id:         id,
		secret:     secret,
		userAgent:  fmt.Sprintf("bigiot/%s (https://github.com/thingful/bigiot)", Version),
		baseURL:    u,
		httpClient: httpClient,
	}

	var err error

	// apply all functional options
	for _, opt := range options {
		err = opt(config)
		if err != nil {
			return nil, err
		}
	}

	// now wrap our transport to add authentication header if accessToken is available
	transport := http.DefaultTransport
	if config.httpClient.Transport != nil {
		transport = config.httpClient.Transport
	}

	config.httpClient.Transport = &authTransport{
		proxied: transport,
		config:  config,
	}

	// setup our graphql client pointing at the specified marketplace, and using
	// our auth enabled http client
	graphqlURL := *config.baseURL
	graphqlURL.Path = "/graphql"

	config.graphqlClient = graphql.NewClient(graphqlURL.String(), config.httpClient, nil)

	return &Provider{Base: config}, nil
}

// Authenticate attempts to connect to the marketplace and authenticate the
// client. Returns any error to the caller.
func (p *Provider) Authenticate() (err error) {
	// deference to make sure we clone our baseURL property rather than modifying
	// the pointed to value
	authURL := *p.baseURL
	authURL.Path = "/accessToken"

	params := &url.Values{
		"clientId":     []string{p.id},
		"clientSecret": []string{p.secret},
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

	p.accessToken = string(body)

	return nil
}

// Offering returns details of an offering on being given the ID of that
// offering. It makes a call to the marketplace API and returns the offering
// details.
// func (p *Provider) Offering(id string) (*Offering, error) {
// 	var query struct {
// 		Offering struct {
// 			ID   graphql.String
// 			Name graphql.String
// 		} `graphql:"offering(id: $id)"`
// 	}
//
// 	variables := map[string]interface{}{
// 		"id": graphql.String(id),
// 	}
//
// 	err := p.graphqlClient.Query(context.Background(), &query, variables)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return &Offering{
// 		ID: string(query.Offering.ID),
// 		OfferingDescription: OfferingDescription{
// 			Name: string(query.Offering.Name),
// 		},
// 	}, nil
// }

func (p *Provider) RegisterOffering(ctx context.Context, offering *AddOffering) (*Offering, error) {
	var mutation struct {
		//AddOffering AddOffering `graphql:"addOffering(input: $addOffering)"`
		AddOffering AddOffering `graphql:"addOffering()"`
	}

	fmt.Println(mutation)

	err := p.Base.Mutate(ctx, &mutation, offering, nil)
	if err != nil {
		return nil, err
	}

	fmt.Println(mutation)

	return nil, nil
}
