package bigiot_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thingful/bigiot"
)

func TestEndpointTypeStringer(t *testing.T) {
	assert.Equal(t, "HTTP_GET", bigiot.HTTPGet.String())
}
