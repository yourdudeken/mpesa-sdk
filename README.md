# Mpesa SDK

Multi-language SDK for Mpesa Daraja API.

[![Node.js](https://img.shields.io/badge/Node.js-16+-339933?style=flat-square&logo=node.js)](https://www.npmjs.com/package/@yourdudeken/mpesa-sdk)
[![Python](https://img.shields.io/badge/Python-3.8+-3776AB?style=flat-square&logo=python)](https://pypi.org/project/yourdudeken-mpesa-sdk/)
[![Go](https://img.shields.io/badge/Go-1.20+-00ADD8?style=flat-square&logo=go)](https://pkg.go.dev/github.com/yourdudeken/mpesa-sdk)

## Supported Languages

| Language | Package | Version |
|----------|---------|---------|
| Node.js/TypeScript | `@yourdudeken/mpesa-sdk` | v1.0.0 |
| Python | `yourdudeken-mpesa-sdk` | v1.0.0 |
| Go | `github.com/yourdudeken/mpesa-sdk` | v1.0.0 |

## Quick Start

### Node.js
```bash
npm install @yourdudeken/mpesa-sdk
```

### Python
```bash
pip install yourdudeken-mpesa-sdk
```

### Go
```bash
go get github.com/yourdudeken/mpesa-sdk
```

## Documentation

- [Node.js SDK](./packages/node/README.md)
- [Python SDK](./packages/python/README.md)
- [Go SDK](./packages/go/README.md)

## API Reference

All SDKs expose a consistent interface across all three languages.

### Node.js

```typescript
import { Mpesa } from '@yourdudeken/mpesa-sdk';

const mpesa = new Mpesa({
  environment: 'sandbox',
  mpesaConsumerKey: '...',
  mpesaConsumerSecret: '...',
  passkey: '...',
  shortcode: '174379',
  initiatorName: 'testapi',
  initiatorPassword: '...',
});

await mpesa.stkpush({ phonenumber, amount, accountNumber });
await mpesa.stkquery(checkoutRequestID);
await mpesa.b2c({ phonenumber, commandId, amount, remarks });
await mpesa.validated_b2c({ phonenumber, commandId, amount, remarks, idNumber });
await mpesa.b2b({ receiverShortcode, commandId, amount, remarks, accountNumber });
await mpesa.c2bregisterURLS({ shortcode, confirmUrl, validateUrl });
await mpesa.c2bsimulate({ phonenumber, amount, shortcode, commandId });
await mpesa.accountBalance({ shortcode, identifierType, remarks });
await mpesa.transactionStatus({ shortcode, transactionId, identifierType, remarks });
await mpesa.reversal({ shortcode, transactionId, amount, remarks });
await mpesa.b2pochi({ phonenumber, amount, remarks });
```

### Python

```python
from mpesa import Mpesa, MpesaConfig

mpesa = Mpesa(MpesaConfig(
    environment='sandbox',
    mpesa_consumer_key='...',
    mpesa_consumer_secret='...',
    passkey='...',
    shortcode='174379',
    initiator_name='testapi',
    initiator_password='...',
))

mpesa.stkpush(phonenumber, amount, account_number)
mpesa.stkquery(checkout_request_id)
mpesa.b2c(phonenumber, command_id, amount, remarks)
mpesa.validated_b2c(phonenumber, command_id, amount, remarks, id_number)
mpesa.b2b(receiver_shortcode, command_id, amount, remarks, account_number)
mpesa.c2b_register_urls(shortcode, confirm_url, validate_url)
mpesa.c2bsimulate(phonenumber, amount, shortcode, command_id)
mpesa.account_balance(shortcode, identifier_type, remarks)
mpesa.transaction_status(shortcode, transaction_id, identifier_type, remarks)
mpesa.reversal(shortcode, transaction_id, amount, remarks)
mpesa.b2pochi(phonenumber, amount, remarks)
```

### Go

```go
import "github.com/yourdudeken/mpesa-sdk/mpesa"

client := mpesa.NewClient(&mpesa.Config{
    Environment:         "sandbox",
    MpesaConsumerKey:    "...",
    MpesaConsumerSecret: "...",
    Passkey:             "...",
    Shortcode:           "174379",
    InitiatorName:       "testapi",
    InitiatorPassword:   "...",
})

client.Stkpush(phonenumber, amount, accountNumber, callbackURL)
client.Stkquery(checkoutRequestID, callbackURL)
client.B2c(phonenumber, commandId, amount, remarks)
client.Validated_b2c(phonenumber, commandId, amount, remarks, idNumber)
client.B2b(receiverShortcode, commandId, amount, remarks, accountNumber)
client.C2bregisterURLS(shortcode, confirmUrl, validateUrl)
client.C2bsimulate(phonenumber, amount, shortcode, commandId, accountNumber)
client.AccountBalance(shortcode, identifierType, remarks)
client.TransactionStatus(shortcode, transactionId, identifierType, remarks)
client.Reversal(shortcode, transactionId, amount, remarks)
client.B2pochi(phonenumber, amount, remarks)
```

## Supported APIs

- **STK Push** - Lipa na Mpesa Online
- **B2C** - Business to Customer
- **B2B** - Business to Business
- **C2B** - Customer to Business
- **B2Pochi** - Business to Pochi
- **Account Balance** - Check till balance
- **Transaction Status** - Query transaction
- **Reversal** - Reverse transaction

## License

MIT License - see [LICENSE.md](./LICENSE.md)
