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

import (
	"bytes"
	"encoding/json"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

// OfferingDescription is the type used to register an offering with the
// marketplace. It contains information about the offerings inputs and outputs,
// its endpoints, license and price. In addition this is how offerings specify
// that they are active.
type OfferingDescription struct {
	providerID    string
	LocalID       string
	Name          string
	Category      string
	Inputs        []DataField
	Outputs       []DataField
	Endpoints     []Endpoint
	SpatialExtent *SpatialExtent
	License       License
	Price         Price
	Activation    *Activation
}

// serialize attempts to serialize it into the string form that the marketplace
// accepts as input to register an offering in the marketplace. Currently this
// implemented by manually building up the query using a bytes.Buffer as the
// existing Go graphql libraries didn't seem able to communicate with the
// marketplace.
func (o *OfferingDescription) serialize(clock Clock) string {
	var buf bytes.Buffer

	buf.WriteString(`mutation addOffering { addOffering ( input: { id: "`)
	buf.WriteString(o.providerID)
	buf.WriteString(`", localId: "`)
	buf.WriteString(o.LocalID)
	buf.WriteString(`", name: "`)
	buf.WriteString(o.Name)

	if o.Activation != nil {
		buf.WriteString(`", activation: `)
		buf.WriteString(o.Activation.serialize(clock))
	}

	buf.WriteString(`, rdfUri: "`)
	buf.WriteString(o.Category)
	buf.WriteString(`"`)

	if len(o.Inputs) > 0 {
		buf.WriteString(`, inputs: [`)
		for i, input := range o.Inputs {
			buf.WriteString(input.serialize(clock))
			if i < len(o.Inputs)-1 {
				buf.WriteString(`, `)
			}
		}
		buf.WriteString(`]`)
	}

	if len(o.Outputs) > 0 {
		buf.WriteString(`, outputs: [`)
		for i, output := range o.Outputs {
			buf.WriteString(output.serialize(clock))
			if i < len(o.Outputs)-1 {
				buf.WriteString(`, `)
			}
		}
		buf.WriteString(`]`)
	}

	if len(o.Endpoints) > 0 {
		buf.WriteString(`, endpoints: [`)
		for i, endpoint := range o.Endpoints {
			buf.WriteString(endpoint.serialize(clock))
			if i < len(o.Endpoints)-1 {
				buf.WriteString(`, `)
			}
		}
		buf.WriteString(`]`)
	}

	// add license
	buf.WriteString(`, license: `)
	buf.WriteString(o.License.String())

	// add price
	buf.WriteString(`, price: `)
	buf.WriteString(o.Price.serialize(clock))

	if o.SpatialExtent != nil {
		// add extent
		buf.WriteString(`, spatialExtent: `)
		buf.WriteString(o.SpatialExtent.serialize(clock))
	}

	buf.WriteString(` } )`)

	// desired returned output
	buf.WriteString(` { id name activation { status expirationTime } } }`)

	return buf.String()
}

// DataField captures information about an offering's inputs or outputs. Used
// when creating an offering.
type DataField struct {
	Name   string
	RdfURI string
}

// serialize is our implementation of Serializable for DataField. Serializes
// into a form that the marketplace understands.
func (d *DataField) serialize(clock Clock) string {
	var buf bytes.Buffer

	buf.WriteString(`{ name: "`)
	buf.WriteString(d.Name)
	buf.WriteString(`", rdfUri: "`)
	buf.WriteString(d.RdfURI)
	buf.WriteString(`" }`)

	return buf.String()
}

// Endpoint captures information about the endpoint of an offering.
type Endpoint struct {
	EndpointType        EndpointType
	URI                 string
	AccessInterfaceType AccessInterfaceType
}

// serialize is Endpoint's implementation of our Serializable interface
func (e *Endpoint) serialize(clock Clock) string {
	var buf bytes.Buffer

	buf.WriteString(`{ uri: "`)
	buf.WriteString(e.URI)
	buf.WriteString(`", endpointType: `)
	buf.WriteString(e.EndpointType.String())
	buf.WriteString(`, accessInterfaceType: `)
	buf.WriteString(e.AccessInterfaceType.String())
	buf.WriteString(` }`)

	return buf.String()
}

// SpatialExtent is how the BIG IoT marketplace defines geographical constraints when
// registering an offering.
type SpatialExtent struct {
	City        string
	BoundingBox *BoundingBox
}

// serialize is our implementation of serializable - to convert into BIG IoT
// graphql form.
func (a *SpatialExtent) serialize(clock Clock) string {
	var buf bytes.Buffer

	buf.WriteString(`{ city: "`)
	buf.WriteString(a.City)
	buf.WriteString(`"`)

	if a.BoundingBox != nil {
		buf.WriteString(`, boundary: `)
		buf.WriteString(a.BoundingBox.serialize(clock))
	}

	buf.WriteString(` }`)

	return buf.String()
}

// BoundingBox is used to represent a geographical bounding box within which an
// offering provides data. It contains two locations representing opposite
// corners of a geospatial box.
type BoundingBox struct {
	Location1 Location
	Location2 Location
}

// serialize is our implementation of serializable - to convert into BIG IoT
// graphql form.
func (b *BoundingBox) serialize(clock Clock) string {
	var buf bytes.Buffer

	buf.WriteString(`{ l1: `)
	buf.WriteString(b.Location1.serialize(clock))
	buf.WriteString(`, l2: `)
	buf.WriteString(b.Location2.serialize(clock))
	buf.WriteString(` }`)

	return buf.String()
}

// Location is used to represent a geographic location expressed as a decimal
// lng/lat pair.
type Location struct {
	Lng float64
	Lat float64
}

// serialize is our implementation of the serializable interface for BIG IoT graphql
func (l *Location) serialize(clock Clock) string {
	var buf bytes.Buffer
	buf.WriteString(`{ lng: `)
	buf.WriteString(strconv.FormatFloat(l.Lng, 'f', -1, 64))
	buf.WriteString(`, lat: `)
	buf.WriteString(strconv.FormatFloat(l.Lat, 'f', -1, 64))
	buf.WriteString(` }`)

	return buf.String()
}

// Price captures information about the pricing of an offering.
type Price struct {
	PricingModel PricingModel
	Money        Money
}

// serialize is our implementation of Serializable for Price objects.
func (p *Price) serialize(clock Clock) string {
	var buf bytes.Buffer

	buf.WriteString(`{ money: `)
	buf.WriteString(p.Money.serialize(clock))
	buf.WriteString(`, pricingModel: `)
	buf.WriteString(p.PricingModel.String())
	buf.WriteString(` }`)

	return buf.String()
}

// Money is used to capture price information for the offering. Note we aren't
// using precise numeric types here so this is not suitable for precision
// calculations.
type Money struct {
	Amount   float64 // TODO: look at more precise numeric type here
	Currency Currency
}

// serialize is our implementation of Serializable for Money objects.
func (m *Money) serialize(clock Clock) string {
	var buf bytes.Buffer

	buf.WriteString(`{ amount: `)
	buf.WriteString(strconv.FormatFloat(m.Amount, 'g', -1, 64))
	buf.WriteString(`, currency: `)
	buf.WriteString(m.Currency.String())
	buf.WriteString(` }`)

	return buf.String()
}

// Activation represents an activation of a resource. This comprises a boolean
// flag, and an expiration time. If the flag is set to true and the expiration
// time is in the future, then the offering is active; otherwise it is inactive.
type Activation struct {
	Status         bool          `json:"status"`
	ExpirationTime time.Time     `json:"-"`
	Duration       time.Duration `json:"-"`
}

// serialize converts our Activation into the structure required to send to the
// marketplace
func (a *Activation) serialize(clock Clock) string {
	var (
		buf            bytes.Buffer
		expirationTime time.Time
	)

	buf.WriteString(`{ status: `)
	buf.WriteString(strconv.FormatBool(a.Status))
	buf.WriteString(`, expirationTime: `)
	if a.ExpirationTime.IsZero() {
		if a.Duration == 0 {
			expirationTime = clock.Now().Add(DefaultActivationDuration)
		} else {
			expirationTime = clock.Now().Add(a.Duration)
		}
	} else {
		expirationTime = a.ExpirationTime
	}
	buf.WriteString(ToEpochMs(expirationTime))
	buf.WriteString(` }`)

	return buf.String()
}

// UnmarshalJSON is an implementation of the json Unmarshaler interface. We add
// a custom implementation to handle converting timestamps from epoch
// milliseconds into golang time.Time objects.
func (a *Activation) UnmarshalJSON(b []byte) error {
	// create anonymous struct for unmarshalling
	d := struct {
		Status         bool  `json:"status"`
		ExpirationTime int64 `json:"expirationTime"`
	}{}

	err := json.Unmarshal(b, &d)
	if err != nil {
		return errors.Wrap(err, "error unmarshalling activation type")
	}

	a.Status = d.Status
	a.ExpirationTime = FromEpochMs(d.ExpirationTime)

	return nil
}

// Offering is an output type used when returning information about an offering.
// This can happen either after creating an offering or if we get information on
// an offering from the marketplace.
type Offering struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Activation Activation `json:"activation"`
}

