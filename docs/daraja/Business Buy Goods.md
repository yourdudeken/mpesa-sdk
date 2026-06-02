# Business Buy Goods API

Receive payments for goods and services through M-Pesa till number.

**Endpoint:** `POST https://sandbox.safaricom.co.ke/mpesa/c2b/v1/simulate`

## Overview
The Business Buy Goods API enables merchants to receive payment requests from customers for retail purchases. Customers pay by entering the till number and transaction is completed instantly.

### Key Features
- Real-time payment processing
- Instant transaction confirmation
- Automatic payment routing
- Customer receipt notifications
- Integrated with M-Pesa USSD menu

## How It Works
1. Customer selects Buy Goods option in M-Pesa
2. Customer enters merchant till number
3. Transaction details displayed to customer
4. Payment is processed and confirmed
5. Funds credited to merchant account
6. Both parties receive transaction confirmation

## Use Cases
- Retail store payments
- Goods and services sales
- Shop purchases
- Marketplace transactions
- POS system integration
- Kiosk operations

## Getting Started

### Prerequisites
- Daraja Account on Safaricom Developer Portal
- Sandbox app with API credentials
- Buy Goods (Till Number) merchant account
- Business Admin/Manager operators setup

### Good to Know
Buy Goods designed for retail transactions. Till numbers typically 5-6 digit codes assigned to physical or online stores.

## Request Body (Simulation)
```json
{
  "ShortCode": "600000",
  "CommandID": "SimulateC2BTrans",
  "Amount": "1000",
  "Msisdn": "254722000000",
  "BillRefNumber": "INV-001"
}
```

## Request Parameter Definition

| Name | Description | Type | Sample |
|------|-------------|------|--------|
| ShortCode | Merchant till number | Numeric | 600000 |
| CommandID | Command type (always SimulateC2BTrans for testing) | String | SimulateC2BTrans |
| Amount | Transaction amount in KES | Numeric | 1000 |
| Msisdn | Customer phone number (format: 2547XXXXXXXX) | Numeric | 254722000000 |
| BillRefNumber | Optional bill reference number | String | INV-001 |

## Response Body
```json
{
  "ResponseDescription": "Accept the service request successfully.",
  "ResponseCode": "0",
  "MerchantRequestID": "29115-34620561-1",
  "CheckoutRequestID": "ws_CO_191220191020363925",
  "CustomerMessage": "success"
}
```

## Response Parameter Definition

| Name | Description | Type | Sample |
|------|-------------|------|--------|
| ResponseDescription | Response message | String | Accept the service request successfully |
| ResponseCode | Success indicator (0=success) | String | 0 |
| MerchantRequestID | Merchant transaction ID | String | 29115-34620561-1 |
| CheckoutRequestID | M-Pesa transaction ID | String | ws_CO_191220191020363925 |
| CustomerMessage | Customer notification message | String | success |

## Callback Payload (Real Transaction)
```json
{
  "Body": {
    "stkCallback": {
      "MerchantRequestID": "29115-34620561-1",
      "CheckoutRequestID": "ws_CO_191220191020363925",
      "ResultCode": 0,
      "ResultDesc": "The service request is processed successfully.",
      "CallbackMetadata": {
        "Item": [
          { "Name": "Amount", "Value": 1000 },
          { "Name": "MpesaReceiptNumber", "Value": "NLJ7RT61SV" },
          { "Name": "TransactionDate", "Value": 20191219102115 },
          { "Name": "PhoneNumber", "Value": 254722000000 }
        ]
      }
    }
  }
}
```

## Account Types Receiving Payments

| Account Type | Description |
|-------------|-------------|
| Till Number | Direct retail point accepting payments |
| Buy Goods Account | Merchant account for goods/services |
| Utility Account | Account receiving customer payments |
| Working Account | Processing account with transaction balance |

## Transaction Flow
1. Customer initiates M-Pesa menu
2. Selects Buy Goods option
3. Enters till number
4. Enters amount
5. Confirms transaction
6. Enters M-Pesa PIN
7. Transaction completes
8. Funds credited to till

## Error Codes

| errorCode | errorMessage | Mitigation |
|-----------|-------------|------------|
| 400.002.02 | Bad Request - Invalid Till | Ensure till number exists |
| 401.002.01 | Invalid Access Token | Regenerate token |
| 404.002.01 | Resource not found | Check till number |
| 500.001.1001 | Till does not exist | Verify till configuration |
| 500.003.02 | System is busy | Retry request |

## Transaction Limits

| Limit | Value |
|-------|-------|
| Minimum Amount | KES 1 |
| Maximum Amount | KES 250,000 |
| Daily Limit | KES 500,000 |
| Account Balance Limit | KES 500,000 |

## Testing
Use Daraja Simulator with predefined till numbers for testing.

## Go Live
Configure live till number and submit Go Live application with business details.

## Support
- **Chatbot:** Daraja Chatbot
- **Email:** apisupport@safaricom.co.ke
