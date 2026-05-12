# Customer To Business (C2B) API

Register URL for Validation/Confirmation and Simulate transaction.

**Endpoint:** `POST https://sandbox.safaricom.co.ke/mpesa/c2b/v2/registerurl`

## Overview
The Customer to Business (C2B) API, also known as the Register URL API, enables merchants to receive notifications for successful payments to their Paybill or Till numbers. Funds originate from the customer wallet and are transferred to the merchant's short code.

Payments can be initiated via: SIM Toolkit, Mpesa App, Safaricom App, USSD, NI Push API, or Dynamic QR Code API.

The C2B API allows you to register callback URLs for payment notifications:
- **Validation URL:** Used when a merchant needs to validate payment details before accepting
- **Confirmation URL:** Receives payment notification after successful completion

> **Note:** C2B Transaction Validation is optional and must be activated by emailing apisupport@safaricom.co.ke or M-pesabusiness@safaricom.co.ke.

## How It Works
1. Customer initiates payment to a Paybill or Till number
2. M-PESA validates the request internally
3. If External Validation is enabled:
   - Sends Validation request to registered Validation URL
   - Merchant validates and responds (within ~8 seconds)
   - M-PESA processes based on response
4. If External Validation is disabled, M-PESA completes the transaction
5. If URLs are not reachable, M-PESA uses the default action value
6. If no URLs registered, M-PESA completes the request
7. SMS notifications sent to both customer and merchant

### URL Requirements
- Use publicly available IP addresses or domain names
- Production URLs must be HTTPS; Sandbox allows HTTP
- Avoid keywords: M-PESA, Safaricom, exe, exec, cmd, SQL, query
- Do not use public URL testers (ngrok, mockbin, requestbin) in production
- Sandbox: URLs can be registered multiple times
- Production: One-time registration; contact apisupport@safaricom.co.ke to change

> **Note:** The words "Cancelled/Completed" in ResponseType must be in sentence case and well-spelled.

## Getting Started

### Prerequisites
- Daraja Account on Safaricom Developer Portal
- Sandbox app to get API credentials
- Consumer Key & Consumer Secret
- Test data from simulator section
- Live M-PESA Paybill/Till number for production

### Good to Know
- The API is asynchronous
- Can be consumed over internet, VPN, or Multiprotocol Switch

## Request Body

### Register URLs
```json
{
  "ShortCode": "600984",
  "ResponseType": "Either Cancelled or Completed",
  "ConfirmationURL": "your confirmation URL",
  "ValidationURL": "your validation URL"
}
```

### Simulate Transactions
```json
{
  "ShortCode": 600984,
  "CommandID": "Either CustomerBuyGoodsOnline or CustomerPayBillOnline",
  "Amount": 1,
  "Msisdn": 254708374149,
  "BillRefNumber": "Account reference for Customer paybills and null for customer buy goods"
}
```

## Request Parameter Definition

| Name | Description | Type | Sample |
|------|-------------|------|--------|
| ValidationURL | URL receiving validation request (only if external validation enabled) | URL | https://ip/domain:port/path |
| ConfirmationURL | URL receiving confirmation request upon payment completion | URL | https://ip/domain:port/path |
| ResponseType | Default action if validation URL unreachable (Completed or Cancelled) | String | Completed, Cancelled |
| ShortCode | Unique M-PESA pay bill/till number | Numeric | 600996 |
| CommandID | Transaction type: CustomerBuyGoodsOnline or CustomerPayBillOnline | String | CustomerBuyGoodsOnline / CustomerPayBillOnline |
| Amount | Amount to be transacted | Numeric | 10 |
| Msisdn | Phone number for debit | Numeric | 254708374149 |
| BillRefNumber | Account reference for paybill; null for till number | String | "Test Ref" |

## Response Body

### Simulate Response
```json
{
  "OriginatorCoversationID": "53e3-4aa8-9fe0-8fb5e4092cdd3405976",
  "ResponseCode": "0",
  "ResponseDescription": "Accept the service request successfully."
}
```

### Register URLs Response
```json
{
  "OriginatorCoversationID": "6e86-45dd-91ac-fd5d4178ab523408729",
  "ResponseCode": "0",
  "ResponseDescription": "Success"
}
```

## Response Parameter Definition

| Name | Description | Type | Sample |
|------|-------------|------|--------|
| OriginatorCoversationID | Global unique identifier from API proxy | Alphanumeric | Alpha-numeric string |
| ResponseCode | Indicates whether request is accepted | Alphanumeric | 0 |
| ResponseDescription | Status of the request | String | Success |

> **NB:** You must generate an access token before making register URL API calls.

## Callback Payload

### Validation Request (only if External Validation enabled)
```json
{
  "TransactionType": "Pay Bill",
  "TransID": "RKL51ZDR4F",
  "TransTime": "20231121121325",
  "TransAmount": "5.00",
  "BusinessShortCode": "600966",
  "BillRefNumber": "Sample Transaction",
  "InvoiceNumber": "",
  "OrgAccountBalance": "25.00",
  "ThirdPartyTransID": "",
  "MSISDN": "2547 ***** 126",
  "FirstName": "NICHOLAS",
  "MiddleName": "",
  "LastName": ""
}
```

## Callback Parameter Definition

