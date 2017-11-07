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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEpochMs(t *testing.T) {
	testcases := []struct {
		name         string
		input        time.Time
		expected     string
		expectedBack time.Time
	}{
		{"zero time", time.Unix(0, 0), "0", time.Unix(0, 0)},
		{"nano time", time.Unix(0, 1509983101577890997), "1509983101577", time.Unix(0, 1509983101577000000)},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			got := toEpochMs(testcase.input)
			assert.Equal(t, testcase.expected, got)
			val, err := strconv.ParseInt(got, 10, 64)
			assert.Nil(t, err)
			back := fromEpochMs(val)
			assert.Equal(t, testcase.expectedBack.UTC(), back)
		})
	}
}
