# Reversals API

Reverses an M-Pesa transaction.

**Endpoint:** `POST https://sandbox.safaricom.co.ke/mpesa/reversal/v1/request`

## Overview
The Reversals API enables the reversal of Customer-to-Business (C2B) transactions.

### How It Works
1. Organization sends a reversal request via the Daraja API Gateway
2. Gateway authenticates and forwards to M-PESA system
3. M-PESA processes the reversal, refunds the customer, sends SMS notification
4. Result returned through Daraja to the organization

> **Note:** This API is asynchronous and can be consumed over the internet, VPN, or Multiprotocol Switch.

## Use Cases
- Reverse erroneous payments made to M-PESA Collection Account
- Reverse double payments
- Reverse payments where services were not fulfilled

## Getting Started

### Prerequisites
- Daraja Account on Safaricom Developer Portal
- Sandbox app with API credentials
- Consumer Key & Consumer Secret
- Test data available in simulator section
- For production: live pay bill/till number with Business Admin/Manager operators

### Good to Know
This API is asynchronous. Can also check status of C2B, B2B, B2C, Reversal, and IMT transactions.

## Request Body
```json
{
  "Initiator": "apiop37",
  "SecurityCredential": "jUb+dOXJiBDui8FnruaFckZJQup3kmmCH5XJ4NY/Oo3KaUTmJbxUiVgzBjqdL533u5Q435MT2VJwr/ /1fuZvA===",
  "CommandID": "TransactionReversal",
  "TransactionID": "PDU91HIVIT",
  "Amount": "200",
  "ReceiverParty": "603021",
  "RecieverIdentifierType": "11",
  "ResultURL": "https://mydomain.com/reversal/result",
  "QueueTimeOutURL": "https://mydomain.com/reversal/queue",
  "Remarks": "Payment reversal"
}
```

## Request Parameter Definition

| Parameter | Description | Type | Required | Sample |
|-----------|-------------|------|----------|--------|
| Initiator | Username of API user on M-PESA portal | String | Yes | johndoe |
| SecurityCredential | Encrypted password for API user | String | Yes | RC6E9WDx9X2c6z3gp0oC5Th== |
| CommandID | Only 'TransactionReversal' allowed | String | Yes | TransactionReversal |
| Amount | Transaction amount | Numeric | Yes | 100 |
| ReceiverParty | Organization Short Code | Numeric | Yes | 600997 |
| RecieverIdentifierType | Type of Organization (should be '11') | Numeric | Yes | 11 |
| Remarks | Additional information (2-100 chars) | String | Yes | Any string |
| QueueTimeOutURL | URL for timeout notification | URL | Yes | https://mydomain.com/reversal/timedout |
| ResultURL | URL for result notification | URL | Yes | https://mydomain.com/reversal/result |
| TransactionID | M-PESA Receipt Number for transaction being reversed | String | Yes | PDU91HIVIT |

## Response Body

### Success Response
```json
{
  "OriginatorConversationID": "f1e2-4b95-a71d-b30d3cdbb7a7735297",
  "ConversationID": "AG_20210706_20106e9209f64bebd05b",
  "ResponseCode": "0",
  "ResponseDescription": "Accept the service request successfully."
}
```

## Response Parameter Definition

| Parameter | Description | Sample Value | Type |
|-----------|-------------|--------------|------|
| ConversationID | Unique global identifier from M-PESA | 4f31-fd2d5deb744c | String |
| OriginatorConversationID | Unique identifier from API proxy | AG_20210706_2010ead4245 | String |
| ResponseCode | Status code (0 = success) | 0 | Numeric |
| ResponseDescription | Acknowledgment message | Accept the service request successfully | String |

## Callback Result Payload

### Successful Callback
```json
{
  "Result": {
    "ResultType": 0,
    "ResultCode": 0,
    "ResultDesc": "The service request is processed successfully.",
    "OriginatorConversationID": "dad6-4c34-8787-c8cb963a496d1268232",
    "ConversationID": "AG_20211114_201018edbbf9f1582eaa",
    "TransactionID": "SKE52PAWR9",
    "ResultParameters": {
      "ResultParameter": [
        { "Key": "DebitAccountBalance", "Value": "Utility Account|KES|7722179.62|7722179.62|0.00|0.00" },
        { "Key": "Amount", "Value": 1.0 },
        { "Key": "TransCompletedTime", "Value": 20211114132711 },
        { "Key": "OriginalTransactionID", "Value": "SKC82PACB8" },
        { "Key": "Charge", "Value": 0.0 },
        { "Key": "CreditPartyPublicName", "Value": "254705912645 - NICHOLAS JOHN SONGOK" },
        { "Key": "DebitPartyPublicName", "Value": "600992 - Safaricom Daraja 992" }
      ]
    },
    "ReferenceData": {
      "ReferenceItem": { "Key": "QueueTimeoutURL", "Value": "https://internalsandbox.safaricom.co.ke/mpesa/reversalresults/v1/submit" }
    }
  }
}
```

