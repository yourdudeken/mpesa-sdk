# Callback URLs

Most Mpesa API operations are asynchronous — Safaricom sends the result to your callback URL when the transaction completes.

## How Callbacks Work

1. You configure callback URLs in the SDK config.
2. When calling an API method, you can pass `callbackUrl` to override the config.
3. Safaricom's servers call your URL with the transaction result.

## Callback Keys (all SDKs)

All three SDKs use the same snake_case keys internally:

| Key | Used By | Purpose |
|-----|---------|---------|
| `callback_url` | STK Push | Result of STK push transaction |
| `b2c_result_url` | B2C, Validated B2C | B2C transaction result |
| `b2c_timeout_url` | B2C, Validated B2C | B2C timeout fallback |
| `b2b_result_url` | B2B | B2B transaction result |
| `b2b_timeout_url` | B2B | B2B timeout fallback |
| `b2pochi_result_url` | B2 Pochi | B2Pochi transaction result |
| `b2pochi_timeout_url` | B2 Pochi | B2Pochi timeout fallback |
| `c2b_confirmation_url` | C2B Register | C2B payment confirmation |
| `c2b_validation_url` | C2B Register | C2B payment validation |
| `balance_result_url` | Account Balance | Balance query result |
| `balance_timeout_url` | Account Balance | Balance query timeout |
| `status_result_url` | Transaction Status | Status query result |
| `status_timeout_url` | Transaction Status | Status query timeout |
| `reversal_result_url` | Reversal | Reversal result |
| `reversal_timeout_url` | Reversal | Reversal timeout |

## Resolution Order

When the SDK builds the API request, it resolves each callback URL in this order:

1. **Parameter passed directly** to the method (highest priority)
2. **Config callbacks map** (if configured during client construction)
3. **Error** — the API call is rejected if no URL is found

## Node.js Example

```typescript
const mpesa = new Mpesa({
  // ...
  callbacks: {
    callbackUrl: 'https://api.example.com/mpesa/callback',
    b2cResultUrl: 'https://api.example.com/mpesa/b2c/result',
    b2cTimeoutUrl: 'https://api.example.com/mpesa/b2c/timeout',
    statusResultUrl: 'https://api.example.com/mpesa/status/result',
    // ... (use either camelCase or snake_case keys)
  },
});

// Override per-call:
await mpesa.stkpush({
  phonenumber: '254712345678',
  amount: 10,
  accountNumber: 'INV-001',
  callbackUrl: 'https://custom.url/callback',  // overrides config
});
```

## Python Example

```python
config = MpesaConfig(
    # ...
    callbacks={
        'callback_url': 'https://api.example.com/mpesa/callback',
        'b2c_result_url': 'https://api.example.com/mpesa/b2c/result',
        'b2c_timeout_url': 'https://api.example.com/mpesa/b2c/timeout',
    },
)

# Override per-call:
mpesa.stkpush('254712345678', 10, 'INV-001', callback_url='https://custom.url/callback')
```

## Go Example

```go
config := &mpesa.Config{
    // ...
    Callbacks: map[string]string{
        "callback_url":   "https://api.example.com/mpesa/callback",
        "b2c_result_url": "https://api.example.com/mpesa/b2c/result",
        "b2c_timeout_url": "https://api.example.com/mpesa/b2c/timeout",
    },
}

// Override per-call:
client.Stkpush("254712345678", 10, "INV-001", "https://custom.url/callback")
```
