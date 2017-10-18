package bigiot

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	// DefaultMarketplaceURL is the URI to the default BIG IoT marketplace
	DefaultMarketplaceURL = "https://market.big-iot.org"

	// DefaultTimeout is the default timeout in seconds to set on on requests to
	// the marketplace
	DefaultTimeout = 10
)

// Provider is our struct for containing auth credentials to interact with the
// BIGIoT marketplace.
type Provider struct {
	ID          string
	Secret      string
	AccessToken string
	BaseURL     *url.URL
	httpClient  *http.Client
}

// NewProvider instantiates and returns a configured Provider instance. The
// required parameters to the function are the provider ID and secret. If you
// want to connect to a marketplace other than the offical marketplace (i.e.
// connecting to a local instance for testing), you can configure this by means
// of the variadic third parameter, which can be used for additional
// configuration.
func NewProvider(id, secret string, options ...func(*Provider) error) (*Provider, error) {
	// this is a known good URL, so we can ignore the error here
	u, _ := url.Parse(DefaultMarketplaceURL)

	httpClient := &http.Client{
		Timeout: time.Second * DefaultTimeout,
	}

	provider := &Provider{
		ID:         id,
		Secret:     secret,
		BaseURL:    u,
		httpClient: httpClient,
	}

	var err error

	// apply all functional options
	for _, opt := range options {
		err = opt(provider)
		if err != nil {
			return nil, err
		}
	}

	// now wrap our transport to add authentication header if accessToken is available
	transport := http.DefaultTransport
	if provider.httpClient.Transport != nil {
		transport = provider.httpClient.Transport
	}

	provider.httpClient.Transport = &authTransport{
		proxied:     transport,
		accessToken: "",
	}

	//provider.httpClient.Transport.(roundTripper).accessToken = "foo"
	return provider, nil
}

// Authenticate attempts to connect to the marketplace and authenticate the
// client. Returns any error to the caller.
func (p *Provider) Authenticate() (err error) {
	// deference to make sure we clone our baseURL property rather than modifying
	// the pointed to value
	authURL := *p.BaseURL
	authURL.Path = "/accessToken"

	req, err := http.NewRequest(http.MethodGet, authURL.String(), nil)
	if err != nil {
		return err
	}

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		if cerr := resp.Body.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return ErrUnexpectedResponse
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	p.AccessToken = string(body)

	return nil
}

// WithMarketplace is a functional configuration option allowing us to
// optionally set a custom marketplace URI when constructing a Provider
// instance.
func WithMarketplace(marketplaceURI string) func(*Provider) error {
	return func(p *Provider) error {
		u, err := url.Parse(marketplaceURI)
		if err != nil {
			return err
		}

		p.BaseURL = u

		return nil
	}
}