### Unsuccessful Callback
```json
{
  "Result": {
    "ResultType": 0,
    "ResultCode": "R000002",
    "ResultDesc": "The OriginalTransactionID is invalid.",
    "OriginatorConversationID": "3124-481d-b706-10bdd6fbc8e21792398",
    "ConversationID": "AG_20211114_2010573069aefb6b625a",
    "TransactionID": "SKE0000000",
    "ReferenceData": {
      "ReferenceItem": { "Key": "QueueTimeoutURL", "Value": "https://internalsandbox.safaricom.co.ke/mpesa/reversalresults/v1/submit" }
    }
  }
}
```

## Response Parameter Definition (Callback)

| Parameter | Description | Type | Optional | Sample |
|-----------|-------------|------|----------|--------|
| Result | Root parameter for result message | JSON Object | No | {} |
| ResultType | Status code (usually 0) | Numeric | No | 0 |
| ResultCode | Status code (0 = success) | String | No | 0 |
| ResultDesc | Status message | String | No | The service request is processed successfully. |
| OriginatorConversationID | Unique identifier for reversal request | String | No | 53e3-4aa8-9fe0-... |
| ConversationID | Unique identifier from M-PESA | String | No | AG_20210707_20106f7a33 |
| TransactionID | M-PESA Receipt Number for reversal | String | No | SKE52PAWR9 |
| ResultParameters | Additional transaction details | JSON Object | Yes | {} |
| DebitAccountBalance | Account balances (Account Type\|Currency\|...) | String | Yes | Utility Account\|KES\|7722179.62 |
| Amount | Transaction amount | Decimal | Yes | 1.00 |
| TransCompletedTime | Completion time (YYYYMMDDhhmmss) | String | Yes | 20211114132711 |
| OriginalTransactionID | TransactionID being reversed | String | Yes | SKC82PACB8 |
| Charge | Total fee amount | String | Yes | 0.00 |
| CreditPartyPublicName | Credit Party public name | String | Yes | 254705912645 - NICHOLAS JOHN SONGOK |
| DebitPartyPublicName | Debit Party public name | String | Yes | 600992 - Safaricom Daraja 992 |
| ReferenceData | Additional request details | JSON Object | Yes | {} |

## Result Codes

| ResultCode | ResultDesc | Explanation |
|-----------|-----------|-------------|
| 0 | The service request is processed successfully | Request processed successfully on M-PESA |
| R000002 | The OriginalTransactionID is invalid | TransactionID invalid or doesn't exist |
| R000001 | The transaction has already been reversed | TransactionID already reversed |
| 11 | The DebitParty is in an invalid state | Organization/short code account not active |
| 21 | The initiator is not allowed to initiate | API user lacks Org Reversals Initiator role |
| 2001 | The initiator information is invalid | Invalid API user credentials |
| 2006 | Declined due to account rule | Organization/short code not active |
| 2028 | Not permitted according to product assignment | Short code has no permission |
| 8006 | The security credential is locked | API user password locked |
| 1 | The balance is insufficient | Short code insufficient funds |

## Error Response

### Example
```json
{
  "requestId": "94fc-460e-a970-797968bf6a851272619",
  "errorCode": "400.002.02",
  "errorMessage": "Bad Request - Invalid TransactionID"
}
```

| Parameter | Description | Sample | Type |
|-----------|-------------|--------|------|
| requestId | Unique identifier by API | 30764-19833054-1 | String |
| errorCode | Unique error code | 400.002.02 | String |
| errorMessage | Descriptive failure message | Bad Request - Invalid TransactionID | String |

## Error Response Codes

| errorCode | errorMessage | Mitigation | HTTP |
|-----------|-------------|------------|------|
| 404.001.03 | Invalid Access Token | Regenerate token before expiry | 404 |
| 400.002.02 | Bad Request – Invalid XXXX | Ensure payload is correct | 400 |
| 404.001.01 | Resource not found | Check API endpoint | 404 |
| 500.001.1001 | Internal Server Error | Check setup and documentation | 500 |
| 500.003.02 | Spike Arrest Violation | Avoid multiple requests violating TPS | 500 |
| 500.003.03 | Quota Violation | Avoid multiple requests violating limits | 500 |

## Testing

### Option 1: Daraja Simulator
Create a test app, select Reversal product, use automated simulator.

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
- **Not receiving callbacks?** Ensure endpoint is publicly accessible
- **Invalid access token?** Token expires hourly; regenerate
- **Determine if reversal successful?** Check ResultCode in callback (0 = success)
- **Required API role?** Org Reversals Initiator
- **Activate pending API user?** Set password via M-PESA portal
- **Initiator information is invalid?** Check username, password encryption, algorithm
- **Security credential locked?** Unlock via Business Administrator
- **Invalid API call - no apiproduct match?** Ensure Reversal product is enabled
- **Reverse B2C transaction?** No, done manually on M-PESA portal
- **Bad Request - Invalid RecieverIdentifierType?** Check Content-Type, parameter name, value
- **Process of creating an Initiator?** Create Business Manager → API operator → Set password
- **Generate Security Credential?** Encrypt base64 initiator password with M-Pesa public key (RSA + PKCS#1.5)
