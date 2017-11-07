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
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Provider is our type for interacting with the marketplace from the
// perspective of a data provider. We embed the base Config type which stores
// our runtime configuration (auth credentials, base url etc.).
type Provider struct {
	*BIGIoT
}

// NewProvider instantiates and returns a configured Provider instance. The
// required parameters to the function are the provider ID and secret. If you
// want to connect to a marketplace other than the offical marketplace (i.e.
// connecting to a local instance for testing), you can configure this by means
// of the variadic third parameter, which can be used for additional
// configuration.
func NewProvider(id, secret string, options ...Option) (*Provider, error) {
	base, err := NewBIGIoT(id, secret, options...)
	if err != nil {
		return nil, err
	}

	return &Provider{BIGIoT: base}, nil
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

func (p *Provider) RegisterOffering(ctx context.Context, offering *OfferingInput) (*Offering, error) {
	offering.providerID = p.ID

	q := &Query{
		Query: offering.Serialize(),
	}

	b, err := json.Marshal(q)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, p.graphqlURL, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	req.Header.Set(contentTypeHeader, applicationJSON)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		if cerr := resp.Body.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, ErrUnexpectedResponse
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	addOffering := addOfferingResponse{}

	err = json.Unmarshal(body, &addOffering)
	if err != nil {
		return nil, err
	}

	return &addOffering.Data.Offering, nil
}

// addOfferingResponse is a unexported type used when parsing the response from
// calling RegisterOffering
type addOfferingResponse struct {
	Data struct {
		Offering Offering `json:"addOffering"`
	} `json:"data"`
}