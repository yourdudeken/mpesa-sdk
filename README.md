# M-Pesa SDK

Production-grade SDK ecosystem for Safaricom M-Pesa Daraja APIs.

[![CI](https://github.com/yourdudeken/mpesa-sdk/actions/workflows/ci.yml/badge.svg)](https://github.com/yourdudeken/mpesa-sdk/actions/workflows/ci.yml)
[![npm version](https://img.shields.io/npm/v/@yourdudeken/mpesa-sdk)](https://www.npmjs.com/package/@yourdudeken/mpesa-sdk)
[![PyPI version](https://img.shields.io/pypi/v/yourdudeken-mpesa-sdk)](https://pypi.org/project/yourdudeken-mpesa-sdk/)
[![Go Reference](https://pkg.go.dev/badge/github.com/yourdudeken/mpesa-sdk.svg)](https://pkg.go.dev/github.com/yourdudeken/mpesa-sdk)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Available SDKs

| Language | Package | Version |
|----------|---------|---------|
| **TypeScript** | `@yourdudeken/mpesa-sdk` | ![npm](https://img.shields.io/npm/v/@yourdudeken/mpesa-sdk) |
| **Python** | `yourdudeken-mpesa-sdk` | ![pypi](https://img.shields.io/pypi/v/yourdudeken-mpesa-sdk) |
| **Go** | `github.com/yourdudeken/mpesa-sdk` | ![go](https://img.shields.io/github/v/tag/yourdudeken/mpesa-sdk?filter=go/v*.*.*) |

## Features

- **OAuth Authentication** with automatic token management
- **STK Push (M-Pesa Express)** with password generation
- **STK Query** - Check transaction status
- **C2B** - Register URLs & simulate transactions
- **B2C** - Business to Customer payments
- **B2B** - Business to Business payments
- **Transaction Reversal** - Reverse C2B transactions
- **Transaction Status Query** - Reconciliation
- **Account Balance Query** - Check balances
- **Dynamic QR** - Generate QR codes
- **Webhook Handling** - Event-driven callbacks
- **Structured Errors** - Typed error hierarchy
- **Retry & Resilience** - Exponential backoff
- **Framework Integrations** - Express, Fastify, FastAPI, Flask, Gin

## Quick Start

### TypeScript

```bash
npm install @yourdudeken/mpesa-sdk axios
```

```typescript
import { Mpesa } from '@yourdudeken/mpesa-sdk';

const mpesa = new Mpesa({
  consumerKey: process.env.MPESA_CONSUMER_KEY!,
  consumerSecret: process.env.MPESA_CONSUMER_SECRET!,
  environment: 'sandbox',
  passkey: process.env.MPESA_PASSKEY!,
});

const response = await mpesa.stkPush.initiate({
  BusinessShortCode: 174379,
  TransactionType: 'CustomerPayBillOnline',
  Amount: 1,
  PartyA: 254722000000,
  PartyB: 174379,
  PhoneNumber: 254722111111,
  CallBackURL: 'https://example.com/callback',
  AccountReference: 'INV-001',
  TransactionDesc: 'Payment',
});
```

### Python

```bash
pip install yourdudeken-mpesa-sdk
```

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

```bash
go get github.com/yourdudeken/mpesa-sdk
```

```go
import (
    "github.com/yourdudeken/mpesa-sdk/client"
    "github.com/yourdudeken/mpesa-sdk/types"
)

mpesa := client.NewClient(types.MpesaConfig{
    ConsumerKey:    os.Getenv("MPESA_CONSUMER_KEY"),
    ConsumerSecret: os.Getenv("MPESA_CONSUMER_SECRET"),
    Environment:    types.Sandbox,
    Passkey:        os.Getenv("MPESA_PASSKEY"),
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

```
mpesa-sdk/
├── openapi/          # OpenAPI specification (single source of truth)
├── typescript/       # TypeScript SDK
├── python/           # Python SDK
├── go/               # Go SDK
├── examples/         # Usage examples per language
├── docs/             # Documentation (Docusaurus)
├── scripts/          # Automation scripts
├── shared/           # Shared utilities
└── .github/          # CI/CD workflows
```

## Documentation

Full documentation: [https://yourdudeken.github.io/mpesa-sdk](https://yourdudeken.github.io/mpesa-sdk)

## License

MIT
