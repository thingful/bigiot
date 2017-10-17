package bigiot_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thingful/bigiot"
)

func TestProviderConstructor(t *testing.T) {
	p := bigiot.NewProvider("id", "secret")
	assert.Equal(t, "id", p.ID)
	assert.Equal(t, "secret", p.Secret)
	assert.Equal(t, bigiot.DefaultMarketplaceURI, p.MarketplaceURI)
}

func TestProviderConstructorWithMarketplace(t *testing.T) {
	p := bigiot.NewProvider(
		"id",
		"secret",
		bigiot.WithMarketplace("https://market-dev.big-iot.org"),
	)

	assert.Equal(t, "https://market-dev.big-iot.org", p.MarketplaceURI)
}

func TestAuthenticate(t *testing.T) {
	p := bigiot.NewProvider("id", "secret")
	assert.Nil(t, p.Authenticate())
}
