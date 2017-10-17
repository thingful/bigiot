# bigiot

Go implementation of the BIG IoT library/SDK

## Planned Features

* Register an offering in the marketplace
* Unregister an offering from the marketplace
* Validating tokens presented by offering subscribers
* Discovering an offering in the marketplace
* Subscribing to an offering

## Example Provider

```go
// create provider and authenticate with marketplace
provider := bigiot.NewProvider(providerID, providerSecret)
err := provider.Authenticate()
if err != nil {
    panic(err)
}

// create the description of the offering
offeringDescription := &bigiot.OfferingDescription{
    Name: "Demo Parking Offering",
    RDFType: RDFType("bigiot:Parking"),
    Endpoints: []bigiot.Endpoint{
        {
            URI: "http://example.com/parking",
        }
    },
    InputData: []bigiot.InputData{
        {
            Name: "longitude",
            RDFType: RDFType("schema:longitude"),
            ValueType: bigiot.Number,
        },
        {
            Name: "latitude",
            RDFType: RDFType("schema:latitude"),
            ValueType: bigiot.Number,
        }, 
        {
            Name: "radius",
            RDFType: RDFType("schema:geoRadius"),
            ValueType: bigiot.Number,
        },
    },
    OutputData: []bigiot.Data{
        {
            Name: "geoCoordinates",
            RDFType: RDFType("schema:geoCoordinates"),
            ValueType: bigiot.Number,
        }
    },
    Region: &bigiot.Region{
        City: "Barcelona",
    }
}

// register the offering - note this registration is timeboxed so should be
// refreshed regularly
offering, err = provider.RegisterOffering(offeringDescription)
if err != nil {
    panic(err)
}

// deregister the offering
err = offering.Deregister()
if err != nil {
    panic(err)
}
```

## Example Provider with custom marketplace

```go
provider := bigiot.NewProvider(
    "id", 
    "secret",
    bigiot.WithMarketplace("https://market-dev.big-iot.org"),
)
```