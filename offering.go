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

type RdfType string

func (r RdfType) String() string {
	return string(r)
}

type EndpointType int

const (
	HTTPGet EndpointType = iota
	HTTPPost
	WebSocket
)

type Endpoint struct {
	URI          string
	EndpointType EndpointType
}

type DataField struct {
	Name    string
	RdfType RdfType
}

type Extent struct {
	City string
}

type OfferingDescription struct {
	Name       string
	RdfType    RdfType
	Endpoints  []Endpoint
	InputData  []DataField
	OutputData []DataField
	Extent     Extent
}

type Offering struct {
	ID string
	OfferingDescription
}
