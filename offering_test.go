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
			expected: fmt.Sprintf("{ status: true, expirationTime: %v }", ToEpochMs(clock.Now().Add(10*time.Minute))),
		},
		{
			label: "with duration",
			input: Activation{
				Status:   true,
				Duration: 15 * time.Minute,
			},
			expected: fmt.Sprintf("{ status: true, expirationTime: %v }", ToEpochMs(clock.Now().Add(15*time.Minute))),
		},
		{
			label: "with neither",
			input: Activation{
				Status: true,
			},
			expected: fmt.Sprintf("{ status: true, expirationTime: %v }", ToEpochMs(clock.Now().Add(10*time.Minute))),
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

func TestSerializeLocation(t *testing.T) {
	clock := mocks.Clock{
		T: time.Now(),
	}

	loc := Location{
		Lng: 23.2,
		Lat: 45.2,
	}

	assert.Equal(t, `{ lng: 23.2, lat: 45.2 }`, loc.serialize(clock))
}

func TestSerializeBoundingBox(t *testing.T) {
	clock := mocks.Clock{
		T: time.Now(),
	}

	bb := BoundingBox{
		Location1: Location{
			Lng: 23.2,
			Lat: 45.2,
		},
		Location2: Location{
			Lng: 24.2,
			Lat: 46.2,
		},
	}

	assert.Equal(t, `{ l1: { lng: 23.2, lat: 45.2 }, l2: { lng: 24.2, lat: 46.2 } }`, bb.serialize(clock))
}

func TestSerializeSpatialExtent(t *testing.T) {
	clock := mocks.Clock{
		T: time.Now(),
	}

	testcases := []struct {
		label    string
		input    SpatialExtent
		expected string
	}{
		{
			label: "only city",
			input: SpatialExtent{
				City: "Edinburgh",
			},
			expected: `{ city: "Edinburgh" }`,
		},
		{
			label: "with bounding box",
			input: SpatialExtent{
				City: "Edinburgh",
				BoundingBox: &BoundingBox{
					Location1: Location{
						Lng: 23.2,
						Lat: 45.2,
					},
					Location2: Location{
						Lng: 24.2,
						Lat: 46.2,
					},
				},
			},
			expected: `{ city: "Edinburgh", boundary: { l1: { lng: 23.2, lat: 45.2 }, l2: { lng: 24.2, lat: 46.2 } } }`,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.label, func(t *testing.T) {
			assert.Equal(t, testcase.expected, testcase.input.serialize(clock))
		})
	}
}

func TestSerializeOfferingDescription(t *testing.T) {
	now := time.Unix(0, 0)
	duration := 10 * time.Minute
	clock := mocks.Clock{T: now}

	testcases := []struct {
		label    string
		input    *OfferingDescription
		expected string
	}{
		{
			label: "duration no bounding",
			input: &OfferingDescription{
				LocalID: "TestOffering",
				Name:    "Test Offering",
				RdfURI:  "urn:proposed:RandomValues",
				Outputs: []DataField{
					{
						Name:   "value",
						RdfURI: "schema:random",
					},
				},
				Endpoints: []Endpoint{
					{
						URI:                 "https://example.com/random",
						EndpointType:        HTTPGet,
						AccessInterfaceType: BIGIoTLib,
					},
				},
				License: OpenDataLicense,
				Price: Price{
					Money: Money{
						Amount:   0.001,
						Currency: EUR,
					},
					PricingModel: PerAccess,
				},
				SpatialExtent: &SpatialExtent{
					City: "Berlin",
				},
				Activation: &Activation{
					Status:   true,
					Duration: duration,
				},
			},
			expected: `mutation addOffering { addOffering ( input: { id: "", localId: "TestOffering", name: "Test Offering", activation: { status: true, expirationTime: 600000 }, rdfUri: "urn:proposed:RandomValues", outputs: [{ name: "value", rdfUri: "schema:random" }], endpoints: [{ uri: "https://example.com/random", endpointType: HTTP_GET, accessInterfaceType: BIGIOT_LIB }], license: OPEN_DATA_LICENSE, price: { money: { amount: 0.001, currency: EUR }, pricingModel: PER_ACCESS }, spatialExtent: { city: "Berlin" } } ) { id name activation { status expirationTime } } }`,
		},
		{
			label: "duration no bounding",
			input: &OfferingDescription{
				LocalID: "TestOffering",
				Name:    "Test Offering",
				RdfURI:  "urn:proposed:RandomValues",
				Inputs: []DataField{
					{
						Name:   "value",
						RdfURI: "schema:random",
					},
				},
				Outputs: []DataField{
					{
						Name:   "value",
						RdfURI: "schema:random",
					},
				},
				Endpoints: []Endpoint{
					{
						URI:                 "https://example.com/random",
						EndpointType:        HTTPGet,
						AccessInterfaceType: BIGIoTLib,
					},
				},
				License: OpenDataLicense,
				Price: Price{
					Money: Money{
						Amount:   0.001,
						Currency: EUR,
					},
					PricingModel: PerAccess,
				},
				SpatialExtent: &SpatialExtent{
					City: "Berlin",
					BoundingBox: &BoundingBox{
						Location1: Location{
							Lng: 0,
							Lat: 0,
						},
						Location2: Location{
							Lng: 1,
							Lat: 1,
						},
					},
				},
				Activation: &Activation{
					Status:   true,
					Duration: duration,
				},
			},
			expected: `mutation addOffering { addOffering ( input: { id: "", localId: "TestOffering", name: "Test Offering", activation: { status: true, expirationTime: 600000 }, rdfUri: "urn:proposed:RandomValues", inputs: [{ name: "value", rdfUri: "schema:random" }], outputs: [{ name: "value", rdfUri: "schema:random" }], endpoints: [{ uri: "https://example.com/random", endpointType: HTTP_GET, accessInterfaceType: BIGIOT_LIB }], license: OPEN_DATA_LICENSE, price: { money: { amount: 0.001, currency: EUR }, pricingModel: PER_ACCESS }, spatialExtent: { city: "Berlin", boundary: { l1: { lng: 0, lat: 0 }, l2: { lng: 1, lat: 1 } } } } ) { id name activation { status expirationTime } } }`,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.label, func(t *testing.T) {
			assert.Equal(t, testcase.expected, testcase.input.serialize(clock))
		})
	}
}
