# Mpesa SDK Documentation

## Overview

This is a multi-language SDK for the Mpesa Daraja API. All SDKs provide a consistent interface across languages.

## Installation

See individual package documentation in `packages/{language}/README.md`

## Configuration

### PHP (Laravel)
```php
// .env
MPESA_ENVIRONMENT=sandbox
MPESA_CONSUMER_KEY=your_consumer_key
MPESA_CONSUMER_SECRET=your_consumer_secret
MPESA_BUSINESS_SHORTCODE=174379
SAFARICOM_PASSKEY=your_passkey
MPESA_INITIATOR_NAME=testapi
MPESA_INITIATOR_PASSWORD=your_password
```

### Node.js
```typescript
const config = {
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
};
```

### Python
```python
config = MpesaConfig(
    environment='sandbox',
    mpesa_consumer_key='your_key',
    mpesa_consumer_secret='your_secret',
    passkey='your_passkey',
    shortcode='174379',
    initiator_name='testapi',
    initiator_password='your_password'
)
```

## API Methods

All SDKs support:
- `stkpush()` - Lipa na Mpesa Online payment
- `b2c()` - Business to Customer payment
- `b2b()` - Business to Business payment
- `c2bregisterURLS()` - Register C2B URLs
- `c2bsimulate()` - Simulate C2B transaction
- `transactionStatus()` - Query transaction status
- `accountBalance()` - Check account balance
- `reversal()` - Reverse a transaction
- `b2pochi()` - Business to Pochi payment

## Error Handling

Each SDK has its own exception types. See individual package documentation.
