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

// Error is used for unmarshalling error messages from the marketplace
type Error struct {
	Message string `json:"message"`
}

// ErrorResponse is used to unmarshal the response from the marketplace in the
// event of an error.
type ErrorResponse struct {
	Errors []Error `json:"errors"`
}
