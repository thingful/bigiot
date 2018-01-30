package bigiot

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/thingful/bigiot/mocks"
)

func TestActivation(t *testing.T) {
	clock := mocks.Clock{
		T: time.Now(),
	}

	testcases := []struct {
		label    string
		input    Activation
		expected string
	}{
		{
			label: "simple expiration",
			input: Activation{
				Status:         true,
				ExpirationTime: clock.Now().Add(10 * time.Minute),
			},
			expected: fmt.Sprintf("{status: true, expirationTime: %v} ", ToEpochMs(clock.Now().Add(10*time.Minute))),
		},
		{
			label: "with duration",
			input: Activation{
				Status:   true,
				Duration: 15 * time.Minute,
			},
			expected: fmt.Sprintf("{status: true, expirationTime: %v} ", ToEpochMs(clock.Now().Add(15*time.Minute))),
		},
		{
			label: "with neither",
			input: Activation{
				Status: true,
			},
			expected: fmt.Sprintf("{status: true, expirationTime: %v} ", ToEpochMs(clock.Now().Add(10*time.Minute))),
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.label, func(t *testing.T) {
			got := testcase.input.serialize(clock)
			assert.Equal(t, testcase.expected, got)
		})
	}
}

func TestActivateOffering(t *testing.T) {
	clock := mocks.Clock{
		T: time.Now(),
	}

	testcases := []struct {
		label    string
		input    ActivateOffering
		expected string
	}{
		{
			label: "simple expiration",
			input: ActivateOffering{
				ID:             "Organisation-Provider-Offering",
				ExpirationTime: clock.Now().Add(10 * time.Minute),
			},
			expected: fmt.Sprintf(`mutation activateOffering { activateOffering ( input: { id: "Organisation-Provider-Offering", expirationTime: %v } ) { id activation { status expirationTime } } }`, ToEpochMs(clock.Now().Add(10*time.Minute))),
		},
		{
			label: "with duration",
			input: ActivateOffering{
				ID:       "Organisation-Provider-Offering",
				Duration: 15 * time.Minute,
			},
			expected: fmt.Sprintf(`mutation activateOffering { activateOffering ( input: { id: "Organisation-Provider-Offering", expirationTime: %v } ) { id activation { status expirationTime } } }`, ToEpochMs(clock.Now().Add(15*time.Minute))),
		},
		{
			label: "with neither",
			input: ActivateOffering{
				ID: "Organisation-Provider-Offering",
			},
			expected: fmt.Sprintf(`mutation activateOffering { activateOffering ( input: { id: "Organisation-Provider-Offering", expirationTime: %v } ) { id activation { status expirationTime } } }`, ToEpochMs(clock.Now().Add(10*time.Minute))),
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.label, func(t *testing.T) {
			got := testcase.input.serialize(clock)
			assert.Equal(t, testcase.expected, got)
		})
	}
}
