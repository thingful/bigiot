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

package bigiot

// EndpointType represents the type of an Endpoint accessible via the BIGIoT
// Marketplace.
type EndpointType string

const (
	// HTTPGet is a const value representing an endpoint type accessible via HTTP
	// GET
	HTTPGet EndpointType = "HTTP_GET"

	// HTTPPost is a const value representing an endpoint type accessible via HTTP
	// POST
	HTTPPost EndpointType = "HTTP_POST"

	// WebSocket is a const value representing an endpoint type accessible via
	// WebSockets
	WebSocket EndpointType = "WEBSOCKET"
)

// String is an implementation of the Stringer interface for EndpointType
// instances.
func (e EndpointType) String() string {
	return string(e)
}

// AccessInterfaceType is a type used to represent the type of an access
// interface. This can be one of BIGIOT_LIB or EXTERNAL.
type AccessInterfaceType string

const (
	// BIGIoTLib is the const value representing an interface type of BIGIOT_LIB
	BIGIoTLib AccessInterfaceType = "BIGIOT_LIB"

	// External is the const value representing an interface type of EXTERNAL
	External AccessInterfaceType = "EXTERNAL"
)

// String is an implementation of Stringer for our AccessInterfaceType type.
func (a AccessInterfaceType) String() string {
	return string(a)
}

// License is a type alias for string used to represent the license being
// applied to an offering.
type License string

const (
	// CreativeCommons is a License instance representing the Creative Commons License
	CreativeCommons License = "CREATIVE_COMMONS"

	// OpenDataLicense is a License instance representing an open data license
	OpenDataLicense License = "OPEN_DATA_LICENSE"

	// NonCommercialDataLicense is a License instance representing a non-commercial
	// data license
	NonCommercialDataLicense License = "NON_COMMERCIAL_DATA_LICENSE"
)

// String is an implementation of Stringer for our License type.
func (l License) String() string {
	return string(l)
}

// PricingModel is a type alias for string used to represent pricing models to
// be applied to BIGIoT offerings.
type PricingModel string

const (
	// Free is a const used to represent a free pricing model.
	Free PricingModel = "FREE"

	// PerMonth is a const used to represent a per month based pricing model.
	PerMonth PricingModel = "PER_MONTH"

	// PerAccess is a const used to represent a per access based pricing model.
	PerAccess PricingModel = "PER_ACCESS"

	// PerByte is a const used to represent a per byte based pricing model.
	PerByte PricingModel = "PER_BYTE"
)

// String is an implementation of Stringer for our PricingModel type.
func (p PricingModel) String() string {
	return string(p)
}

// Currency is a type alias for string used to represent currencies
type Currency string

const (
	// EUR is a currency instance representing the Euro currency
	EUR Currency = "EUR"
)

// String is an implementation of Stringer for our Currency type.
func (c Currency) String() string {
	return string(c)
}
