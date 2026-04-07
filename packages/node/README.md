# Mpesa Node.js SDK

A TypeScript SDK for Mpesa Daraja API.

## Installation

```bash
npm install @yourdudeken/mpesa-sdk
```

## Usage

```typescript
import { Mpesa } from '@yourdudeken/mpesa-sdk';

const mpesa = new Mpesa({
  environment: 'sandbox',
  mpesaConsumerKey: 'your_key',
  mpesaConsumerSecret: 'your_secret',
  shortcode: '174379',
  passkey: 'your_passkey',
  initiatorName: 'testapi',
  initiatorPassword: 'your_password',
  callbacks: {
    callbackUrl: 'https://yourdomain.com/callback'
  }
});

// STK Push
const response = await mpesa.stkpush({
  phonenumber: '254712345678',
  amount: 100,
  accountNumber: 'ORDER123'
});

// B2C
const b2cResponse = await mpesa.b2c({
  phonenumber: '254712345678',
  commandId: 'BusinessPayment',
  amount: 1000,
  remarks: 'Payment'
});
```

## API Reference

- `stkpush()` - Lipa na Mpesa Online
- `stkquery()` - Query STK Push status
- `b2c()` - Business to Customer
- `validated_b2c()` - B2C with ID validation
- `b2b()` - Business to Business
- `c2bregisterURLS()` - Register C2B URLs
- `c2bsimulate()` - Simulate C2B
- `transactionStatus()` - Query transaction
- `accountBalance()` - Check balance
- `reversal()` - Reverse transaction
- `b2pochi()` - Business to Pochi

## License

MIT License
