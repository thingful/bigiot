package bigiot_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/thingful/bigiot"
)

func TestActivation(t *testing.T) {
	clock := mockClock{
		t: time.Now(),
	}

	testcases := []struct {
		label    string
		input    bigiot.Activation
		expected string
	}{
		{
			label: "simple expiration",
			input: bigiot.Activation{
				Status:         true,
				ExpirationTime: clock.Now().Add(10 * time.Minute),
			},
			expected: fmt.Sprintf("{status: true, expirationTime: %v} ", bigiot.ToEpochMs(clock.Now().Add(10*time.Minute))),
		},
		{
			label: "with duration",
			input: bigiot.Activation{
				Status:   true,
				Duration: 15 * time.Minute,
			},
			expected: fmt.Sprintf("{status: true, expirationTime: %v} ", bigiot.ToEpochMs(clock.Now().Add(15*time.Minute))),
		},
		{
			label: "with neither",
			input: bigiot.Activation{
				Status: true,
			},
			expected: fmt.Sprintf("{status: true, expirationTime: %v} ", bigiot.ToEpochMs(clock.Now().Add(10*time.Minute))),
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.label, func(t *testing.T) {
			got := testcase.input.Serialize(clock)
			assert.Equal(t, testcase.expected, got)
		})
	}
}
