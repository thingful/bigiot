package bigiot_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/thingful/bigiot"
	"github.com/thingful/simular"
)

func TestProviderConstructorInvalidMarketplace(t *testing.T) {
	_, err := bigiot.NewProvider(
		"id",
		"secret",
		bigiot.WithMarketplace("http ://market-dev.big-iot.org"),
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
					"User-Agent": []string{"bigiot/" + bigiot.Version + " (https://github.com/thingful/bigiot)"},
					"Accept":     []string{"text/plain"},
				},
			),
		),
	)

	p, _ := bigiot.NewProvider("id", "secret")
	assert.Nil(t, p.Authenticate())

	assert.Nil(t, simular.AllStubsCalled())
}
func TestAuthenticateDifferentMarketplace(t *testing.T) {
	simular.Activate()
	defer simular.DeactivateAndReset()

	simular.RegisterStubRequests(
		simular.NewStubRequest(
			http.MethodGet,
			"https://market-dev.big-iot.org/accessToken?clientId=id&clientSecret=secret",
			simular.NewStringResponder(200, "1234abcd"),
			simular.WithHeader(
				&http.Header{
					"User-Agent": []string{"bigiot/" + bigiot.Version + " (https://github.com/thingful/bigiot)"},
					"Accept":     []string{"text/plain"},
				},
			),
		),
	)

	p, _ := bigiot.NewProvider(
		"id",
		"secret",
		bigiot.WithMarketplace("https://market-dev.big-iot.org"),
	)

	assert.Nil(t, p.Authenticate())

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

	p, _ := bigiot.NewProvider("id", "secret")
	err := p.Authenticate()
	assert.Equal(t, bigiot.ErrUnexpectedResponse, err)
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

	p, _ := bigiot.NewProvider(
		"id",
		"secret",
		bigiot.WithUserAgent("foo"),
	)
	assert.Nil(t, p.Authenticate())

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

	p, _ := bigiot.NewProvider(
		"id",
		"secret",
		bigiot.WithHTTPClient(client),
	)
	assert.Nil(t, p.Authenticate())

	assert.Nil(t, simular.AllStubsCalled())
}

func TestValidateTokenTimes(t *testing.T) {
	now := time.Date(2018, 1, 1, 9, 4, 0, 0, time.UTC)

	tokenStr := `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTQ3OTc1MDAsIm5iZiI6MTUxNDc5NzQ0MCwic3ViIjoiQ29uc3VtZXItUXVlcnk9PVByb3ZpZGVyLU9mZmVyaW5nIn0.KQAhL8d3ynvX3YmGIrScq11p-q_Y61aQybOS9lo5H7c`

	times := []struct {
		label string
		time  time.Time
		valid bool
	}{
		{"now", now, true},
		{"late", now.Add(time.Duration(60) * time.Second), true},
		{"early", now.Add(time.Duration(-60) * time.Second), true},
		{"too late", now.Add(time.Duration(121) * time.Second), false},
		{"too early", now.Add(time.Duration(-61) * time.Second), false},
	}

	for _, tm := range times {
		t.Run(tm.label, func(t *testing.T) {
			clock := mockClock{t: tm.time}

			p, err := bigiot.NewProvider(
				"id",
				"CF72ABfRTqy1FQS1zBaevw==",
				bigiot.WithClock(clock))

			assert.Nil(t, err)

			if tm.valid {
				assert.Nil(t, p.ValidateToken(tokenStr))
			} else {
				err = p.ValidateToken(tokenStr)
				fmt.Println("ERROR", err)
				assert.NotNil(t, err)
			}
		})
	}
}

func TestValidateTokenData(t *testing.T) {
	now := time.Date(2018, 1, 1, 9, 4, 0, 0, time.UTC)
	clock := mockClock{t: now}

	testcases := []struct {
		label    string
		secret   string
		tokenStr string
		valid    bool
	}{
		{
			"valid",
			"CF72ABfRTqy1FQS1zBaevw==",
			`eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTQ3OTc1MDAsIm5iZiI6MTUxNDc5NzQ0MCwic3ViIjoiQ29uc3VtZXItUXVlcnk9PVByb3ZpZGVyLU9mZmVyaW5nIn0.KQAhL8d3ynvX3YmGIrScq11p-q_Y61aQybOS9lo5H7c`,
			true,
		},
		{
			"invalid secret",
			"CF72ABfRTqy1FOS1zBaevw==",
			`eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTQ3OTc1MDAsIm5iZiI6MTUxNDc5NzQ0MCwic3ViIjoiQ29uc3VtZXItUXVlcnk9PVByb3ZpZGVyLU9mZmVyaW5nIn0.KQAhL8d3ynvX3YmGIrScq11p-q_Y61aQybOS9lo5H7c`,
			false,
		},
		{
			"invalid character in token",
			"CF72ABfRTqy1FQS1zBaevw==",
			`eyJhbGciOiJIUZI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTQ3OTc1MDAsIm5iZiI6MTUxNDc5NzQ0MCwic3ViIjoiQ29uc3VtZXItUXVlcnk9PVByb3ZpZGVyLU9mZmVyaW5nIn0.KQAhL8d3ynvX3YmGIrScq11p-q_Y61aQybOS9lo5H7c`,
			false,
		},
		{
			"missing token part",
			"CF72ABfRTqy1FQS1zBaevw==",
			`eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTQ3OTc1MDAsIm5iZiI6MTUxNDc5NzQ0MCwic3ViIjoiQ29uc3VtZXItUXVlcnk9PVByb3ZpZGVyLU9mZmVyaW5nIn0`,
			false,
		},
		{
			"non base64 key",
			"hello",
			`eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTQ3OTc1MDAsIm5iZiI6MTUxNDc5NzQ0MCwic3ViIjoiQ29uc3VtZXItUXVlcnk9PVByb3ZpZGVyLU9mZmVyaW5nIn0.KQAhL8d3ynvX3YmGIrScq11p-q_Y61aQybOS9lo5H7c`,
			false,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.label, func(t *testing.T) {
			p, err := bigiot.NewProvider(
				"id",
				testcase.secret,
				bigiot.WithClock(clock))
			assert.Nil(t, err)

			if testcase.valid {
				assert.Nil(t, p.ValidateToken(testcase.tokenStr))
			} else {
				assert.NotNil(t, p.ValidateToken(testcase.tokenStr))
			}
		})
	}

}
