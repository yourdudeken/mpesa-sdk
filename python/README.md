# M-Pesa SDK — Python

Production-grade Python SDK for Safaricom M-Pesa Daraja API.

## Installation

```bash
pip install yourdudeken-mpesa-sdk
```

Requires Python 3.11+.

## Quick Start

```python
from mpesa import Mpesa

client = Mpesa({
    "consumer_key": "your_consumer_key",
    "consumer_secret": "your_consumer_secret",
    "environment": "sandbox",
    "passkey": "your_passkey",
})

response = client.stk_push({
    "BusinessShortCode": 174379,
    "TransactionType": "CustomerPayBillOnline",
    "Amount": 1,
    "PartyA": 254722000000,
    "PartyB": 174379,
    "PhoneNumber": 254722111111,
    "CallBackURL": "https://example.com/callback",
    "AccountReference": "INV-001",
    "TransactionDesc": "Payment",
})

print(f"Checkout ID: {response.CheckoutRequestID}")
```

## Features

- OAuth authentication with automatic token management
- STK Push (M-Pesa Express)
- STK Query
- C2B (Register URL & Simulate)
- B2C (Business to Customer)
- B2B (Business to Business)
- Transaction Reversal
- Transaction Status Query
- Account Balance Query
- Dynamic QR Generation
- Webhook handling with event-driven architecture
- Structured error hierarchy
- Retry with exponential backoff
- Pydantic v2 models for request/response validation

## Documentation

Full documentation: [https://yourdudeken.github.io/mpesa-sdk](https://yourdudeken.github.io/mpesa-sdk)
