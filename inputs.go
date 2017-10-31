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

type Input interface{}

type ActivationInput struct {
	Status         Boolean
	ExpirationTime Long
}

type DataFieldInput struct {
	Name   String
	RdfURI String
}

type EndpointInput struct {
	EndpointType        EndpointType
	URI                 String
	AccessInterfaceType AccessInterfaceType
}

type AddressInput struct {
	City String
}

type PriceInput struct {
	PricingModel PricingModel
	Money        MoneyInput
}

type MoneyInput struct {
	Amount   BigDecimal
	Currency Currency
}

type AddOffering struct {
	ID         String
	LocalID    String
	Name       String
	RdfURI     String
	InputData  []DataFieldInput
	OutputData []DataFieldInput
	Endpoints  []EndpointInput
	Extent     AddressInput
	License    License
	Price      PriceInput
	Activation ActivationInput
}
