# bigiot

Go implementation of the BIG IoT library/SDK

## Implemented Features

* Register an offering in the marketplace
* Unregister an offering from the marketplace

## Planned Features

* Validating tokens presented by offering subscribers
* Discovering an offering in the marketplace
* Subscribing to an offering

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