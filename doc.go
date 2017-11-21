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

Planned functionality:
  * validating tokens presented by offering subscribers
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

Then in order to register an offering a client would first create a description of the offering.

	addOfferingInput := &bigiot.AddOffering{
		LocalID: "ParkingOffering",
		Name:    "Demo Parking Offering",
		InputData: []bigiot.DataField{
			{
				Name:   "longitude",
				RdfURI: "schema:longitude",
			},
			{
				Name:   "latitude",
				RdfURI: "schema:latitude",
			},
		},
		OutputData: []bigiot.DataField{
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
		Extent: bigiot.Address{
			City: "Berlin",
		},
		Price: bigiot.Price{
			Money: bigiot.Money{
				Amount:   0.01,
				Currency: bigiot.EUR,
			},
			PricingModel: bigiot.PerAccess,
		},
		Activation: bigiot.Activation{
			Status:         true,
			ExpirationTime: expirationTime,
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
*/
package bigiot
