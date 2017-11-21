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
	assert.Equal(t, "id", p.id)
	assert.Equal(t, "secret", p.secret)
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
