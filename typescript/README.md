# @yourdudeken/mpesa-sdk — TypeScript

Production-grade TypeScript SDK for Safaricom M-Pesa Daraja API.

## Installation

```bash
npm install @yourdudeken/mpesa-sdk axios
```

## Quick Start

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

console.log(response.CheckoutRequestID);
```

## API

All methods return typed responses.

| Service | Methods |
|---------|---------|
| `mpesa.stkPush` | `initiate()`, `query()` |
| `mpesa.c2b` | `registerURL()`, `simulate()` |
| `mpesa.b2c` | `send()` |
| `mpesa.b2b` | `send()` |
| `mpesa.reversal` | `reverse()` |
| `mpesa.transactionStatus` | `query()` |
| `mpesa.accountBalance` | `query()` |
| `mpesa.dynamicQR` | `generate()` |
| `mpesa.webhooks` | `on()`, `off()`, `handleEvent()` |

## Framework Integrations

- **Express**: `createExpressMiddleware()`
- **Fastify**: `createFastifyPlugin()`

## Documentation

Full documentation at [https://yourdudeken.github.io/mpesa-sdk](https://yourdudeken.github.io/mpesa-sdk)
