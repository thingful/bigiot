# bigiot

Go implementation of the BIG IoT library/SDK

## Planned Features

* Register an offering in the marketplace
* Unregister an offering from the marketplace
* Validating tokens presented by offering subscribers
* Discovering an offering in the marketplace
* Subscribing to an offering

## Create Offering

```go
// create provider and authenticate with marketplace
provider, err := bigiot.NewProvider(providerID, providerSecret)
if err != nil {
    panic(err)
}

err := provider.Authenticate()
if err != nil {
    panic(err)
}

// create the description of the offering
offeringInput := &bigiot.OfferingInput{
    LocalID: "ParkingOffering",
    Name: "Demo Parking Offering",
    Endpoints: []bigiot.Endpoint{
        {
            URI:                 "http://example.com/parking",
            EndpointType:        bigiot.HTTPGet,
            AccessInterfaceType: bigiot.External,
        },
    },
    InputData: []bigiot.DataField{
        {
            Name: "longitude",
            RdfURI: "schema:longitude",
        },
        {
            Name: "latitude",
            RdfURI: "schema:latitude",
        },
        {
            Name: "radius",
            RdfURI: "schema:geoRadius",
        },
    },
    OutputData: []bigiot.Data{
        {
            Name: "geoCoordinates",
            RdfURI: "schema:geoCoordinates",
        }
    },
    Extent: bigiot.Address{
        City: "Barcelona",
    },
    Activation: bigiot.Activation{
        Status: true,
        ExpirationTime: expirationTime,
    }
}

// register the offering - note this registration is timeboxed so should be
// refreshed regularly
offering, err = provider.RegisterOffering(ctx, offeringInput)
if err != nil {
    panic(err)
}

// deregister the offering
err = offering.Deregister()
if err != nil {
    panic(err)
}
```

## Create provider client with custom marketplace

```go
provider := bigiot.NewProvider(
    "id",
    "secret",
    bigiot.WithMarketplace("https://market-dev.big-iot.org"),
)
```

## List offerings

```go
provider := bigiot.NewProvider("id", "secret")
err := provider.Authenticate()
if err != nil {
    panic(err)
}

offerings, err := provider.Offerings()
if err != nil {
    panic(err)
}

for _, offering := range offerings {
    fmt.Println(offering.Name)
}
```