package bigiot

type RdfType string

type EndpointType int

const (
	HTTPGet EndpointType = iota
	HTTPPost
	WebSocket
)

type Endpoint struct {
	URI          string
	EndpointType EndpointType
}

type DataField struct {
	Name   string
	RdfURI string
}

type OfferingDescription struct {
	ID         string
	Name       string
	RdfType    RdfType
	Endpoints  []Endpoint
	InputData  []DataField
	OutputData []DataField
}
