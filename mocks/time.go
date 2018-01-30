package mocks

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

import "time"

// Clock is an implementation of the Clock interface for use in tests.
// Returns a canned time for "now".
type Clock struct {
	T time.Time
}

// Now is our implementation of the Clock Now() function that in the real case
// returns the current time, but here we just return the canned time value set
// on the struct.
func (c Clock) Now() time.Time {
	return c.T
}
