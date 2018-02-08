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

// +build integration

package bigiot_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/thingful/bigiot"
)

func Getenv(t *testing.T, key string) string {
	t.Helper()
	if val, ok := os.LookupEnv(key); ok {
		return val
	}

	t.Fatalf("Missing required env variable: %s", key)
	return ""
}

func TestProviderIntegration(t *testing.T) {
	providerID := Getenv(t, "PROVIDER_ID")
	providerSecret := Getenv(t, "PROVIDER_SECRET")

	provider, err := bigiot.NewProvider(providerID, providerSecret)
	assert.Nil(t, err)

	err = provider.Authenticate()
	assert.Nil(t, err)

	expirationTime := time.Now().Add(60 * time.Minute)

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

	_, err = provider.RegisterOffering(context.Background(), offeringInput)
	assert.Nil(t, err)
	//assert.NotNil(t, offering.ID)

	//deleteOffering := &bigiot.DeleteOffering{
	//	ID: offering.ID,
	//}

	//err = provider.DeleteOffering(context.Background(), deleteOffering)
	//assert.Nil(t, err)
}
