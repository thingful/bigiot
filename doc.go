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

/*

Package bigiot is an attempt at porting the BIGIot client library from Java
to Go, adapting the library where appropriate to better fit Go idioms and
practices. This is very much a work in progress, so currently is a long way
from supporting the same range of functionality as the Java library.

Implemented functionality:

  * register an offering in the marketplace
	* delete or unregister an offering from the marketplace
	* reactivating offerings from the marketplace
	* validating tokens presented by offering subscribers

Planned functionality:
  * discovering an offering in the marketplace
  * subscribing to an offering

Here's an example of how you can create a local Provider client, and
authenticate with the marketplace. This assumes that you have already created
an account on the marketkplace and you have extracted from the marketplace
your providerID and secret.

	provider, err := bigiot.NewProvider(providerID, providerSecret)
	if err != nil {
		panic(err) // handle error properly
	}

	err = provider.Authenticate()
	if err != nil {
		panic(err) // handle error properly
	}

Then in order to register an offering a client would first create a
description of the offering.

	addOfferingInput := &bigiot.OfferingDescription{
		LocalID:  "ParkingOffering",
		Name:     "Demo Parking Offering",
		Category: "urn:big-iot:ParkingSpaces"
		Inputs: []bigiot.DataField{
			{
				Name:   "longitude",
				RdfURI: "schema:longitude",
			},
			{
				Name:   "latitude",
				RdfURI: "schema:latitude",
			},
		},
		Outputs: []bigiot.DataField{
			{
				Name:   "geoCoordinates",
				RdfURI: "schema:geoCoordinates",
			}
		},
		Endpoints: []bigiot.Endpoint{
			{
				URI:                 "https://example.com/parking",
				EndpointType:        bigiot.HTTPGet,
				AccessInterfaceType: bigiot.External,
			}
		},
		License: bigiot.OpenDataLicense,
		SpatialExtent: &bigiot.SpatialExtent{
			City: "Berlin",
			BoundingBox: &bigiot.BoundingBox{
				Location1: bigiot.Location{
					Lng: 2.33,
					Lat: 54.5,
				},
				Location2: bigiot.Location{
					Lng: 2.38,
					Lat: 54.53,
				},
			},
		},
		Price: bigiot.Price{
			Money: bigiot.Money{
				Amount:   0.01,
				Currency: bigiot.EUR,
			},
			PricingModel: bigiot.PerAccess,
		},
		Activation: &bigiot.Activation{
			Status:   true,
			Duration: 15 * time.Minute,
		},
	}

Then we can use the above description to register our offering on the marketplace:

	offering, err := provider.RegisterOffering(context.Background(), addOfferingInput)
	if err != nil {
		panic(err) // handle error properly
	}

To delete an offering we need to invoke the DeleteOffering method:

	deleteOfferingInput := &bigiot.DeleteOffering{
		ID: offering.ID,
	}

	err = provider.DeleteOffering(context.Background(), deleteOfferingInput)
	if err != nil {
		panic(err) // handle error properly
	}

To re-activate an offering (a provider will want to do this as they are
expected to continously re-activate to show an offering is still alive) you
can use the following method on the Provider:

	activateOfferingInput := &bigiot.ActivateOffering{
		ID: offering.ID,
		Duration: 15 * time.Minute,
	}

	err := provider.ActivateOffering(context.Background(), activateOfferingInput)
	if err != nil {
		panic(err) // handle error properly
	}

To validate incoming tokens presented by a consumer, we expose a
ValidateToken method. This takes as input a token string encoded via the
compact JWT serialization form, and returns either the ID of the offering
being requested, or an error should the incoming token be invalid.

	offeringID, err := provider.ValidateToken(tokenStr)
	if err != nil {
		panic(err) // handle error properly
	}
*/
package bigiot
