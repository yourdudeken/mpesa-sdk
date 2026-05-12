# Business To Customer (B2C) API

Transact between an M-Pesa short code to a phone number registered on M-Pesa.

**Endpoint:** `POST https://sandbox.safaricom.co.ke/mpesa/b2c/v3/paymentrequest`

## Overview
B2C API is used to make payments from a Business to Customers' number, also known as Bulk Disbursements.

### How It Works
1. The Merchant prepares and sends a payment request to the M-PESA B2C API endpoint
2. The API Management Platform validates, authorizes, authenticates, and forwards to M-PESA
3. M-PESA validates the request and initiator details, then processes the transaction
4. M-PESA sends the transaction response via callback URL
5. The Customer receives an SMS notification confirming the payment

> **Note:** To go live, you must apply for a Bulk Disbursement Account and obtain a B2C Short Code or convert an existing pay bill/till number to a one account (can receive and disburse).

## Use Cases
- **Salary Payments:** Disbursement of employee wages
- **Cashback Payments:** Automated refunds or loyalty rewards
- **Promotional Payments:** Incentive payouts for campaigns
- **Winnings:** Contest/game payouts
- **Financial Institutions Withdrawal of Funds:** Transfers for cash-out
- **Loan Disbursements:** Distribution of approved loan amounts

## Getting Started

