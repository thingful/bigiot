package bigiot_test

import (
	"net/http"
	"testing"

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
			"https://market.big-iot.org/accessToken",
			httpmock.NewStringResponder(200, "1234abcd"),
		).WithHeader(
			&http.Header{
				"User-Agent": []string{"bigiot-go"},
			},
		),
	)

	p, _ := bigiot.NewProvider("id", "secret")
	assert.Nil(t, p.Authenticate())
	assert.Equal(t, "1234abcd", p.AccessToken)
	assert.Equal(t, bigiot.DefaultMarketplaceURL, p.BaseURL.String())

	assert.Nil(t, httpmock.AllStubsCalled())
}
