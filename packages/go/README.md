# Mpesa Go SDK

A Go SDK for Mpesa Daraja API.

## Installation

```bash
go get github.com/yourdudeken/mpesa-sdk
```

## Usage

```go
package main

import (
    "github.com/yourdudeken/mpesa-sdk/mpesa"
)

func main() {
    config := &mpesa.Config{
        Environment:       "sandbox",
        MpesaConsumerKey:  "your_key",
        MpesaConsumerSecret: "your_secret",
        Shortcode:         "174379",
        Passkey:           "your_passkey",
        InitiatorName:     "testapi",
        InitiatorPassword: "your_password",
    }

    client := mpesa.NewClient(config)
    stkService := mpesa.NewSTKPushService(client.HTTPClient, client.Auth, config)
    stkService.Push("254712345678", 100, "ORDER123")
}
```

## License

MIT License
