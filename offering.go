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
)

// OfferingInput is the type used to register an offering with the marketplace.
// It contains information about the offerings inputs and outputs, its
// endpoints, license and price. In addition this is how offerings specify that
// they are active.
type OfferingInput struct {
	providerID string
	LocalID    string
	Name       string
	RdfURI     string
	InputData  []DataField
	OutputData []DataField
	Endpoints  []Endpoint
	Extent     Address
	License    License
	Price      Price
	Activation Activation
}

// Serialize attempts to serialize it into the string form that the marketplace
// accepts as input to register an offering in the marketplace. Currently this
// implemented by manually building up the query using a bytes.Buffer as the
// existing Go graphql libraries didn't seem able to communicate with the
// marketplace.
func (o *OfferingInput) Serialize() string {
	var buf bytes.Buffer

	buf.WriteString(`mutation addOffering { addOffering ( input: { id: "`)
	buf.WriteString(o.providerID)
	buf.WriteString(`", localId: "`)
	buf.WriteString(o.LocalID)
	buf.WriteString(`", name: "`)
	buf.WriteString(o.Name)
	buf.WriteString(`", activation: `)
	buf.WriteString(o.Activation.Serialize())
	buf.WriteString(`, rdfUri: "`)
	buf.WriteString(o.RdfURI)
	buf.WriteString(`", inputData: [`)
	// serialized inputData goes here
	for _, input := range o.InputData {
		buf.WriteString(input.Serialize())
		buf.WriteString(" ")
	}
	buf.WriteString(`], outputData: [`)
	// serialized outputData goes here
	for _, output := range o.OutputData {
		buf.WriteString(output.Serialize())
		buf.WriteString(" ")
	}
	buf.WriteString(`], endpoints: [`)
	// serialized endpoints
	for _, endpoint := range o.Endpoints {
		buf.WriteString(endpoint.Serialize())
		buf.WriteString(" ")
	}
	buf.WriteString(`], license: `)
	buf.WriteString(o.License.String())
	buf.WriteString(`, price: `)
	// serialized price
	buf.WriteString(o.Price.Serialize())
	buf.WriteString(`, extent: `)
	// serialized address
	buf.WriteString(o.Extent.Serialize())
	buf.WriteString(` } ) `)

	// desired returned output
	buf.WriteString(`{ id name activation { status expirationTime } } }`)

	return buf.String()
}

type Offering struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Activation Activation `json:"activation"`
}

type DataField struct {
	Name   string
	RdfURI string
}

func (d *DataField) Serialize() string {
	var buf bytes.Buffer

	buf.WriteString(`{name: "`)
	buf.WriteString(d.Name)
	buf.WriteString(`", rdfUri: "`)
	buf.WriteString(d.RdfURI)
	buf.WriteString(`"}`)

	return buf.String()
}

type Endpoint struct {
	EndpointType        EndpointType
	URI                 string
	AccessInterfaceType AccessInterfaceType
}

func (e *Endpoint) Serialize() string {
	var buf bytes.Buffer

	buf.WriteString(`{uri: "`)
	buf.WriteString(e.URI)
	buf.WriteString(`", endpointType: `)
	buf.WriteString(e.EndpointType.String())
	buf.WriteString(`, accessInterfaceType: `)
	buf.WriteString(e.AccessInterfaceType.String())
	buf.WriteString(`}`)

	return buf.String()
}

type Address struct {
	City string
}

func (a *Address) Serialize() string {
	var buf bytes.Buffer

	buf.WriteString(`{city: "`)
	buf.WriteString(a.City)
	buf.WriteString(`"}`)

	return buf.String()
}

type Price struct {
	PricingModel PricingModel
	Money        Money
}

func (p *Price) Serialize() string {
	var buf bytes.Buffer

	buf.WriteString(`{money: `)
	buf.WriteString(p.Money.Serialize())
	buf.WriteString(`, pricingModel: `)
	buf.WriteString(p.PricingModel.String())
	buf.WriteString(`}`)

	return buf.String()
}

type Money struct {
	Amount   float64 // TODO: look at more precise numeric type here
	Currency Currency
}

func (m *Money) Serialize() string {
	var buf bytes.Buffer

	buf.WriteString(`{amount: `)
	buf.WriteString(strconv.FormatFloat(m.Amount, 'g', -1, 64))
	buf.WriteString(`, currency: `)
	buf.WriteString(m.Currency.String())
	buf.WriteString(`}`)

	return buf.String()
}

// Activation represents an activation of a resource. This comprises a boolean
// flag, and an expiration time. If the flag is set to true and the expiration
// time is in the future, then the offering is active; otherwise it is inactive.
type Activation struct {
	Status         bool      `json:"status"`
	ExpirationTime time.Time `json:"-"`
}

// Serialize converts our Activation into the structure required to send to the
// marketplace
func (a *Activation) Serialize() string {
	var buf bytes.Buffer

	buf.WriteString(`{status: `)
	buf.WriteString(strconv.FormatBool(a.Status))
	buf.WriteString(`, expirationTime: `)
	buf.WriteString(toEpochMs(a.ExpirationTime))
	buf.WriteString(`} `)

	return buf.String()
}

func (a *Activation) UnmarshalJSON(b []byte) error {
	// create anonymous struct for unmarshalling
	d := struct {
		Status         bool  `json:"status"`
		ExpirationTime int64 `json:"expirationTime"`
	}{}

	err := json.Unmarshal(b, &d)
	if err != nil {
		return err
	}

	a.Status = d.Status
	a.ExpirationTime = fromEpochMs(d.ExpirationTime)

	return nil
}
