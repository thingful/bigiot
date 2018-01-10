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
	"bytes"
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/thingful/bigiot"
	"github.com/thingful/simular"
)

// mockClock is an implementation of the Clock interface for use in tests.
// Returns a canned time for "now".
type mockClock struct {
	t time.Time
}

// Now is our implementation of the Clock Now() function that in the real case
// returns the current time, but here we just return the canned time value set
// on the struct.
func (m mockClock) Now() time.Time {
	return m.t
}

func TestRegisterOffering(t *testing.T) {
	simular.Activate()
	defer simular.DeactivateAndReset()

	expirationTime := time.Unix(0, 1509983101577000000)

	simular.RegisterStubRequests(
		simular.NewStubRequest(
			http.MethodGet,
			"https://market.big-iot.org/accessToken?clientId=Provider&clientSecret=secret",
			simular.NewStringResponder(200, "1234abcd"),
		),
		simular.NewStubRequest(
			http.MethodPost,
			"https://market.big-iot.org/graphql",
			simular.NewStringResponder(200, `{"data": {"addOffering": {"id": "Organization-Provider-TestOffering", "activation": { "status": true, "expirationTime": 1509983101577}}}}`),
			simular.WithBody(
				bytes.NewBufferString(`{"query":"mutation addOffering { addOffering ( input: { id: \"Provider\", localId: \"TestOffering\", name: \"Test Offering\", activation: {status: true, expirationTime: 1509983101577} , rdfUri: \"\", inputData: [], outputData: [{name: \"value\", rdfUri: \"schema:random\"} ], endpoints: [{uri: \"https://example.com/random\", endpointType: HTTP_GET, accessInterfaceType: BIGIOT_LIB} ], license: OPEN_DATA_LICENSE, price: {money: {amount: 0.001, currency: EUR}, pricingModel: PER_ACCESS}, extent: {city: \"Berlin\"} } ) { id name activation { status expirationTime } } }"}`),
			),
		),
	)

	provider, err := bigiot.NewProvider("Provider", "secret")
	assert.Nil(t, err)

	err = provider.Authenticate()
	assert.Nil(t, err)

	offeringInput := &bigiot.OfferingDescription{
		LocalID: "TestOffering",
		Name:    "Test Offering",
		OutputData: []bigiot.DataField{
			{
				Name:   "value",
				RdfURI: "schema:random",
			},
		},
		Endpoints: []bigiot.Endpoint{
			{
				URI:                 "https://example.com/random",
				EndpointType:        bigiot.HTTPGet,
				AccessInterfaceType: bigiot.BIGIoTLib,
			},
		},
		License: bigiot.OpenDataLicense,
		Price: bigiot.Price{
			Money: bigiot.Money{
				Amount:   0.001,
				Currency: bigiot.EUR,
			},
			PricingModel: bigiot.PerAccess,
		},
		Extent: bigiot.Address{
			City: "Berlin",
		},
		Activation: bigiot.Activation{
			Status:         true,
			Duration:       10 * time.Minute,
			ExpirationTime: expirationTime,
		},
	}

	offering, err := provider.RegisterOffering(context.Background(), offeringInput)
	assert.Nil(t, err)
	assert.Equal(t, "Organization-Provider-TestOffering", offering.ID)
	assert.True(t, offering.Activation.Status)
	assert.Equal(t, expirationTime.UTC(), offering.Activation.ExpirationTime)
}