// DeleteOffering is an input type used to delete or unregister an offering.
type DeleteOffering struct {
	ID string
}

// serialize is our implementation of Serializable for DeleteOffering objects.
func (d *DeleteOffering) serialize(clock Clock) string {
	var buf bytes.Buffer

	buf.WriteString(`mutation deleteOffering { deleteOffering ( input: { id: "`)
	buf.WriteString(d.ID)
	buf.WriteString(`" } ) { id } }`)

	return buf.String()
}

// ActivateOffering is an input type used to reactivate an existing offering.
type ActivateOffering struct {
	ID             string
	ExpirationTime time.Time
	Duration       time.Duration
}

// serialize is our implementation of the serializable interface
func (a *ActivateOffering) serialize(clock Clock) string {
	var (
		buf            bytes.Buffer
		expirationTime time.Time
	)

	buf.WriteString(`mutation activateOffering { activateOffering ( input: { id: "`)
	buf.WriteString(a.ID)
	buf.WriteString(`", expirationTime: `)
	if a.ExpirationTime.IsZero() {
		if a.Duration == 0 {
			expirationTime = clock.Now().Add(DefaultActivationDuration)
		} else {
			expirationTime = clock.Now().Add(a.Duration)
		}
	} else {
		expirationTime = a.ExpirationTime
	}
	buf.WriteString(ToEpochMs(expirationTime))
	buf.WriteString(` } ) { id activation { status expirationTime } } }`)

	return buf.String()
}
