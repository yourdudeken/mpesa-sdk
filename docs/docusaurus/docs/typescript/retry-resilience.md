---
sidebar_position: 10
---

# Retry & Resilience

The SDK implements automatic retry with exponential backoff for transient failures.

## Retry Configuration

### TypeScript

```typescript
const mpesa = new Mpesa({
  consumerKey: '...',
  consumerSecret: '...',
  retryConfig: {
    maxRetries: 3,
    baseDelayMs: 1000,
    maxDelayMs: 30000,
  },
  timeout: 30000,
});
```

### Python

```python
client = Mpesa({
    "consumer_key": "...",
    "consumer_secret": "...",
    "timeout": 30,
    "max_retries": 3,
})
```

### Go

```go
mpesa := client.NewClient(types.MpesaConfig{
    ConsumerKey:    "...",
    ConsumerSecret: "...",
    Timeout:        30 * time.Second,
    RetryConfig: types.RetryConfig{
        MaxRetries:  3,
        BaseDelayMs: 1000,
        MaxDelayMs:  30000,
    },
})
```

## Retryable Status Codes

The SDK retries on these HTTP status codes:

- `408` — Request Timeout
- `429` — Rate Limited
- `500` — Internal Server Error
- `502` — Bad Gateway
- `503` — Service Unavailable
- `504` — Gateway Timeout

## Exponential Backoff

```
Delay = baseDelayMs * 2^attempt + random_jitter
```

| Attempt | Delay Range |
|---------|-------------|
| 0 | ~1000ms |
| 1 | ~2000ms |
| 2 | ~4000ms |
| 3 | ~8000ms (capped at maxDelayMs) |

## Timeout Handling

The SDK will throw a `TimeoutError` when a request exceeds the configured timeout.

```typescript
import { TimeoutError } from '@yourdudeken/mpesa-sdk';

try {
  await mpesa.stkPush.initiate({...});
} catch (error) {
  if (error instanceof TimeoutError) {
    console.log('Request timed out, please retry');
  }
}
```

## Rate Limit Handling

When a 429 response is received, the SDK reads the `Retry-After` header and throws a `RateLimitError` with the retry duration:

```typescript
import { RateLimitError } from '@yourdudeken/mpesa-sdk';

try {
  await mpesa.stkPush.initiate({...});
} catch (error) {
  if (error instanceof RateLimitError) {
    console.log(`Retry after ${error.retryAfter}s`);
  }
}
```
