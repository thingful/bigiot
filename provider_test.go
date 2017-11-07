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
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/thingful/simular"
)

func TestNewProvider(t *testing.T) {
	p, err := NewProvider("id", "secret")
	assert.Nil(t, err)
	assert.Equal(t, "id", p.ID)
	assert.Equal(t, "secret", p.Secret)
	assert.Equal(t, DefaultMarketplaceURL, p.baseURL.String())
}

func TestProviderConstructorWithMarketplace(t *testing.T) {
	p, err := NewProvider(
		"id",
		"secret",
		WithMarketplace("https://market-dev.big-iot.org"),
	)
	assert.Nil(t, err)
	assert.Equal(t, "https://market-dev.big-iot.org", p.baseURL.String())
}

func TestProviderConstructorInvalidMarketplace(t *testing.T) {
	_, err := NewProvider(
		"id",
		"secret",
		WithMarketplace("http ://market-dev.big-iot.org"),
	)
	assert.NotNil(t, err)
}

func TestAuthenticate(t *testing.T) {
	simular.Activate()
	defer simular.DeactivateAndReset()

	simular.RegisterStubRequests(
		simular.NewStubRequest(
			http.MethodGet,
			"https://market.big-iot.org/accessToken?clientId=id&clientSecret=secret",
			simular.NewStringResponder(200, "1234abcd"),
			simular.WithHeader(
				&http.Header{
					"User-Agent": []string{"bigiot/" + Version + " (https://github.com/thingful/bigiot)"},
					"Accept":     []string{"text/plain"},
				},
			),
		),
	)

	p, _ := NewProvider("id", "secret")
	assert.Nil(t, p.Authenticate())
	assert.Equal(t, "1234abcd", p.accessToken)
	assert.Equal(t, DefaultMarketplaceURL, p.baseURL.String())

	assert.Nil(t, simular.AllStubsCalled())
}

func TestAuthenticateUnexpectedResponse(t *testing.T) {
	simular.Activate()
	defer simular.DeactivateAndReset()

	simular.RegisterStubRequests(
		simular.NewStubRequest(
			http.MethodGet,
			"https://market.big-iot.org/accessToken?clientId=id&clientSecret=secret",
			simular.NewStringResponder(403, "Forbidden"),
		),
	)

	p, _ := NewProvider("id", "secret")
	err := p.Authenticate()
	assert.Equal(t, ErrUnexpectedResponse, err)
	assert.Equal(t, "Unexpected HTTP response code", err.Error())
}

func TestAuthenticateCustomUserAgent(t *testing.T) {
	simular.Activate()
	defer simular.DeactivateAndReset()

	simular.RegisterStubRequests(
		simular.NewStubRequest(
			http.MethodGet,
			"https://market.big-iot.org/accessToken?clientId=id&clientSecret=secret",
			simular.NewStringResponder(200, "1234abcd"),
			simular.WithHeader(
				&http.Header{
					"User-Agent": []string{"foo"},
					"Accept":     []string{"text/plain"},
				},
			),
		),
	)

	p, _ := NewProvider(
		"id",
		"secret",
		WithUserAgent("foo"),
	)
	assert.Nil(t, p.Authenticate())
	assert.Equal(t, "1234abcd", p.accessToken)
	assert.Equal(t, DefaultMarketplaceURL, p.baseURL.String())

	assert.Nil(t, simular.AllStubsCalled())
}

type testTripper struct {
	proxied http.RoundTripper
}

func (t testTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.proxied.RoundTrip(req)
}

func TestAuthenticateCustomTransport(t *testing.T) {
	simular.Activate()
	defer simular.DeactivateAndReset()

	simular.RegisterStubRequests(
		simular.NewStubRequest(
			http.MethodGet,
			"https://market.big-iot.org/accessToken?clientId=id&clientSecret=secret",
			simular.NewStringResponder(200, "1234abcd"),
			simular.WithHeader(
				&http.Header{
					"Accept": []string{"text/plain"},
				},
			),
		),
	)

	client := &http.Client{
		Timeout:   1 * time.Second,
		Transport: testTripper{proxied: http.DefaultTransport},
	}

	p, _ := NewProvider(
		"id",
		"secret",
		WithHTTPClient(client),
	)
	assert.Nil(t, p.Authenticate())
	assert.Equal(t, "1234abcd", p.accessToken)
	assert.Equal(t, DefaultMarketplaceURL, p.baseURL.String())

	assert.Nil(t, simular.AllStubsCalled())
}
