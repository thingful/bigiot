package bigiot_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/thingful/bigiot"
	"github.com/thingful/bigiot/mocks"
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
			simular.NewStringResponder(403, "ClientDoesNotExist: id"),
		),
	)

	p, _ := bigiot.NewProvider("id", "secret")
	err := p.Authenticate()
	assert.Equal(t, "ClientDoesNotExist: id", err.Error())
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
	offeringID := "Provider-Offering"
	now := time.Date(2018, 1, 1, 9, 4, 0, 0, time.UTC)

	// encoded token string for the above date, that expires after 1 minute
	tokenStr := `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTQ3OTc1MDAsIm5iZiI6MTUxNDc5NzQ0MCwic3ViIjoiQ29uc3VtZXItUXVlcnk9PVByb3ZpZGVyLU9mZmVyaW5nIiwic3Vic2NyaWJhYmxlSWQiOiJQcm92aWRlci1PZmZlcmluZyIsInN1YnNjcmliZXJJZCI6IkNvbnN1bWVyLVF1ZXJ5In0.US13MsSFVLvl7NMndSSKERhIz5d6K2hXH4w06_ZVkZw`

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
			clock := mocks.Clock{T: tm.time}

			p, err := bigiot.NewProvider(
				"id",
				"CF72ABfRTqy1FQS1zBaevw==",
				bigiot.WithClock(clock))

			assert.Nil(t, err)

			if tm.valid {
				gotID, err := p.ValidateToken(tokenStr)
				assert.Nil(t, err)
				assert.Equal(t, offeringID, gotID)
			} else {
				_, err = p.ValidateToken(tokenStr)
				fmt.Println("ERROR", err)
				assert.NotNil(t, err)
			}
		})
	}
}

func TestValidateTokenData(t *testing.T) {
	offeringID := "Provider-Offering"
	now := time.Date(2018, 1, 1, 9, 4, 0, 0, time.UTC)
	clock := mocks.Clock{T: now}

	testcases := []struct {
		label    string
		secret   string
		tokenStr string
		valid    bool
	}{
		{
			"valid",
			"CF72ABfRTqy1FQS1zBaevw==",
			`eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTQ3OTc1MDAsIm5iZiI6MTUxNDc5NzQ0MCwic3ViIjoiQ29uc3VtZXItUXVlcnk9PVByb3ZpZGVyLU9mZmVyaW5nIiwic3Vic2NyaWJhYmxlSWQiOiJQcm92aWRlci1PZmZlcmluZyIsInN1YnNjcmliZXJJZCI6IkNvbnN1bWVyLVF1ZXJ5In0.US13MsSFVLvl7NMndSSKERhIz5d6K2hXH4w06_ZVkZw`,
			true,
		},
		{
			"invalid secret",
			"CF72ABfRTqy1FOS1zBaevw==",
			`eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTQ3OTc1MDAsIm5iZiI6MTUxNDc5NzQ0MCwic3ViIjoiQ29uc3VtZXItUXVlcnk9PVByb3ZpZGVyLU9mZmVyaW5nIiwic3Vic2NyaWJhYmxlSWQiOiJQcm92aWRlci1PZmZlcmluZyIsInN1YnNjcmliZXJJZCI6IkNvbnN1bWVyLVF1ZXJ5In0.US13MsSFVLvl7NMndSSKERhIz5d6K2hXH4w06_ZVkZw`,
			false,
		},
		{
			"invalid character in token",
			"CF72ABfRTqy1FQS1zBaevw==",
			`eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.fyJleHAiOjE1MTQ3OTc1MDAsIm5iZiI6MTUxNDc5NzQ0MCwic3ViIjoiQ29uc3VtZXItUXVlcnk9PVByb3ZpZGVyLU9mZmVyaW5nIiwic3Vic2NyaWJhYmxlSWQiOiJQcm92aWRlci1PZmZlcmluZyIsInN1YnNjcmliZXJJZCI6IkNvbnN1bWVyLVF1ZXJ5In0.US13MsSFVLvl7NMndSSKERhIz5d6K2hXH4w06_ZVkZw`,
			false,
		},
		{
			"missing token part",
			"CF72ABfRTqy1FQS1zBaevw==",
			`eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTQ3OTc1MDAsIm5iZiI6MTUxNDc5NzQ0MCwic3ViIjoiQ29uc3VtZXItUXVlcnk9PVByb3ZpZGVyLU9mZmVyaW5nIiwic3Vic2NyaWJhYmxlSWQiOiJQcm92aWRlci1PZmZlcmluZyIsInN1YnNjcmliZXJJZCI6IkNvbnN1bWVyLVF1ZXJ5In0`,
			false,
		},
		{
			"non base64 key",
			"hello",
			`eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTQ3OTc1MDAsIm5iZiI6MTUxNDc5NzQ0MCwic3ViIjoiQ29uc3VtZXItUXVlcnk9PVByb3ZpZGVyLU9mZmVyaW5nIiwic3Vic2NyaWJhYmxlSWQiOiJQcm92aWRlci1PZmZlcmluZyIsInN1YnNjcmliZXJJZCI6IkNvbnN1bWVyLVF1ZXJ5In0.US13MsSFVLvl7NMndSSKERhIz5d6K2hXH4w06_ZVkZw`,
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
				gotID, err := p.ValidateToken(testcase.tokenStr)
				assert.Nil(t, err)
				assert.Equal(t, offeringID, gotID)
			} else {
				_, err := p.ValidateToken(testcase.tokenStr)
				assert.NotNil(t, err)
			}
		})
	}

}
