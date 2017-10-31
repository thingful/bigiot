package bigiot

// EndpointType represents the type of an Endpoint accessible via the BIGIoT
// Marketplace
type EndpointType string

const (
	// HTTPGet is a const value representing an endpoint type accessible via HTTP
	// GET
	HTTPGet EndpointType = "HTTP_GET"

	// HTTPPost is a const value representing an endpoint type accessible via HTTP
	// POST
	HTTPPost EndpointType = "HTTP_POST"

	// WebSocket is a const value representing an endpoint type accessible via
	// WebSockets
	WebSocket EndpointType = "WEBSOCKET"
)

// String is our implementation of the Stringer interface for EndpointType
// instances.
func (e EndpointType) String() string {
	return string(e)
}

type AccessInterfaceType string

const (
	BIGIoTLib AccessInterfaceType = "BIGIOT_LIB"

	External AccessInterfaceType = "EXTERNAL"
)

func (a AccessInterfaceType) String() string {
	return string(a)
}

type License string

const (
	CreativeCommons          License = "CREATIVE_COMMONS"
	OpenDataLicense          License = "OPEN_DATA_LICENSE"
	NonCommercialDataLicense License = "NON_COMMERCIAL_DATA_LICENSE"
)

func (l License) String() string {
	return string(l)
}

type PricingModel string

const (
	Free      PricingModel = "FREE"
	PerMonth  PricingModel = "PER_MONTH"
	PerAccess PricingModel = "PER_ACCESS"
	PerByte   PricingModel = "PER_BYTE"
)

func (p PricingModel) String() string {
	return string(p)
}

type Currency string

const (
	EUR Currency = "EUR"
)

func (c Currency) String() string {
	return string(c)
}
