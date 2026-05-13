---
sidebar_position: 1
---

# Introduction

**M-Pesa SDK** is a production-grade SDK ecosystem for integrating with Safaricom's M-Pesa Daraja API. Available in TypeScript, Python, and Go.

## Features

- **OAuth Authentication** with automatic token management and refresh
- **STK Push (M-Pesa Express)** with password generation
- **STK Query** to check transaction status
- **C2B APIs** for Customer-to-Business payments
- **B2C APIs** for Business-to-Customer disbursements
- **B2B APIs** for Business-to-Business transfers
- **Transaction Reversal** with full callback parsing
- **Transaction Status Query** for reconciliation
- **Account Balance Query** with structured balance parsing
- **Dynamic QR Code** generation
- **Webhook Handling** with event-driven architecture
- **Structured Errors** with typed error hierarchy
- **Retry & Resilience** with exponential backoff
- **Environment Switching** (sandbox/production)
- **Security** with credential masking and input validation

## Supported Languages

| Language | Package | Status |
|----------|---------|--------|
| TypeScript | `@yourdudeken/mpesa-sdk` | Beta |
| Python | `yourdudeken-mpesa-sdk` | Beta |
| Go | `github.com/yourdudeken/mpesa-sdk` | Beta |

## Quick Comparison

### TypeScript
```typescript
import { Mpesa } from '@yourdudeken/mpesa-sdk';

const mpesa = new Mpesa({
  consumerKey: process.env.MPESA_CONSUMER_KEY,
  consumerSecret: process.env.MPESA_CONSUMER_SECRET,
  environment: 'sandbox',
  passkey: process.env.MPESA_PASSKEY,
});

const response = await mpesa.stkPush.initiate({
  BusinessShortCode: 174379,
  TransactionType: 'CustomerPayBillOnline',
  Amount: 100,
  PartyA: 254722000000,
  PartyB: 174379,
  PhoneNumber: 254722111111,
  CallBackURL: 'https://example.com/callback',
  AccountReference: 'INV-001',
  TransactionDesc: 'Payment',
});
```

### Python
```python
from mpesa import Mpesa

client = Mpesa({
    "consumer_key": "...",
    "consumer_secret": "...",
    "environment": "sandbox",
    "passkey": "...",
})

response = client.stk_push({
    "BusinessShortCode": 174379,
    "TransactionType": "CustomerPayBillOnline",
    "Amount": 100,
    "PartyA": 254722000000,
    "PartyB": 174379,
    "PhoneNumber": 254722111111,
    "CallBackURL": "https://example.com/callback",
    "AccountReference": "INV-001",
    "TransactionDesc": "Payment",
})
```

### Go
```go
import (
    "github.com/yourdudeken/mpesa-sdk/client"
    "github.com/yourdudeken/mpesa-sdk/types"
)

mpesa := client.NewClient(types.MpesaConfig{
    ConsumerKey:    "...",
    ConsumerSecret: "...",
    Environment:    types.Sandbox,
    Passkey:        "...",
})

resp, err := mpesa.STKPush(ctx, types.STKPushRequest{
    BusinessShortCode: 174379,
    TransactionType:   types.CustomerPayBillOnline,
    Amount:            100,
    PartyA:            254722000000,
    PartyB:            174379,
    PhoneNumber:       254722111111,
    CallBackURL:       "https://example.com/callback",
    AccountReference:  "INV-001",
    TransactionDesc:   "Payment",
})
```

## Architecture

The SDKs follow a consistent architecture across all languages:

- **Client** - Core HTTP client with auth, retry, and logging
- **Services** - Domain-specific service classes per API endpoint
- **Types/Models** - Strongly typed request/response structures
- **Errors** - Hierarchical error classes with rich context
- **Webhooks** - Event-driven callback handling
- **Utils** - Shared utilities (password gen, validation, etc.)

## Next Steps

- [Installation Guide](./installation)
- [Authentication Guide](./authentication)
- [TypeScript SDK](./typescript/)
- [Python SDK](./python/)
- [Go SDK](./go/)
