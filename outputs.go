package bigiot

type Activation struct {
	Status         Boolean
	ExpirationTime Boolean
}

type RdfType struct {
	URI  String
	Name String
}

type RdfContext struct {
	Context  String
	Prefixes []Prefix
}

type Prefix struct {
	Prefix String
	URI    String
}

type Endpoint struct {
	EndpointType        EndpointType
	URI                 String
	AccessInterfaceType AccessInterfaceType
}

type DataField struct {
	Name    String
	RdfType RdfType
}

type Address struct {
	City String
}

type Price struct {
	PricingModel PricingModel
	Money        MoneyInput
}

type Offering struct {
	ID         String
	Name       String
	Activation Activation
	RdfType    RdfType
	RdfContext RdfContext
	Endpoints  []Endpoint
	OutputData []DataField
	InputData  []DataField
	Extent     Address
	License    License
	Price      Price
}
