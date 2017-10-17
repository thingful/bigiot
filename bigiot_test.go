package bigiot_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thingful/bigiot"
)

func TestProviderConstructor(t *testing.T) {
	p := bigiot.NewProvider("id", "secret")
	assert.NotNil(t, p)
}
