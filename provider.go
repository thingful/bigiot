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
	"encoding/base64"
	"encoding/json"

	"github.com/pkg/errors"
	"gopkg.in/square/go-jose.v2/jwt"
)

// Provider is our type for interacting with the marketplace from the
// perspective of a data provider. We embed the base Config type which stores
// our runtime configuration (auth credentials, base url etc.).
type Provider struct {
	*base
}

// NewProvider instantiates and returns a configured Provider instance. The
// required parameters to the function are the provider ID and secret. If you
// want to connect to a marketplace other than the official marketplace (i.e.
// connecting to a local instance for testing), you can configure this by means
// of the variadic third parameter, which can be used for additional
// configuration.
func NewProvider(id, secret string, options ...Option) (*Provider, error) {
	b, err := newBase(id, secret, options...)
	if err != nil {
		return nil, err
	}

	return &Provider{base: b}, nil
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
func (p *Provider) RegisterOffering(ctx context.Context, offering *OfferingDescription) (*Offering, error) {
	offering.providerID = p.id

	body, err := p.query(ctx, offering)
	if err != nil {
		return nil, errors.Wrap(err, "error registering offering")
	}

	response := addOfferingResponse{}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshalling register offering json")
	}

	return &response.Data.Offering, nil
}

// DeleteOffering attempts to delete or unregister an offering on the
// marketplace. It is called with a context, and a DeleteOffering instance. This
// instance is serialized and the query executed against the GraphQL server. The
// function returns an error if anything goes wrong.
func (p *Provider) DeleteOffering(ctx context.Context, offering *DeleteOffering) error {
	_, err := p.query(ctx, offering)
	if err != nil {
		return errors.Wrap(err, "error deleting offering")
	}

	return nil
}

func (p *Provider) ActivateOffering(ctx context.Context, activation *ActivateOffering) (*Offering, error) {
	body, err := p.query(ctx, activation)
	if err != nil {
		return nil, errors.Wrap(err, "error re-activating offering")
	}

	response := activateOfferingResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshalling activation response")
	}

	return &response.Data.Offering, nil
}

// ValidateToken takes as input a string which should be JWT token generated by
// the marketplace and given to the client before it is allowed to access data
// from an offering. It takes as input the encoded token string, extracts its
// component parts and verifies the signature using the secret of the provider.
func (p *Provider) ValidateToken(tokenStr string) error {
	key, err := base64.StdEncoding.DecodeString(p.secret)
	if err != nil {
		return errors.Wrap(err, "decoding secret failed")
	}

	token, err := jwt.ParseSigned(tokenStr)
	if err != nil {
		return errors.Wrap(err, "error parsing token string")
	}

	cl := jwt.Claims{}
	err = token.Claims(key, &cl)
	if err != nil {
		return errors.Wrap(err, "error extracting claims from token")
	}

	// the only claim we validate for now is that the token has neither expired nor
	// is not valid yet. Note the jwt library allows leeway of one minute before
	// marking a token as invalid, I presume to allow for clock inconsistencies,
	// i.e. if the token expires 17:05, then Validate will still allow it up to
	// 17:06. The same applies before the token is technically valid.
	err = cl.Validate(jwt.Expected{
		Time: p.clock.Now(),
	})
	if err != nil {
		return errors.Wrap(err, "error validating claims")
	}

	// all good
	return nil
}

// addOfferingResponse is a unexported type used when parsing the response from
// calling RegisterOffering
type addOfferingResponse struct {
	Data struct {
		Offering Offering `json:"addOffering"`
	} `json:"data"`
}

// activateOfferingResponse is an unexported type used when parsing the response
// calling ActivateOffering
type activateOfferingResponse struct {
	Data struct {
		Offering Offering `json:"activateOffering"`
	} `json:"data"`
}
