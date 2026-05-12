---
sidebar_position: 7
---

# Production Deployment

## Environment Setup

### Prerequisites

1. **M-PESA Account** with PayBill/Till/B2C capabilities
2. **Daraja Portal App** with live credentials
3. **Business Administrator** and **API Operator** users
4. **M-PESA Public Key Certificate** for security credentials
5. **SSL Certificate** for your servers

### Go Live Process

1. Create a Daraja app and select the required products
2. Configure your shortcode (PayBill/Till/B2C)
3. Register callback URLs (one-time in production)
4. Set up API operators with appropriate roles
5. Test with small amounts first

## Production Configuration

### TypeScript

```typescript
const mpesa = new Mpesa({
  consumerKey: process.env.MPESA_CONSUMER_KEY!,
  consumerSecret: process.env.MPESA_CONSUMER_SECRET!,
  environment: 'production',
  passkey: process.env.MPESA_PASSKEY!,
  timeout: 30000,
  retryConfig: {
    maxRetries: 3,
    baseDelayMs: 1000,
    maxDelayMs: 30000,
  },
});
```

## Error Monitoring

Set up logging hooks to monitor API errors:

```typescript
const mpesa = new Mpesa({
  ...config,
  logging: {
    onError: (errorLog) => {
      // Send to your monitoring system
      console.error(errorLog);
    },
  },
});
```

## Rate Limiting

M-Pesa APIs have rate limits. Implement:

- Exponential backoff (built into SDK)
- Queue management for bulk operations
- Monitoring for 429 responses
