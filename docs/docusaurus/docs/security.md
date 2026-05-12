---
sidebar_position: 6
---

# Security Best Practices

## Credential Management

- **Never hardcode credentials** in your source code
- Use environment variables or a secrets manager
- Rotate consumer keys/secrets regularly

```bash
# .env file
MPESA_CONSUMER_KEY=your_key_here
MPESA_CONSUMER_SECRET=your_secret_here
MPESA_PASSKEY=your_passkey_here
```

## Credential Masking

The SDK automatically masks sensitive data in logs:

```typescript
// Logged: { consumerKey: "abc1****" }
// Not logged: { consumerKey: "abc12345secret" }
```

## Webhook Security

- **Verify webhook signatures** using HMAC-SHA256
- **Use HTTPS** for all callback URLs in production
- **Validate payload structure** before processing

## IP Whitelisting

Add these Safaricom IPs to your firewall whitelist:

```
196.201.214.200
196.201.214.206
196.201.213.114
196.201.214.207
196.201.214.208
196.201.213.44
196.201.212.127
196.201.212.138
196.201.212.129
196.201.212.136
196.201.212.74
196.201.212.69
```

## Production Checklist

- [ ] Use HTTPS for all endpoints
- [ ] Store credentials in secrets manager
- [ ] Enable webhook signature verification
- [ ] Monitor rate limits
- [ ] Set up alerting on error rates
- [ ] Log sanitized request/response data
- [ ] Use environment-specific configurations
- [ ] Implement idempotency for critical operations