func TestRegisterOfferingWithDuration(t *testing.T) {
	simular.Activate()
	defer simular.DeactivateAndReset()

	now := time.Unix(0, 0)
	duration := 10 * time.Minute
	expirationTime := now.Add(duration)
	clock := mockClock{t: now}

	offeringInput := &bigiot.OfferingDescription{
		LocalID: "TestOffering",
		Name:    "Test Offering",
		OutputData: []bigiot.DataField{
			{
				Name:   "value",
				RdfURI: "schema:random",
			},
		},
		Endpoints: []bigiot.Endpoint{
			{
				URI:                 "https://example.com/random",
				EndpointType:        bigiot.HTTPGet,
				AccessInterfaceType: bigiot.BIGIoTLib,
			},
		},
		License: bigiot.OpenDataLicense,
		Price: bigiot.Price{
			Money: bigiot.Money{
				Amount:   0.001,
				Currency: bigiot.EUR,
			},
			PricingModel: bigiot.PerAccess,
		},
		Extent: bigiot.Address{
			City: "Berlin",
		},
		Activation: bigiot.Activation{
			Status:   true,
			Duration: duration,
		},
	}

	t.Run("with valid response", func(t *testing.T) {
		simular.RegisterStubRequests(
			simular.NewStubRequest(
				http.MethodGet,
				"https://market.big-iot.org/accessToken?clientId=Provider&clientSecret=secret",
				simular.NewStringResponder(200, "1234abcd"),
			),
			simular.NewStubRequest(
				http.MethodPost,
				"https://market.big-iot.org/graphql",
				simular.NewStringResponder(200, `{"data": {"addOffering": {"id": "Organization-Provider-TestOffering", "activation": { "status": true, "expirationTime": 600000}}}}`),
				simular.WithBody(
					bytes.NewBufferString(`{"query":"mutation addOffering { addOffering ( input: { id: \"Provider\", localId: \"TestOffering\", name: \"Test Offering\", activation: {status: true, expirationTime: 600000} , rdfUri: \"\", inputData: [], outputData: [{name: \"value\", rdfUri: \"schema:random\"} ], endpoints: [{uri: \"https://example.com/random\", endpointType: HTTP_GET, accessInterfaceType: BIGIOT_LIB} ], license: OPEN_DATA_LICENSE, price: {money: {amount: 0.001, currency: EUR}, pricingModel: PER_ACCESS}, extent: {city: \"Berlin\"} } ) { id name activation { status expirationTime } } }"}`),
				),
			),
		)

		provider, err := bigiot.NewProvider(
			"Provider",
			"secret",
			bigiot.WithClock(clock),
		)
		assert.Nil(t, err)

		err = provider.Authenticate()
		assert.Nil(t, err)

		offering, err := provider.RegisterOffering(context.Background(), offeringInput)
		assert.Nil(t, err)
		assert.Equal(t, "Organization-Provider-TestOffering", offering.ID)
		assert.True(t, offering.Activation.Status)
		assert.Equal(t, expirationTime.UTC(), offering.Activation.ExpirationTime)
	})

	t.Run("with error response", func(t *testing.T) {
		simular.RegisterStubRequests(
			simular.NewStubRequest(
				http.MethodGet,
				"https://market.big-iot.org/accessToken?clientId=Provider&clientSecret=secret",
				simular.NewStringResponder(200, "1234abcd"),
			),
			simular.NewStubRequest(
				http.MethodPost,
				"https://market.big-iot.org/graphql",
				simular.NewStringResponder(500, `error`),
				simular.WithBody(
					bytes.NewBufferString(`{"query":"mutation addOffering { addOffering ( input: { id: \"Provider\", localId: \"TestOffering\", name: \"Test Offering\", activation: {status: true, expirationTime: 600000} , rdfUri: \"\", inputData: [], outputData: [{name: \"value\", rdfUri: \"schema:random\"} ], endpoints: [{uri: \"https://example.com/random\", endpointType: HTTP_GET, accessInterfaceType: BIGIOT_LIB} ], license: OPEN_DATA_LICENSE, price: {money: {amount: 0.001, currency: EUR}, pricingModel: PER_ACCESS}, extent: {city: \"Berlin\"} } ) { id name activation { status expirationTime } } }"}`),
				),
			),
		)

		provider, err := bigiot.NewProvider(
			"Provider",
			"secret",
			bigiot.WithClock(clock),
		)
		assert.Nil(t, err)

		err = provider.Authenticate()
		assert.Nil(t, err)

		_, err = provider.RegisterOffering(context.Background(), offeringInput)
		assert.NotNil(t, err)
	})

	t.Run("with invalid json", func(t *testing.T) {
		simular.RegisterStubRequests(
			simular.NewStubRequest(
				http.MethodGet,
				"https://market.big-iot.org/accessToken?clientId=Provider&clientSecret=secret",
				simular.NewStringResponder(200, "1234abcd"),
			),
			simular.NewStubRequest(
				http.MethodPost,
				"https://market.big-iot.org/graphql",
				simular.NewStringResponder(200, `{"data": {"addOffering": {"id": "Organization-Provider-TestOffering"`),
				simular.WithBody(
					bytes.NewBufferString(`{"query":"mutation addOffering { addOffering ( input: { id: \"Provider\", localId: \"TestOffering\", name: \"Test Offering\", activation: {status: true, expirationTime: 600000} , rdfUri: \"\", inputData: [], outputData: [{name: \"value\", rdfUri: \"schema:random\"} ], endpoints: [{uri: \"https://example.com/random\", endpointType: HTTP_GET, accessInterfaceType: BIGIOT_LIB} ], license: OPEN_DATA_LICENSE, price: {money: {amount: 0.001, currency: EUR}, pricingModel: PER_ACCESS}, extent: {city: \"Berlin\"} } ) { id name activation { status expirationTime } } }"}`),
				),
			),
		)

		provider, err := bigiot.NewProvider(
			"Provider",
			"secret",
			bigiot.WithClock(clock),
		)
		assert.Nil(t, err)

		err = provider.Authenticate()
		assert.Nil(t, err)

		_, err = provider.RegisterOffering(context.Background(), offeringInput)
		assert.NotNil(t, err)
	})
}

