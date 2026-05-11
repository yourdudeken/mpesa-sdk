# Configuration Reference

All three SDKs accept similar configuration with language-appropriate naming conventions.

## Common Fields

| Field | Required | Description |
|-------|----------|-------------|
| `environment` | Yes | `"sandbox"` or `"production"` |
| `mpesaConsumerKey` / `mpesa_consumer_key` | Yes | API consumer key from Daraja portal |
| `mpesaConsumerSecret` / `mpesa_consumer_secret` | Yes | API consumer secret from Daraja portal |
| `passkey` | Yes | STK Push passkey (Lipa na Mpesa) |
| `shortcode` | Yes | Business shortcode / Till number for PAYBILL |
| `initiatorName` / `initiator_name` | Yes | API initiator username |
| `initiatorPassword` / `initiator_password` | Yes | API initiator password (used for security credential) |

## Optional Fields

| Field | Purpose |
|-------|---------|
| `b2cShortcode` / `b2c_shortcode` | Shortcode for B2C/B2B/B2Pochi transactions |
| `tillNumber` / `till_number` | Till number for TILL transaction type |
| `b2cConsumerKey` / `b2c_consumer_key` | Separate consumer key for B2C/B2B operations |
| `b2cConsumerSecret` / `b2c_consumer_secret` | Separate consumer secret for B2C/B2B operations |
| `callbacks` / `callbacks` | Map of callback URLs (see [Callbacks](./callbacks.md)) |

## Language-specific

### Node.js
```typescript
new Mpesa({
  environment: 'sandbox',
  mpesaConsumerKey: '...',
  mpesaConsumerSecret: '...',
  passkey: '...',
  shortcode: '174379',
  initiatorName: 'testapi',
  initiatorPassword: '...',
  b2cShortcode: '600000',
  tillNumber: '123456',
  b2cConsumerKey: '...',
  b2cConsumerSecret: '...',
  callbacks: {
    callback_url: 'https://example.com/callback',
    // ... see callbacks docs
  },
});
```

### Python
```python
MpesaConfig(
    environment='sandbox',
    mpesa_consumer_key='...',
    mpesa_consumer_secret='...',
    passkey='...',
    shortcode='174379',
    initiator_name='testapi',
    initiator_password='...',
    b2c_shortcode='600000',
    till_number='123456',
    b2c_consumer_key='...',
    b2c_consumer_secret='...',
    callbacks={
        'callback_url': 'https://example.com/callback',
        # ... see callbacks docs
    },
)
```

### Go
```go
&mpesa.Config{
    Environment:         "sandbox",
    MpesaConsumerKey:    "...",
    MpesaConsumerSecret: "...",
    Passkey:             "...",
    Shortcode:           "174379",
    InitiatorName:       "testapi",
    InitiatorPassword:   "...",
    B2cShortcode:        "600000",
    TillNumber:          "123456",
    B2cConsumerKey:      "...",
    B2cConsumerSecret:   "...",
    Callbacks: map[string]string{
        "callback_url": "https://example.com/callback",
        // ... see callbacks docs
    },
}
```

## Environment Variables

Keep credentials out of source code. Example env vars:

```
MPESA_ENV=sandbox
MPESA_CONSUMER_KEY=your_key
MPESA_CONSUMER_SECRET=your_secret
MPESA_PASSKEY=your_passkey
MPESA_SHORTCODE=174379
MPESA_INITIATOR_NAME=testapi
MPESA_INITIATOR_PASSWORD=your_password
MPESA_B2C_SHORTCODE=600000
MPESA_TILL_NUMBER=123456
MPESA_CALLBACK_URL=https://example.com/callback
```
