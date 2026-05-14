# mpesa-sdk — Go

Production-grade Go SDK for Safaricom M-Pesa Daraja API.

## Installation

```bash
go get github.com/yourdudeken/mpesa-sdk/go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "os"

    "github.com/yourdudeken/mpesa-sdk/go/client"
    "github.com/yourdudeken/mpesa-sdk/go/types"
)

func main() {
    mpesa := client.NewClient(types.MpesaConfig{
        ConsumerKey:    os.Getenv("MPESA_CONSUMER_KEY"),
        ConsumerSecret: os.Getenv("MPESA_CONSUMER_SECRET"),
        Environment:    types.Sandbox,
        Passkey:        os.Getenv("MPESA_PASSKEY"),
    })

    resp, err := mpesa.STKPush(context.Background(), types.STKPushRequest{
        BusinessShortCode: 174379,
        TransactionType:   types.CustomerPayBillOnline,
        Amount:            1,
        PartyA:            254722000000,
        PartyB:            174379,
        PhoneNumber:       254722111111,
        CallBackURL:       "https://example.com/callback",
        AccountReference:  "INV-001",
        TransactionDesc:   "Payment",
    })

    if err != nil {
        fmt.Printf("STK Push failed: %v\n", err)
        os.Exit(1)
    }
    fmt.Printf("Checkout ID: %s\n", resp.CheckoutRequestID)
}
```

## API

All methods accept `context.Context` and return typed responses.

| Method | Description |
|--------|-------------|
| `STKPush()` | Initiate STK Push |
| `STKQuery()` | Query STK Push status |
| `C2BRegisterURL()` | Register C2B URLs |
| `C2BSimulate()` | Simulate C2B transaction |
| `B2C()` | Send B2C payment |
| `B2B()` | Send B2B payment |
| `Reversal()` | Reverse a transaction |
| `TransactionStatus()` | Query transaction status |
| `AccountBalance()` | Query account balance |
| `DynamicQR()` | Generate Dynamic QR |

## Packages

- `client/` — HTTP client with auth, retry, and all API methods
- `types/` — Shared request/response types
- `errors/` — Structured error types
- `webhooks/` — Webhook event manager
- `services/` — Higher-level service layer

## Documentation

Full documentation at [https://yourdudeken.github.io/mpesa-sdk](https://yourdudeken.github.io/mpesa-sdk)
