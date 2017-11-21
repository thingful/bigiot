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

// query is a type used when composing GraphQL queries. We use it when
// marshalling our graphql queries before sending to the marketplace.
type query struct {
	Query string `json:"query"`
}

// serializable is an interface for an instance that can serialize itself into
// some form that the BIG IoT Marketplace will accept as input for either query
// or mutatation.
type serializable interface {
	Serialize() string
}
