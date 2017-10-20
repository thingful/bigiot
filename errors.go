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

const (
	// ErrUnexpectedResponse is an error returned when we receive an unexpected
	// response from the BIG IoT API.
	ErrUnexpectedResponse = Error("Unexpected HTTP response code")
)

// Error is a type alias for string, allowing us to export const error values
type Error string

// Error is the implementation of the error interface, allowing us to use our
// Error type as an error.
func (e Error) Error() string {
	return string(e)
}
