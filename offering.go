package bigiot

type RdfType string

func (r RdfType) String() string {
	return string(r)
}

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
	Name    string
	RdfType RdfType
}

type Extent struct {
	City string
}

type OfferingDescription struct {
	Name       string
	RdfType    RdfType
	Endpoints  []Endpoint
	InputData  []DataField
	OutputData []DataField
	Extent     Extent
}

type Offering struct {
	ID string
	OfferingDescription
}
