package bigiot

import (
	"net/http"
	"testing"
	"time"

	"github.com/smulube/httpmock"
	"github.com/stretchr/testify/assert"
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
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterStubRequest(
		httpmock.NewStubRequest(
			http.MethodGet,
			"https://market.big-iot.org/accessToken?clientId=id&clientSecret=secret",
			httpmock.NewStringResponder(200, "1234abcd"),
		).WithHeader(
			&http.Header{
				"User-Agent": []string{"bigiot/" + Version + " (https://github.com/thingful/bigiot)"},
				"Accept":     []string{"text/plain"},
			},
		),
	)

	p, _ := NewProvider("id", "secret")
	assert.Nil(t, p.Authenticate())
	assert.Equal(t, "1234abcd", p.accessToken)
	assert.Equal(t, DefaultMarketplaceURL, p.baseURL.String())

	assert.Nil(t, httpmock.AllStubsCalled())
}

func TestAuthenticateUnexpectedResponse(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterStubRequest(
		httpmock.NewStubRequest(
			http.MethodGet,
			"https://market.big-iot.org/accessToken?clientId=id&clientSecret=secret",
			httpmock.NewStringResponder(403, "Forbidden"),
		),
	)

	p, _ := NewProvider("id", "secret")
	err := p.Authenticate()
	assert.Equal(t, ErrUnexpectedResponse, err)
	assert.Equal(t, "Unexpected HTTP response code", err.Error())
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

	p, _ := NewProvider(
		"id",
		"secret",
		WithUserAgent("foo"),
	)
	assert.Nil(t, p.Authenticate())
	assert.Equal(t, "1234abcd", p.accessToken)
	assert.Equal(t, DefaultMarketplaceURL, p.baseURL.String())

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

	p, _ := NewProvider(
		"id",
		"secret",
		WithHTTPClient(client),
	)
	assert.Nil(t, p.Authenticate())
	assert.Equal(t, "1234abcd", p.accessToken)
	assert.Equal(t, DefaultMarketplaceURL, p.baseURL.String())

	assert.Nil(t, httpmock.AllStubsCalled())
}

//func TestRegisterOffering(t *testing.T) {
//	id := "Thingful-Temperature_Service"
//	secret := "mT1WkU1DQ56SX61l4mtIvg=="
//
//	provider, err := bigiot.NewProvider(id, secret, bigiot.WithMarketplace("https://market-dev.big-iot.org"))
//	assert.Nil(t, err)
//
//	offeringDescription := &bigiot.OfferingDescription{
//		Name:    "Simple Weather",
//		RdfType: bigiot.RdfType("bigiot:Weather"),
//		Endpoints: []bigiot.Endpoint{
//			{
//				URI:          "http://example.com/weather",
//				EndpointType: bigiot.HTTPGet,
//			},
//		},
//		InputData: []bigiot.DataField{
//			{
//				Name:    "longitude",
//				RdfType: bigiot.RdfType("schema:longitude"),
//			},
//			{
//				Name:    "latitude",
//				RdfType: bigiot.RdfType("schema:latitude"),
//			},
//		},
//		OutputData: []bigiot.DataField{
//			{
//				Name:    "temperature",
//				RdfType: bigiot.RdfType("schema:airTemperatureValue"),
//			},
//		},
//		Extent: bigiot.Extent{
//			City: "Edinburgh",
//		},
//	}
//
//	assert.NotNil(t, offeringDescription)
//
//	_, err = provider.RegisterOffering(offeringDescription)
//	assert.Nil(t, err)
//	//assert.NotNil(t, offering)
//}
//
//func TestOffering(t *testing.T) {
//	httpmock.Activate()
//	defer httpmock.DeactivateAndReset()
//
//	httpmock.RegisterStubRequest(
//		httpmock.NewStubRequest(
//			http.MethodGet,
//			"https://market.big-iot.org/accessToken?clientId=id&clientSecret=secret",
//			httpmock.NewStringResponder(200, "1234abcd"),
//		).WithHeader(
//			&http.Header{
//				"Accept": []string{"text/plain"},
//			},
//		),
//	)
//
//	httpmock.RegisterStubRequest(
//		httpmock.NewStubRequest(
//			http.MethodPost,
//			"https://market.big-iot.org/graphql",
//			httpmock.NewStringResponder(200, `{
//				"data": {
//					"offering": {
//						"id": "offeringID",
//						"name": "offering name"
//					}
//				}
//			}`),
//		).WithHeader(
//			&http.Header{
//				"Authorization": []string{"Bearer 1234abcd"},
//			},
//		),
//	)
//
//	provider, err := bigiot.NewProvider("id", "secret")
//	assert.Nil(t, err)
//
//	err = provider.Authenticate()
//	assert.Nil(t, err)
//
//	offering, err := provider.Offering("offeringID")
//	assert.Nil(t, err)
//	assert.NotNil(t, offering)
//	assert.Equal(t, "offeringID", offering.ID)
//}
//