func TestDeleteOffering(t *testing.T) {
	simular.Activate()
	defer simular.DeactivateAndReset()

	simular.RegisterStubRequests(
		simular.NewStubRequest(
			http.MethodGet,
			"https://market.big-iot.org/accessToken?clientId=Provider&clientSecret=secret",
			simular.NewStringResponder(200, "1234abcd"),
		),
		simular.NewStubRequest(
			http.MethodPost,
			"https://market.big-iot.org/graphql",
			simular.NewStringResponder(200, `{"data": {"deleteOffering": {"id": "Organization-Provider-TestOffering"}}}`),
			simular.WithBody(
				bytes.NewBufferString(`{"query":"mutation deleteOffering { deleteOffering ( input: { id: \"Organization-Provider-TestOffering\" } ) { id } }"}`),
			),
		),
	)

	provider, err := bigiot.NewProvider("Provider", "secret")
	assert.Nil(t, err)

	err = provider.Authenticate()
	assert.Nil(t, err)

	deleteOffering := &bigiot.DeleteOffering{
		ID: "Organization-Provider-TestOffering",
	}

	err = provider.DeleteOffering(context.Background(), deleteOffering)
	assert.Nil(t, err)
}

func TestDeleteOfferingError(t *testing.T) {
	simular.Activate()
	defer simular.DeactivateAndReset()

	simular.RegisterStubRequests(
		simular.NewStubRequest(
			http.MethodGet,
			"https://market.big-iot.org/accessToken?clientId=Provider&clientSecret=secret",
			simular.NewStringResponder(200, "1234abcd"),
		),
		simular.NewStubRequest(
			http.MethodPost,
			"https://market.big-iot.org/graphql",
			simular.NewStringResponder(500, `error`),
			simular.WithBody(
				bytes.NewBufferString(`{"query":"mutation deleteOffering { deleteOffering ( input: { id: \"Organization-Provider-TestOffering\" } ) { id } }"}`),
			),
		),
	)

	provider, err := bigiot.NewProvider("Provider", "secret")
	assert.Nil(t, err)

	err = provider.Authenticate()
	assert.Nil(t, err)

	deleteOffering := &bigiot.DeleteOffering{
		ID: "Organization-Provider-TestOffering",
	}

	err = provider.DeleteOffering(context.Background(), deleteOffering)
	assert.NotNil(t, err)
}

func TestActivatingOffering(t *testing.T) {
	now := time.Now()
	clock := mockClock{now}

	simular.Activate()
	defer simular.DeactivateAndReset()

	simular.RegisterStubRequests(
		simular.NewStubRequest(
			http.MethodGet,
			"https://market.big-iot.org/accessToken?clientId=Provider&clientSecret=secret",
			simular.NewStringResponder(200, "1234abcd"),
		),
		simular.NewStubRequest(
			http.MethodPost,
			"https://market.big-iot.org/graphql",
			simular.NewStringResponder(200, fmt.Sprintf(`{"data": {"activateOffering": {"id": "Organization-Provider-TestOffering", "activation": { "status": true, "expirationTime": %v}}}}`, bigiot.ToEpochMs(now.Add(10*time.Minute)))),
			simular.WithBody(
				bytes.NewBufferString(fmt.Sprintf(`{"query":"mutation activateOffering { activateOffering ( input: { id: \"Organization-Provider-TestOffering\", expirationTime: %v } ) { id activation { status expirationTime } } }"}`, bigiot.ToEpochMs(now.Add(10*time.Minute)))),
			),
		),
	)

	provider, err := bigiot.NewProvider(
		"Provider",
		"secret",
		bigiot.WithClock(clock),
	)
	assert.Nil(t, err)

	err = provider.Authenticate()
	assert.Nil(t, err)

	activateOffering := &bigiot.ActivateOffering{
		ID: "Organization-Provider-TestOffering",
	}

	_, err = provider.ActivateOffering(context.Background(), activateOffering)
	assert.Nil(t, err)
}

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
