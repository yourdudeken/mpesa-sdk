# Business Pay Bill API

Receive payments through M-Pesa Pay Bill merchant account.

**Endpoint:** `POST https://sandbox.safaricom.co.ke/mpesa/c2b/v1/simulate`

## Overview
The Business Pay Bill API enables businesses to receive regular payments from customers through M-Pesa Pay Bill accounts. Customers use this for utility bills, subscriptions, loans, and other recurring payments.

### Key Features
- Recurring payment collection
- Automatic payment routing
- Instant transaction processing
- Bill reference tracking
- Customer receipt notifications

## How It Works
1. Customer selects Pay Bill in M-Pesa menu
2. Customer enters business short code
3. Customer enters account reference/bill number
4. Customer enters payment amount
5. Transaction is processed and confirmed
6. Funds routed to business Pay Bill account
7. Both parties receive confirmation

## Use Cases
- Utility bill payments (water, electricity)
- Insurance premium payments
- Loan repayments
- Subscription fees
- Rent collection
- School fees
- Membership fees

## Getting Started

### Prerequisites
- Daraja Account on Safaricom Developer Portal
- Sandbox app with API credentials
- Pay Bill merchant account
- Business Admin/Manager operators setup
- Configured bill reference system

### Good to Know
Pay Bill designed for recurring payments. Customers can specify account reference to link payment to their account.

## Request Body (Simulation)
```json
{
  "ShortCode": "600000",
  "CommandID": "SimulateC2BTrans",
  "Amount": "5000",
  "Msisdn": "254722000000",
  "BillRefNumber": "ACC-12345"
}
```

## Request Parameter Definition

| Name | Description | Type | Sample |
|------|-------------|------|--------|
| ShortCode | Business short code (Pay Bill number) | Numeric | 600000 |
| CommandID | Command type (SimulateC2BTrans for testing) | String | SimulateC2BTrans |
| Amount | Payment amount in KES | Numeric | 5000 |
| Msisdn | Customer phone number (format: 2547XXXXXXXX) | Numeric | 254722000000 |
| BillRefNumber | Account reference for payment tracking | String | ACC-12345 |

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
      "MerchantRequestID": "f1e2-4b95-a71d-b30d3cdbb7a7942864",
      "CheckoutRequestID": "ws_CO_21072024125243250722943992",
      "ResultCode": 0,
      "ResultDesc": "The service request is processed successfully.",
      "CallbackMetadata": {
        "Item": [
          { "Name": "Amount", "Value": 5000 },
          { "Name": "MpesaReceiptNumber", "Value": "NLJ7RT61SV" },
          { "Name": "TransactionDate", "Value": 20240721125243 },
          { "Name": "AccountReference", "Value": "ACC-12345" },
          { "Name": "PhoneNumber", "Value": 254722000000 }
        ]
      }
    }
  }
}
```

## Payment Flow
1. Customer launches M-Pesa menu
2. Selects Pay Bill option
3. Enters business short code
4. Enters account reference (if required)
5. Enters payment amount
6. Confirms transaction details
7. Enters M-Pesa PIN
8. Payment processed
9. Funds credited to Pay Bill account

## Account Types

| Account Type | Description |
|-------------|-------------|
| Pay Bill Account | Primary account receiving payments |
| Utility Account | Sub-account for payment routing |
| Working Account | Processing account with balance |
| Settlement Account | Account for charge deduction |

## Bill Reference Format
Account references can be:
- Customer account numbers
- Invoice numbers
- Reference IDs (max 12 alphanumeric characters)

## Error Codes

| errorCode | errorMessage | Mitigation |
|-----------|-------------|------------|
| 400.002.02 | Bad Request - Invalid Code | Verify short code |
| 401.002.01 | Invalid Access Token | Regenerate token |
| 404.002.01 | Resource not found | Check short code |
| 500.001.1001 | Short code not found | Confirm account setup |
| 500.003.02 | System busy | Retry request |

## Transaction Limits

| Limit | Value |
|-------|-------|
| Minimum Amount | KES 1 |
| Maximum Amount | KES 250,000 |
| Daily Limit | KES 500,000 |
| Max Balance | KES 500,000 |

## Reconciliation
Received payments reconciled using:
- Transaction Status Query API
- Account Balance API
- M-Pesa Portal statements
- API callback notifications

## Testing
Use Daraja Simulator to test with predefined short codes.

## Go Live
Configure live Pay Bill account and submit Go Live application.

## Support
- **Chatbot:** Daraja Chatbot
- **Email:** apisupport@safaricom.co.ke
