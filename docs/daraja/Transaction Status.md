# Transaction Status API

Check the status of a transaction.

**Endpoint:** `POST https://sandbox.safaricom.co.ke/mpesa/transactionstatus/v1/query`

## Overview
The Transaction Status API can be used as a secondary reconciliation mechanism when Callbacks are not received. To check the status of a transaction, you need either an M-Pesa Receipt number or an Originator Conversation ID of the transaction.

### How It Works
1. Organization sends API request to check transaction status to Daraja
2. Daraja authenticates and pushes the request to M-PESA system
3. M-PESA checks transaction status and responds to Daraja
4. Daraja pushes the response back to the organization

> **Note:** This API is asynchronous and can be consumed over the internet, VPN, or Multiprotocol Switch.

## Use Cases
Used to check status of a transaction, especially for transactions with delayed callbacks, helping organizations make informed decisions.

## Getting Started

### Prerequisites
- Daraja Account on Safaricom Developer Portal
- Sandbox app with API credentials
- Consumer Key & Consumer Secret

### Good to Know
This API is asynchronous. Can be used for C2B, B2B, B2C, Reversal, and IMT transactions.

## Request Body
```json
{
  "Initiator": "testapiuser",
  "SecurityCredential": "ClONZiMYBpc65lmpJ7nvnrDmUe0WvHvA5QbOsPjEo92B6IGFwDdvdeJIFL0kgwsEKWu6SQKG4ZZUxjC",
  "CommandID": "TransactionStatusQuery",
  "TransactionID": "NEF61H8J60",
  "OriginalConversationID": "7071-4170-a0e5-8345632bad442144258",
  "PartyA": "600782",
  "IdentifierType": "4",
  "ResultURL": "http://myservice:8080/transactionstatus/result",
  "QueueTimeOutURL": "http://myservice:8080/timeout",
  "Remarks": "OK",
  "Occasion": "OK"
}
```

## Request Parameter Definition

| Name | Description | Type | Sample |
|------|-------------|------|--------|
| CommandID | Only 'TransactionStatusQuery' | String | TransactionStatusQuery |
| PartyA | Organization/MSISDN receiving the transaction | Numeric | Short code (6-9 digits) or MSISDN (12 digits) |
| IdentifierType | Type of organization receiving the transaction | Numeric | 4 – Organization shortcode |
| Remarks | Comments sent with the transaction | String | Up to 100 characters |
| Initiator | Name of the initiator for authentication | Alpha-Numeric | Username for API user |
| SecurityCredential | Encrypted password for the initiator | String | Encrypted password |
| QueueTimeoutURL | Path storing timeout transaction info | URL | https://ip:port or domain:port/path |
| TransactionID | Mpesa Receipt Number | Alpha-Numeric | LXXXXXX1234 |
| ResultURL | Path storing transaction info | URL | https://ip:port/path or domain:port/path |
| Occasion | Optional parameter | String | Up to 100 characters |
| OriginalConversationID | Originator Conversation ID of the transaction | String | 7071-4170-a0e5-8345632bad442144258 |

## Response Body
```json
{
  "OriginatorConversationID": "1236-7134259-1",
  "ConversationID": "AG_20210709_1234409f86436c583e3f",
  "ResponseCode": "0",
  "ResponseDescription": "Accept the service request successfully."
}
```

## Response Parameter Definition

| Name | Description | Parameter Type | Sample |
|------|-------------|----------------|--------|
| OriginatorConversationID | Unique request ID for tracking | Alpha-Numeric | 1236-7134259-1 |
| ConversationID | Unique request ID from M-PESA | Alpha-Numeric | AG_20210709_1234409f86436c583e3f |
| ResponseCode | Status code (0 = success) | Number | 0 |
| ResponseDescription | Response description | String | Accept the service request successfully |

## Callback Payload
```json
{
  "Result": {
    "ConversationID": "AG_20180223_0000493344ae97d86f75",
    "OriginatorConversationID": "3213-416199-2",
    "ReferenceData": {
      "ReferenceItem": { "Key": "Occasion" }
    },
    "ResultCode": 0,
    "ResultDesc": "The service request is processed successfully.",
    "ResultParameters": {
      "ResultParameter": [
        { "Key": "DebitPartyName", "Value": "600310 - Safaricom333" },
        { "Key": "DebitPartyName", "Value": "254708374149 - John Doe" },
        { "Key": "OriginatorConversationID", "Value": "3211-416020-3" },
        { "Key": "InitiatedTime", "Value": "20180223054112" },
        { "Key": "DebitAccountType", "Value": "Utility Account" },
        { "Key": "DebitPartyCharges", "Value": "Fee For B2C Payment|KES|22.40" },
        { "Key": "ReasonType", "Value": "Business Payment to Customer via API" },
        { "Key": "TransactionStatus", "Value": "Completed" },
        { "Key": "FinalisedTime", "Value": "20180223054112" },
        { "Key": "Amount", "Value": "300" },
        { "Key": "ConversationID", "Value": "AG_20180223_000041b09c22e613d6c9" },
        { "Key": "ReceiptNo", "Value": "MBN31H462N" }
      ]
    },
    "ResultType": 0,
    "TransactionID": "MBN0000000"
  }
}
```

