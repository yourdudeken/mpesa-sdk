# Query Org Info API

Retrieve organization/business information and account details.

**Endpoint:** `POST https://sandbox.safaricom.co.ke/mpesa/queryorginfo/v1/query`

## Overview
The Query Org Info API enables retrieval of organization account details, business information, and M-Pesa account configuration. This is useful for account validation and business information verification.

### Key Features
- Organization details retrieval
- Account configuration access
- Business information verification
- Account status checking
- API access permissions view
- Account hierarchy information

## How It Works
1. Business queries organization information
2. API validates credentials
3. M-Pesa retrieves organization data
4. Account details returned
5. Business verifies information
6. Integration with management systems

## Use Cases
- Account information verification
- Business detail confirmation
- API permissions checking
- Account status monitoring
- Multi-account management
- Integration validation
- Compliance verification

## Getting Started

### Prerequisites
- Daraja Account on Safaricom Developer Portal
- Sandbox app with API credentials
- Consumer Key & Consumer Secret
- Organization/business account
- Admin/Manager access level

### Good to Know
Query Org Info provides metadata about your M-Pesa business account and API access.

## Request Body
```json
{
  "AccessToken": "YYhZ20EF2nlgD2ekqK1Sy70b3eY"
}
```

## Request Parameter Definition

| Name | Description | Type | Sample |
|------|-------------|------|--------|
| AccessToken | OAuth access token from Authorization API | String | YYhZ20EF2nlgD2ekqK1Sy70b3eY |

## Response Body
```json
{
  "ResponseCode": "0",
  "ResponseDescription": "Success",
  "Organization": {
    "OrgName": "ABC Business Limited",
    "ShortCode": "600000",
    "AccountType": "PayBill",
    "Status": "Active",
    "Industry": "Retail",
    "Region": "Kenya"
  },
  "Accounts": [
    {
      "AccountNumber": "600000",
      "AccountType": "PayBill",
      "Status": "Active",
      "Currency": "KES",
      "CreatedDate": "2020-01-15"
    }
  ],
  "APIAccess": {
    "Permissions": ["B2C", "C2B", "TransactionStatus"],
    "Status": "Active"
  }
}
```

## Response Parameter Definition

| Name | Description | Type | Sample |
|------|-------------|------|--------|
| ResponseCode | Query result | String | 0 |
| ResponseDescription | Status message | String | Success |
| OrgName | Organization name | String | ABC Business Limited |
| ShortCode | Primary business short code | String | 600000 |
| AccountType | Type of account (PayBill, BuyGoods, B2C, etc) | String | PayBill |
| Status | Account status | String | Active |
| Industry | Business industry | String | Retail |
| Region | Operating region | String | Kenya |
| AccountNumber | Specific account number | String | 600000 |
| Currency | Account currency | String | KES |
| CreatedDate | Account creation date | Date | 2020-01-15 |
| Permissions | Available API permissions | Array | [B2C, C2B] |

## Organization Status

| Status | Description |
|--------|-------------|
| Active | Account is operational |
| Inactive | Account not in use |
| Suspended | Account temporarily suspended |
| Closed | Account closed |
| Pending | Account activation pending |

## Account Types

| Type | Description |
|------|-------------|
| PayBill | Regular bill payment account |
| BuyGoods | Retail goods/services account |
| B2C | Business to Customer disbursement |
| Utility | Utility bill collection |
| Till | Retail point-of-sale account |

## API Permissions Available

| Permission | Description |
|-----------|-------------|
| B2C | Business to Customer disbursement |
| C2B | Customer to Business collection |
| B2B | Business to Business transfer |
| TransactionStatus | Transaction status queries |
| Reversals | Transaction reversals |
| AccountBalance | Balance inquiries |
| BillManager | Bill reference management |

## Error Codes

| errorCode | errorMessage | Mitigation |
|-----------|-------------|------------|
| 401.002.01 | Invalid Access Token | Regenerate token |
| 404.002.01 | Organization not found | Verify account |
| 500.003.02 | Service unavailable | Retry request |

## Use Cases
- Pre-integration validation
- Multi-account management
- Permission verification
- Account status monitoring
- Business information updates
- Compliance reporting

## Testing
Use Daraja Simulator to test org info queries.

## Support
- **Chatbot:** Daraja Chatbot
- **Email:** apisupport@safaricom.co.ke
