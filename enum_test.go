package bigiot_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thingful/bigiot"
)

func TestEndpointTypeStringer(t *testing.T) {
	assert.Equal(t, "HTTP_GET", bigiot.HTTPGet.String())
	assert.Equal(t, "HTTP_POST", bigiot.HTTPPost.String())
	assert.Equal(t, "WEBSOCKET", bigiot.WebSocket.String())
}

func TestAccessInterfaceTypeStringer(t *testing.T) {
	assert.Equal(t, "BIGIOT_LIB", bigiot.BIGIoTLib.String())
	assert.Equal(t, "EXTERNAL", bigiot.External.String())
}