## Results Parameter Definition

| Name | Description | Type | Sample |
|------|-------------|------|--------|
| ConversationID | Unique identifier from M-PESA | String | AG_20180223_0000493344ae97d86f75 |
| OriginatorConversationID | Unique request identifier | String | 3213-416199-2 |
| ReferenceData | Reference data for transaction log | ReferenceData | n/a |
| ReferenceItem | Reference data item | ParameterType | n/a |
| ResultCode | Processing status (0 = success) | String | 0 |
| ResultDesc | Description of Result Code | String | The service request is processed successfully. |
| ResultParameters | Specific parameters | n/a | n/a |
| Key | Parameter name | String | DebitPartyName |
| Value | Parameter value | String | 600310 - Safaricom333 |
| ResultType | 0: completed, 1: waiting | Integer | 0 |
| TransactionID | Unique transaction identifier | String | MBN0000000 |

### Key Result Parameters

| Key | Description |
|-----|-------------|
| DebitPartyName | Name of debit party (organization or customer) |
| OriginatorConversationID | Original conversation ID |
| InitiatedTime | When transaction was initiated (YYYYMMDDHHmmss) |
| DebitAccountType | Type of account debited (e.g., Utility Account) |
| DebitPartyCharges | Fees charged |
| ReasonType | Reason/type of transaction |
| TransactionStatus | Status of transaction (e.g., Completed, Cancelled, Declined, Expired) |
| FinalisedTime | When transaction was finalized (YYYYMMDDHHmmss) |
| Amount | Transaction amount |
| ReceiptNo | M-PESA receipt number |

## Error Codes

| Error | Possible Cause | Mitigation |
|-------|----------------|------------|
| 500.003.1001 Internal Server Error | Server failure | Check setup and endpoints |
| 400.003.01 Invalid Access Token | Wrong/expired token | Regenerate token |
| 400.003.02 Bad Request | Missing something | Check API documentation |
| 500.003.03 Quota Violation | Multiple requests violating TPS | Send reasonable number of requests |
| 500.003.02 Spike Arrest Violation | Endpoints generating many errors | Ensure endpoint is accessible and responsive |
| 404.003.01 Resource not found | Wrong endpoint | Check M-PESA API endpoint |
| 404.001.04 Invalid Authentication Header | Wrong HTTP method | Use POST for all except Authorization (GET) |
| 400.002.05 Invalid Request Payload | Improper request body | Submit correct payload, avoid typos |

## Response Codes

| Response Code | Response Description |
|---------------|---------------------|
| 0 | Success |

## Result Codes

| Result Code | Result Description |
|-------------|-------------------|
| 0 | Success |
| SFC_IC0003 | The operator does not exist |

## Transaction Statuses
All transactions have three statuses:
1. **Initiated** - Pending revalidation
2. **Authorized/Pending Authorized** - Depends on validation requirements
3. **Final status** (Cancelled, Declined, Completed, or Expired) - After pre-validation/validation and credit party feedback

## Testing

### Option 1: Daraja Simulator
Create a test app, select Transaction Status product, use automated simulator.

### Option 2: Postman
Generate access token and initiate transactions.

**Sandbox Token Endpoint:** `https://sandbox.safaricom.co.ke/oauth/v1/generate?grant_type=client_credentials`
**Production Token Endpoint:** `https://api.safaricom.co.ke/oauth/v1/generate?grant_type=client_credentials`

## Go Live
Attach integration to a live pay bill/till number. Navigate to GO LIVE tab and fill in live data.

## Support
- **Chatbot:** Daraja Chatbot
- **Production Issues:** Incident Management Page or apisupport@safaricom.co.ke

## FAQs
- **How to test Transaction Status API?** Get app credentials, generate token, initiate request with TransactionID or OriginalConversationID and PartyA
- **Can I query B2C and B2B transactions?** Yes, Transaction Status works for C2B, B2B, B2C, IMT, or Reversal transactions
- **Different transaction statuses?** Initiated → Authorized/Pending Authorized → Final (Cancelled, Declined, Completed, Expired)
- **API reversal failing with 'initiator information is invalid'?** Use correct API user (access channel: API), validate username, ensure user is active
- **Process of creating Initiator?** Create Business Manager → Create API operator → Set API user password (valid for 90 days)
- **How to generate Security Credential?** Encrypt base64 initiator password with M-Pesa public key (RSA + PKCS#1.5 padding)
- **What is a short code?** Unique number allocated to pay bill or buy goods organization
- **How to log in to M-PESA portal?** Launch https://org.ke.m-pesa.com, enter shortcode, username, password, verification code, OTP, set new password
- **What to do when Business Administrator role is dormant?** Email M-PESABusiness@Safaricom.co.ke for activation
