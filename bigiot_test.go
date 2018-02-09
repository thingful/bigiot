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
	"github.com/thingful/bigiot/mocks"
	"github.com/thingful/simular"
)

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
				bytes.NewBufferString(`{"query":"mutation addOffering { addOffering ( input: { id: \"Provider\", localId: \"TestOffering\", name: \"Test Offering\", activation: { status: true, expirationTime: 1509983101577 }, rdfUri: \"urn:proposed:RandomValues\", outputs: [{ name: \"value\", rdfUri: \"schema:random\" }], endpoints: [{ uri: \"https://example.com/random\", endpointType: HTTP_GET, accessInterfaceType: BIGIOT_LIB }], license: OPEN_DATA_LICENSE, price: { money: { amount: 0.001, currency: EUR }, pricingModel: PER_ACCESS }, spatialExtent: { city: \"Berlin\", boundary: { l1: { lng: -2.25, lat: 54.53 }, l2: { lng: -2.26, lat: 54.96 } } } } ) { id name activation { status expirationTime } } }"}`),
			),
		),
	)

	provider, err := bigiot.NewProvider("Provider", "secret")
	assert.Nil(t, err)

	err = provider.Authenticate()
	assert.Nil(t, err)

	offeringInput := &bigiot.OfferingDescription{
		LocalID:  "TestOffering",
		Name:     "Test Offering",
		Category: "urn:proposed:RandomValues",
		Outputs: []bigiot.DataField{
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
		SpatialExtent: &bigiot.SpatialExtent{
			City: "Berlin",
			BoundingBox: &bigiot.BoundingBox{
				Location1: bigiot.Location{
					Lng: -2.25,
					Lat: 54.53,
				},
				Location2: bigiot.Location{
					Lng: -2.26,
					Lat: 54.96,
				},
			},
		},
		Activation: &bigiot.Activation{
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
	clock := mocks.Clock{T: now}

	offeringInput := &bigiot.OfferingDescription{
		LocalID:  "TestOffering",
		Name:     "Test Offering",
		Category: "urn:proposed:RandomValues",
		Outputs: []bigiot.DataField{
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
		SpatialExtent: &bigiot.SpatialExtent{
			City: "Berlin",
		},
		Activation: &bigiot.Activation{
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
					bytes.NewBufferString(`{"query":"mutation addOffering { addOffering ( input: { id: \"Provider\", localId: \"TestOffering\", name: \"Test Offering\", activation: { status: true, expirationTime: 600000 }, rdfUri: \"urn:proposed:RandomValues\", outputs: [{ name: \"value\", rdfUri: \"schema:random\" }], endpoints: [{ uri: \"https://example.com/random\", endpointType: HTTP_GET, accessInterfaceType: BIGIOT_LIB }], license: OPEN_DATA_LICENSE, price: { money: { amount: 0.001, currency: EUR }, pricingModel: PER_ACCESS }, spatialExtent: { city: \"Berlin\" } } ) { id name activation { status expirationTime } } }"}`),
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
				simular.NewStringResponder(400, `{"data":null,"errors":[{"message":"bad request"}]}`),
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
		assert.Regexp(t, "Error registering offering", err.Error())
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
			simular.NewStringResponder(400, `{"data":null,"errors":[{"message":"bad request"}]}`),
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
	assert.Equal(t, "Error deleting offering: bad request", err.Error())
}

func TestActivatingOffering(t *testing.T) {
	now := time.Now()
	clock := mocks.Clock{T: now}

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
