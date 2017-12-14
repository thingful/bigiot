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
	"strconv"
	"time"
)

// toEpochMs takes a time.Time and returns this time as a epoch milliseconds
// formatted as a string.
func toEpochMs(t time.Time) string {
	return strconv.FormatInt(t.UnixNano()*int64(time.Nanosecond)/int64(time.Millisecond), 10)
}

// fromEpochMs takes as input an int value (as returned from toEpochMs) and then
// returns this as a time.Time in UTC.
func fromEpochMs(v int64) time.Time {
	nanosec := v * 1e6
	return time.Unix(0, nanosec).UTC()
}

// Clock is an interface used to make it possible to test time related code more
// easily.
type Clock interface {
	Now() time.Time
}

// Clock is an interface used to make it possible to test time related code more
// easily.
type Clock interface {
	Now() time.Time
}

// realClock is our implementation of the Clock interface that returns the real
// time.
type realClock struct{}

// Now is our implementation of the Clock function interface that returns the
// current time
func (r realClock) Now() time.Time { return time.Now() }
