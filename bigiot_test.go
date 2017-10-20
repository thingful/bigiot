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

package bigiot_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/thingful/httpmock"

	"github.com/thingful/bigiot"
)

func TestProviderConstructor(t *testing.T) {
	p, err := bigiot.NewProvider("id", "secret")
	assert.Nil(t, err)
	assert.Equal(t, "id", p.ID)
	assert.Equal(t, "secret", p.Secret)
	assert.Equal(t, bigiot.DefaultMarketplaceURL, p.BaseURL.String())
}

func TestProviderConstructorWithMarketplace(t *testing.T) {
	p, err := bigiot.NewProvider(
		"id",
		"secret",
		bigiot.WithMarketplace("https://market-dev.big-iot.org"),
	)
	assert.Nil(t, err)
	assert.Equal(t, "https://market-dev.big-iot.org", p.BaseURL.String())
}

func TestProviderConstructorInvalidMarketplace(t *testing.T) {
	_, err := bigiot.NewProvider(
		"id",
		"secret",
		bigiot.WithMarketplace("http ://market-dev.big-iot.org"),
	)
	assert.NotNil(t, err)
}

func TestAuthenticate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterStubRequest(
		httpmock.NewStubRequest(
			http.MethodGet,
			"https://market.big-iot.org/accessToken?clientId=id&clientSecret=secret",
			httpmock.NewStringResponder(200, "1234abcd"),
		).WithHeader(
			&http.Header{
				"User-Agent": []string{"bigiot/" + bigiot.Version + " (https://github.com/thingful/bigiot)"},
				"Accept":     []string{"text/plain"},
			},
		),
	)

	p, _ := bigiot.NewProvider("id", "secret")
	assert.Nil(t, p.Authenticate())
	assert.Equal(t, "1234abcd", p.AccessToken)
	assert.Equal(t, bigiot.DefaultMarketplaceURL, p.BaseURL.String())

	assert.Nil(t, httpmock.AllStubsCalled())
}

func TestAuthenticateCustomUserAgent(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterStubRequest(
		httpmock.NewStubRequest(
			http.MethodGet,
			"https://market.big-iot.org/accessToken?clientId=id&clientSecret=secret",
			httpmock.NewStringResponder(200, "1234abcd"),
		).WithHeader(
			&http.Header{
				"User-Agent": []string{"foo"},
				"Accept":     []string{"text/plain"},
			},
		),
	)

	p, _ := bigiot.NewProvider(
		"id",
		"secret",
		bigiot.WithUserAgent("foo"),
	)
	assert.Nil(t, p.Authenticate())
	assert.Equal(t, "1234abcd", p.AccessToken)
	assert.Equal(t, bigiot.DefaultMarketplaceURL, p.BaseURL.String())

	assert.Nil(t, httpmock.AllStubsCalled())
}

type testTripper struct {
	proxied http.RoundTripper
}

func (t testTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.proxied.RoundTrip(req)
}

func TestAuthenticateCustomTransport(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterStubRequest(
		httpmock.NewStubRequest(
			http.MethodGet,
			"https://market.big-iot.org/accessToken?clientId=id&clientSecret=secret",
			httpmock.NewStringResponder(200, "1234abcd"),
		).WithHeader(
			&http.Header{
				"Accept": []string{"text/plain"},
			},
		),
	)

	client := &http.Client{
		Timeout:   1 * time.Second,
		Transport: testTripper{proxied: http.DefaultTransport},
	}

	p, _ := bigiot.NewProvider(
		"id",
		"secret",
		bigiot.WithHTTPClient(client),
	)
	assert.Nil(t, p.Authenticate())
	assert.Equal(t, "1234abcd", p.AccessToken)
	assert.Equal(t, bigiot.DefaultMarketplaceURL, p.BaseURL.String())

	assert.Nil(t, httpmock.AllStubsCalled())
}

func TestRegisterOffering(t *testing.T) {
	id := "Thingful-Temperature_Service"
	secret := "mT1WkU1DQ56SX61l4mtIvg=="

	provider, err := bigiot.NewProvider(id, secret, bigiot.WithMarketplace("https://market-dev.big-iot.org"))
	assert.Nil(t, err)

	offeringDescription := &bigiot.OfferingDescription{
		Name:    "Simple Weather",
		RdfType: bigiot.RdfType("bigiot:Weather"),
		Endpoints: []bigiot.Endpoint{
			{
				URI:          "http://example.com/weather",
				EndpointType: bigiot.HTTPGet,
			},
		},
		InputData: []bigiot.DataField{
			{
				Name:    "longitude",
				RdfType: bigiot.RdfType("schema:longitude"),
			},
			{
				Name:    "latitude",
				RdfType: bigiot.RdfType("schema:latitude"),
			},
		},
		OutputData: []bigiot.DataField{
			{
				Name:    "temperature",
				RdfType: bigiot.RdfType("schema:airTemperatureValue"),
			},
		},
		Extent: bigiot.Extent{
			City: "Edinburgh",
		},
	}

	assert.NotNil(t, offeringDescription)

	_, err = provider.RegisterOffering(offeringDescription)
	assert.Nil(t, err)
	//assert.NotNil(t, offering)
}

func TestOffering(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterStubRequest(
		httpmock.NewStubRequest(
			http.MethodGet,
			"https://market.big-iot.org/accessToken?clientId=id&clientSecret=secret",
			httpmock.NewStringResponder(200, "1234abcd"),
		).WithHeader(
			&http.Header{
				"Accept": []string{"text/plain"},
			},
		),
	)

	httpmock.RegisterStubRequest(
		httpmock.NewStubRequest(
			http.MethodPost,
			"https://market.big-iot.org/graphql",
			httpmock.NewStringResponder(200, `{
				"data": {
					"offering": {
						"id": "offeringID",
						"name": "offering name"
					}
				}
			}`),
		).WithHeader(
			&http.Header{
				"Authorization": []string{"Bearer 1234abcd"},
			},
		),
	)

	provider, err := bigiot.NewProvider("id", "secret")
	assert.Nil(t, err)

	err = provider.Authenticate()
	assert.Nil(t, err)

	offering, err := provider.Offering("offeringID")
	assert.Nil(t, err)
	assert.NotNil(t, offering)
	assert.Equal(t, "offeringID", offering.ID)
}
