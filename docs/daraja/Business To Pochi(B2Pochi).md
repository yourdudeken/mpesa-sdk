# Business To Pochi (B2Pochi) API

Enables businesses to send money to M-Pesa Pochi La Biashara merchant accounts.

**Endpoint:** `POST https://sandbox.safaricom.co.ke/mpesa/b2pochi/v1/paymentrequest`

## Overview
The Business To Pochi (B2Pochi) API allows registered businesses to send payments directly to M-Pesa Pochi La Biashara merchant accounts. This facilitates merchant-to-merchant payments and business disbursements.

### Key Features
- Direct merchant-to-merchant transfers
- Instant payment processing
- Secure transaction handling
- Real-time confirmation notifications
- Automated reconciliation support

## How It Works
1. Business initiates Pochi payment request
2. API validates merchant credentials
3. M-Pesa processes the payment
4. Pochi merchant receives funds instantly
5. Both parties receive transaction confirmation
6. Payment reflected in Pochi balance

## Use Cases
- Merchant commission payments
- Vendor disbursements
- Agent payouts
- Multi-level distributor payments
- Franchise fee transfers
- Business partner settlements

## Getting Started

### Prerequisites
- Daraja Account on Safaricom Developer Portal
- Sandbox app with API credentials
- Business merchant account
- Pochi merchant account details
- Business Admin/Manager operators setup

### Good to Know
Pochi La Biashara is M-Pesa merchant account for business transactions and deposits.

## Request Body
```json
{
  "InitiatorName": "testuser",
  "SecurityCredential": "SAFVNChNHfVtXEZMBuVo+a1Hwr+DtrUVN3zVg==",
  "CommandID": "BusinessToBusinessTransfer",
  "Amount": "5000",
  "SenderIdentifier": "4",
  "ReceiverIdentifier": "2",
  "PartyA": "600000",
  "PartyB": "600100",
  "AccountReference": "pochi_payment",
  "Remarks": "Pochi transfer",
  "QueueTimeOutURL": "http://myservice:8080/queuetimeouturl",
  "ResultURL": "http://myservice:8080/result"
}
```

## Request Parameter Definition

| Name | Description | Type | Sample |
|------|-------------|------|--------|
| InitiatorName | Initiator username | Alpha-Numeric | testapi772 |
| SecurityCredential | Encrypted password | String | EToK4lNR... |
| CommandID | Command type | String | BusinessToBusinessTransfer |
| Amount | Transfer amount in KES | Numeric | 5000 |
| SenderIdentifier | Sender type (4: Shortcode) | Numeric | 4 |
| ReceiverIdentifier | Receiver type (2: Till) | Numeric | 2 |
| PartyA | Sender shortcode | Numeric | 600000 |
| PartyB | Pochi merchant till/account | Numeric | 600100 |
| AccountReference | Reference for tracking | String | pochi_payment |
| Remarks | Transaction comments | String | Pochi transfer |
| QueueTimeoutURL | Timeout notification URL | URL | https://ip:port/path |
| ResultURL | Result notification URL | URL | https://ip:port/path |

## Response Body
```json
{
  "OriginatorConversationID": "515-5258779-3",
  "ConversationID": "AG_20200123_0000417fed8ed666e976",
  "ResponseCode": "0",
  "ResponseDescription": "Accept the service request successfully"
}
```

## Response Parameter Definition

| Name | Description | Type | Sample |
|------|-------------|------|--------|
| OriginatorConversationID | Unique request identifier | String | 515-5258779-3 |
| ConversationID | M-Pesa unique identifier | String | AG_20200123_0000417fed8ed666e976 |
| ResponseCode | Success indicator | String | 0 |
| ResponseDescription | Status message | String | Accept the service request successfully |

## Pochi Merchant Requirements
- Registered M-Pesa Pochi account
- Valid merchant till number
- Verified business details
- Active merchant status

## Error Codes

| errorCode | errorMessage | Mitigation |
|-----------|-------------|------------|
| 400.002.02 | Bad Request - Invalid Data | Check parameters |
| 401.002.01 | Invalid Access Token | Regenerate token |
| 404.002.01 | Pochi account not found | Verify account details |
| 500.001.1001 | Invalid merchant | Confirm Pochi setup |
| 500.003.02 | System is busy | Retry request |

## Testing
Use Daraja Simulator with test Pochi merchant accounts.

## Go Live
Submit Live Pochi account details and merchant information.

## Support
- **Chatbot:** Daraja Chatbot
- **Email:** apisupport@safaricom.co.ke
