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
	"encoding/json"
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

// RegisterOffering allows calles to register an offering on the marketplace.
// When registering the caller will supply an activation lifetime for the
// Offering as part of the input AddOffering instance. The function returns a
// populated Offering instance or nil and an error.
func (p *Provider) RegisterOffering(ctx context.Context, offering *AddOffering) (*Offering, error) {
	offering.providerID = p.ID

	body, err := p.Query(ctx, offering)
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

// DeleteOffering attempts to delete or unregister an offering on the
// marketplace. It is called with a context, and a DeleteOffering instance. This
// instance is serialized and the query executed against the GraphQL server. The
// function returns an error if anything goes wrong.
func (p *Provider) DeleteOffering(ctx context.Context, offering *DeleteOffering) error {
	_, err := p.Query(ctx, offering)
	if err != nil {
		return err
	}

	return nil
}

// addOfferingResponse is a unexported type used when parsing the response from
// calling RegisterOffering
type addOfferingResponse struct {
	Data struct {
		Offering Offering `json:"addOffering"`
	} `json:"data"`
}