### Prerequisites
- Daraja Account on Safaricom Developer Portal
- Sandbox app with API credentials (Consumer Key & Consumer Secret)
- Initiator Username (API operator's username)
- Initiator Password (password with limited special characters: #, &, %, $)
- Public Key Certificate for encrypting initiator password
- Live M-PESA B2C shortcode/Paybill/Till number for production

### Good to Know
This API is asynchronous.

## Request Body
```json
{
  "OriginatorConversationID": "600997_Test_32et3241ed8yu",
  "InitiatorName": "testapi",
  "SecurityCredential": "RC6E9WDxXR4b9X2c6z3gp0oC5Th ==",
  "CommandID": "BusinessPayment",
  "Amount": "10",
  "PartyA": "600992",
  "PartyB": "254705912645",
  "Remarks": "remarked",
  "QueueTimeOutURL": "https://mydomain.com/path",
  "ResultURL": "https://mydomain.com/path",
  "Occassion": "ChristmasPay"
}
```

## Request Parameter Definition

| Name | Description | Type | Optional | Sample |
|------|-------------|------|----------|--------|
| OriginatorConversationID | Unique string to avoid double disbursement | String | No | 600997*Test*32et3241ed8yu |
| InitiatorName | Username of API user on M-PESA portal | String | No | johndoe |
| SecurityCredential | Encrypted password of API user | String | No | RC6E9WDx9X2c6z3gp0oC5Th== |
| CommandID | Transaction type: SalaryPayment, BusinessPayment, PromotionPayment | String | No | BusinessPayment |
| Amount | Transaction amount | Numeric | No | 11 |
| PartyA | B2C organization short code | Numeric | No | 600997 |
| PartyB | Customer's mobile number (12-digit MSISDN, e.g. 254705912645) | Numeric | No | 254705912645 |
| Remarks | Additional information (2-100 chars) | String | No | Any string |
| QueueTimeOutURL | URL for timeout notification | URL | No | https://mydomain.com/b2c/timedout |
| ResultURL | URL for processing result | URL | No | https://mydomain.com/b2c/result |
| Occassion | Additional information (1-100 chars) | String | Yes | ChristmasPay |

## Response Body
```json
{
  "ConversationID": "AG_20240706_20106e9209f64bebd05b",
  "OriginatorConversationID": "600997_Test_32et3241ed8yu",
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

## Successful Callback Payload
```json
{
  "Result": {
    "ResultType": 0,
    "ResultCode": 0,
    "ResultDesc": "The service request is processed successfully.",
    "OriginatorConversationID": "53e3-4aa8-9fe0-8fb5e4092cdd3533373",
    "ConversationID": "AG_20240706_2010364430d9bbdaf872",
    "TransactionID": "SG632NMUAB",
    "ResultParameters": {
      "ResultParameter": [
        { "Key": "TransactionAmount", "Value": 10 },
        { "Key": "TransactionReceipt", "Value": "SG632NMUAB" },
        { "Key": "ReceiverPartyPublicName", "Value": "254705912645 - NICHOLAS JOHN SONGOK" },
        { "Key": "TransactionCompletedDateTime", "Value": "06.07.2024 22:48:52" },
        { "Key": "B2CUtilityAccountAvailableFunds", "Value": 8959269.6 },
        { "Key": "B2CWorkingAccountAvailableFunds", "Value": 1199371.0 },
        { "Key": "B2CRecipientIsRegisteredCustomer", "Value": "Y" },
        { "Key": "B2CChargesPaidAccountAvailableFunds", "Value": -1980.0 }
      ]
    },
    "ReferenceData": {
      "ReferenceItem": { "Key": "QueueTimeoutURL", "Value": "https://internalsandbox.safaricom.co.ke/mpesa/b2cresults/v1/submit" }
    }
  }
}
```

## Unsuccessful Callback Payload
```json
{
  "Result": {
    "ResultType": 0,
    "ResultCode": 2001,
    "ResultDesc": "The initiator information is invalid.",
    "OriginatorConversationID": "53e3-4aa8-9fe0-8fb5e4092cdd3544366",
    "ConversationID": "AG_20240707_201062f6f6f5804f7a33",
    "TransactionID": "SG722NMVXQ",
    "ReferenceData": {
      "ReferenceItem": { "Key": "QueueTimeoutURL", "Value": "https://internalsandbox.safaricom.co.ke/mpesa/b2cresults/v1/submit" }
    }
  }
}
```

## Results Parameter Definition

| Parameter | Description | Type | Optional | Sample |
|-----------|-------------|------|----------|--------|
| Result | Root object for the result message | JSON Object | No | { ... } |
| ResultType | Status code (0 = sent to listener) | Numeric | No | 0 |
| ResultCode | Transaction processing status (0 = success) | String | No | 0 |
| ResultDesc | Status description | String | No | The service request is processed successfully. |
| OriginatorConversationID | Unique identifier from API proxy | String | No | 53e3-4aa8-9fe0-... |
| ConversationID | Unique identifier from M-PESA | String | No | AG_20240707_20106f7a33 |
| TransactionID | Unique M-PESA transaction ID | String | No | SG722NMVXQ |

### Additional Parameters (Successful Requests)

| Parameter | Description | Type | Optional |
|-----------|-------------|------|----------|
| ResultParameters | Object holding more transaction details | JSON Object | Yes |
| ResultParameter | Array within ResultParameters | JSON Array | Yes |
| TransactionAmount | Amount transacted | Numeric | Yes |
| TransactionReceipt | Unique M-PESA transaction ID | String | Yes |
| ReceiverPartyPublicName | Customer name and phone | String | Yes |
| TransactionCompletedDateTime | Completion date/time | String | Yes |
| B2CUtilityAccountAvailableFunds | Utility account balance | Decimal | Yes |
| B2CWorkingAccountAvailableFunds | Working account balance | Decimal | Yes |
| B2CRecipientIsRegisteredCustomer | Registration status (Y/N) | Character | Yes |
| B2CChargesPaidAccountAvailableFunds | Charges paid account balance | Decimal | Yes |

### Reference Data Parameters

| Parameter | Description | Type |
|-----------|-------------|------|
| ReferenceData | Object holding additional request details | JSON Object |
| ReferenceItem | Object within ReferenceData | JSON Object |

## Sample Error Response
```json
{
  "requestId": "1c5b-4ba8-815c-ac45c57a3db01469899",
  "errorCode": "500.002.1001",
  "errorMessage": "Duplicate OriginatorConversationID."
}
```

| Name | Description | Sample | Type |
|------|-------------|--------|------|
| requestId | Unique identifier by API gateway | 30764-19833054-1 | String |
| errorCode | Unique error code | 401.001 | String |
| errorMessage | Descriptive failure message | Bad Request - Invalid amount | String |

## Response Codes

| Response Code | Response Description |
|---------------|---------------------|
| 0 | Success |

## Result Codes

| ResultCode | ResultDesc | Explanation |
|-----------|-----------|-------------|
| 0 | The service request is processed successfully. | B2C transaction processed successfully |
| 1 | The balance is insufficient for the transaction. | Insufficient utility account balance |
| 2 | Declined due to limit rule | Amount smaller than minimum allowed |
| 3 | Declined due to limit rule: greater than the maximum transaction amount. | Amount exceeds max (Ksh 250,000) |
| 4 | Declined due to limit rule: would exceed daily transfer limit | Exceeds daily limit (Ksh 500,000) |
| 8 | Declined due to limit rule: would exceed the maximum balance. | Exceeds max customer balance (Ksh 500,000) |
| 11 | The DebitParty is in an invalid state. | B2C account not active |
| 21 | The initiator is not allowed to initiate this request | API user lacks ORG B2C API initiator role |
| 2001 | The initiator information is invalid. | Invalid API user credentials |
| 2006 | Declined due to account rule | B2C account not active |
| 2028 | The request is not permitted according to product assignment. | PartyA has no B2C permission |
| 2040 | Credit Party customer type can't be supported by the service. | Customer not registered |
| 8006 | The security credential is locked | API user password locked |
| SFC_IC0003 | The operator does not exist. | Invalid phone number |

## Testing

### Option 1: Daraja Simulator
Create a test app, select B2C product, and use the simulator.

### Option 2: Postman
Generate access token and initiate transactions.

**Sandbox Token Endpoint:** `https://sandbox.safaricom.co.ke/oauth/v1/generate?grant_type=client_credentials`
**Production Token Endpoint:** `https://api.safaricom.co.ke/oauth/v1/generate?grant_type=client_credentials`

## Go Live
Attach integration to a live pay bill/till number. Navigate to GO LIVE tab and fill in live data.

## M-PESA Organization Portal

### Access
URL: https://org.ke.m-pesa.com/orglogin.action

### First-Time Login
1. Launch `https://org.ke.m-pesa.com`
2. Enter Short code (Bulk payment number)
3. Enter Business Administrator username
4. Enter first-time password (case-sensitive)
5. Enter Verification Code, click login
6. Enter OTP to change password
7. Set new password, security questions, submit

### Account Types
- **MMF/Working/M-PESA Account:** For business withdrawals
- **Utility Account:** Receives customer payments
- **Charges Paid Account:** Debited for transaction charges
- **Organization Settlement Account:** Settles charges and moves balance

### Portal Roles

#### Business Administrator
- Creates users and assigns roles
- Cannot view transactions
- Created by Safaricom

#### Business Manager
- Views statements, initiates/approves transactions, withdraws funds

#### API User Creation
1. Log in as Business Administrator
2. Select operators → Add → Enter API initiator username
3. Select access channel as API
4. Assign API roles, submit KYC info
5. Set password via Business Manager

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
- **Why no callbacks?** Ensure ResultURL is publicly accessible
- **Invalid access token?** Token expires hourly; regenerate
- **Can I use actual short code on simulator?** No, only sandbox test short codes
- **How to know if B2C is successful?** ResultCode 0 in callback
- **Is there a test environment?** Yes, sandbox environment
- **Balance insufficient yet money exists?** B2C debits Utility account (not MMF/Working)
- **How to move funds to Utility?** Via M-PESA portal or B2B API (BusinessTransferFromMMFToUtility)
- **What are transaction limits?** Min Ksh 10, Max per transaction Ksh 250,000, Daily max Ksh 500,000
- **Which API role is needed?** ORG B2C API initiator
- **How to activate pending API user?** Set password via user with Set Restricted ORG API PASSWORD role
- **Initiator information is invalid?** Check username, encrypted password, algorithm/certificate
- **Security credential locked?** Multiple failed attempts; unlock via Business Administrator
- **Invalid API call - no apiproduct match?** Ensure B2C product is enabled on Daraja app
- **Passkey for B2C?** No, only for M-PESA Express/STK Push
- **Can I reverse a B2C transaction?** No, must be done manually on M-PESA portal
- **Bad Request - Invalid InitiatorName?** Check Content-Type, parameter name, value
- **Can B2C and C2B coexist?** Yes, if both enabled for short code
- **Generate security credential for every request?** No, can reuse
- **Manual approval needed?** No
- **Check B2C transaction status?** Yes, via Transaction Status Query API
- **API for moving money to B2C short code?** Yes, B2C Account Top Up API