| Name | Description | Type | Sample |
|------|-------------|------|--------|
| TransactionType | Type specified during payment request | String | Buy Goods or Pay Bill |
| TransID | Unique M-Pesa transaction ID | Alpha-numeric | LHG31AA5TX |
| TransTime | Timestamp (YYYYMMDDHHmmss) | Time | 20170813154301 |
| TransAmount | Amount paid by customer | Numeric | 100 |
| BusinessShortCode | Organization's shortcode (5-6 digits) | String | 654321 |
| BillRefNumber | Account number for the payment | String | Alpha-numeric up to 20 chars |
| OrgAccountBalance | Utility account balance after payment | Decimal | 30671 |
| ThirdPartyTransID | Partner's transaction ID | String | 1234567890 |
| MSISDN | Masked customer number | String | 2547 * 126 |
| FirstName | Customer's first name | String | John |
| MiddleName | Customer's middle name | String | null |
| LastName | Customer's last name | String | null |

## Validation Response

### Accept Transaction
```json
{
  "ResultCode": "0",
  "ResultDesc": "Accepted"
}
```

### Reject Transaction
```json
{
  "ResultCode": "C2B00011",
  "ResultDesc": "Rejected"
}
```

## Result Codes

| ResultCode | ResultDesc |
|-----------|-----------|
| C2B00011 | Invalid MSISDN |
| C2B00012 | Invalid Account Number |
| C2B00013 | Invalid Amount |
| C2B00014 | Invalid KYC Details |
| C2B00015 | Invalid Short code |
| C2B00016 | Other Error |

## Error Codes

| HTTP | Error Message | Possible Cause | Mitigation |
|------|---------------|----------------|------------|
| 500 | 500.003.1001 Internal Server Error | Server failure | Check correct setup and endpoints |
| 500 | 500.003.1001 URLs already registered | Existing URL registered | Request deletion and re-register |
| 400 | 400.003.01 Invalid Access Token | Wrong/expired token | Regenerate token |
| 400 | 400.003.02 Bad Request | Missing something | Check API documentation |
| 500 | 500.003.03 Quota Violation | Multiple requests violating TPS | Send reasonable number of requests |
| 500 | 500.003.02 Spike Arrest Violation | Endpoints generating many errors | Ensure endpoint is accessible |
| 404 | 404.003.01 Resource not found | Wrong endpoint | Check API endpoint |
| 404 | 404.001.04 Invalid Authenticator Header | Wrong HTTP method | Use POST for all except Authorization (GET) |
| 400 | 400.002.05 Invalid Request Payload | Improper request body | Submit correct payload |
| 500 | Duplicate notification info | Existing URLs on aggregator platform | Request deletion from aggregator |

## Testing

### Option 1: Daraja Simulator
Create a test app, select C2B product. Register URLs before each simulation.

### Option 2: Postman
Generate access token and initiate transactions.

**Sandbox Token Endpoint:** `https://sandbox.safaricom.co.ke/oauth/v1/generate?grant_type=client_credentials`
**Production Token Endpoint:** `https://api.safaricom.co.ke/oauth/v1/generate?grant_type=client_credentials`

> **Note:** "Simulate C2B Request" is only available in Sandbox.

## Go Live
Attach integration to a live Paybill/Till number. Fill in live data (short code, organization name, M-PESA admin/manager username).

## M-PESA Organization Portal

### Access
URL: https://org.ke.m-pesa.com/orglogin.action

### Account Types
- **MMF/Working/M-PESA Account:** For business withdrawals
- **Utility Account:** Receives customer payments
- **Charges Paid Account:** Debited for transaction charges
- **Organization Settlement Account:** Settles charges and moves balance automatically

### Portal Roles
- **Business Administrator:** Creates users, assigns roles (cannot view transactions)
- **Business Manager:** Approves transactions, views statements, withdraws funds

#### Creating Business Manager
1. Log in as Business Administrator → Select operators → Add
2. Enter username, select access channel as Web
3. Assign role, set password, submit KYC info

#### Creating API User
1. Log in as Business Administrator → Select operators → Add
2. Enter API initiator username, select access channel as API
3. Assign API roles, submit KYC info
4. Set password via Business Manager

### API Roles
- B2C: ORG B2C API Initiator
- Business Pay Bill: Business Paybill Org API initiator
- Business Buy Goods: Business Buy Goods Org API initiator
- Transaction Status: Transaction Status query ORG API
- Reversals: Org Reversals Initiator
- Tax Remittance: Tax Remittance to KRA API
- Set Password: Set Restricted ORG API PASSWORD

## Support
- **Chatbot:** Daraja Chatbot
- **Production Issues:** Incident Management Page or apisupport@safaricom.co.ke

## FAQs
1. **What is a short code?** Unique number allocated to a pay bill or buy goods organization
2. **What is C2B?** Customer to Business payment from customer wallet to merchant's short code
3. **C2B v1 vs v2?** v1 has SHA256 hashed MSISDN; v2 has masked MSISDN
4. **How often to register URLs?** Sandbox: before each simulation; Production: once
5. **How to delete URLs?** Self-managed on Daraja portal under Self Services → URL Management
6. **How to enable validation?** Email APISupport@safaricom.co.ke (~6 hours)
7. **Not receiving notifications?** Check URL validity and requirements
8. **What is validation URL?** Receives validation request before payment completion
9. **What is confirmation URL?** Receives notification upon payment completion
10. **How to get Business Administrator username?** Send official letter to M-PESABusiness@Safaricom.co.ke
