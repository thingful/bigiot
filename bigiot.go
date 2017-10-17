package bigiot

// Provider is our struct for containing auth credentials to interact with the
// BIGIoT marketplace.
type Provider struct {
	ID             string
	Secret         string
	MarketplaceURI string
}

// NewProvider instantiates and returns a configured Provider instance. The
// required parameters to the function are the provider ID and secret. If you
// want to connect to a marketplace other than the offical marketplace (i.e.
// connecting to a local instance for testing), you can configure this by means
// of the variadic third parameter, which can be used for additional
// configuration.
func NewProvider(id, secret string, options ...func(*Provider)) *Provider {
	provider := &Provider{
		ID:     id,
		Secret: secret,
	}

	for _, opt := range options {
		opt(provider)
	}

	return provider
}

// Authenticate attempts to connect to the marketplace and authenticate the
// client. Returns any error to the caller.
func (p *Provider) Authenticate() error {
	return nil
}

// WithMarketplace is a functional configuration option allowing us to
// optionally set a custom marketplace URI when constructing a Provider
// instance.
func WithMarketplace(marketplaceURI string) func(*Provider) {
	return func(p *Provider) {
		p.MarketplaceURI = marketplaceURI
	}
}
