# Bill Manager API

Manage and track bills for Pay Bill and Buy Goods accounts.

**Endpoint:** `POST https://sandbox.safaricom.co.ke/mpesa/billmanager/v1/updatebillreference`

## Overview
The Bill Manager API allows businesses to create, update, and manage bill references for their Pay Bill and Buy Goods accounts. This simplifies bill tracking and payment reconciliation for customers.

### Key Features
- Create bill references for customers
- Update existing bill information
- Track bill payment status
- Automated bill notifications
- Integration with M-Pesa payment flow

## How It Works
1. Business creates/updates bill reference in system
2. Bill information stored and linked to customer
3. When customer makes payment, bill reference used
4. Payment system matches payment to bill
5. Business receives payment notification with bill details
6. Automatic reconciliation and reporting

## Use Cases
- Utility bill management
- Invoice tracking
- Subscription billing
- Loan payment tracking
- Insurance premium billing
- Student fee management

## Getting Started

### Prerequisites
- Daraja Account on Safaricom Developer Portal
- Sandbox app with API credentials
- Pay Bill or Buy Goods merchant account
- Customer database with reference IDs

### Good to Know
Bill references help match customer payments to specific invoices or accounts in your system.

## Request Body
```json
{
  "BillRefName": "customer_ref_2024",
  "DueDate": "2024-12-31",
  "Amount": "10000",
  "InvoiceNumber": "INV-001-2024",
  "AccountReference": "ACC-12345",
  "PhoneNumber": "254722000000",
  "Email": "customer@example.com",
  "Description": "Monthly subscription bill"
}
```

## Request Parameter Definition

| Name | Description | Type | Sample |
|------|-------------|------|--------|
| BillRefName | Name of bill reference | String | customer_ref_2024 |
| DueDate | Bill due date (YYYY-MM-DD) | Date | 2024-12-31 |
| Amount | Bill amount in KES | Numeric | 10000 |
| InvoiceNumber | Invoice reference number | String | INV-001-2024 |
| AccountReference | Customer account reference | String | ACC-12345 |
| PhoneNumber | Customer phone number | String | 254722000000 |
| Email | Customer email address | String | customer@example.com |
| Description | Bill description | String | Monthly subscription |

## Response Body
```json
{
  "OriginatorConversationID": "515-5258779-3",
  "ConversationID": "AG_20200123_0000417fed8ed666e976",
  "ResponseCode": "0",
  "ResponseDescription": "Bill reference created successfully"
}
```

## Response Parameter Definition

| Name | Description | Type | Sample |
|------|-------------|------|--------|
| OriginatorConversationID | Unique identifier of request | String | 515-5258779-3 |
| ConversationID | Unique identifier from M-Pesa | String | AG_20200123_0000417fed8ed666e976 |
| ResponseCode | Success indicator | String | 0 |
| ResponseDescription | Status description | String | Bill reference created successfully |

## Bill Reference Status

| Status | Description |
|--------|-------------|
| Active | Bill is active and ready for payment |
| Paid | Bill has been paid in full |
| Overdue | Bill payment date has passed |
| Partial | Bill has been partially paid |
| Cancelled | Bill has been cancelled |

## Error Codes

| errorCode | errorMessage | Mitigation |
|-----------|-------------|------------|
| 400.002.02 | Bad Request - Invalid Data | Ensure all fields are correct |
| 400.002.03 | Duplicate Bill Reference | Use unique reference per bill |
| 401.002.01 | Invalid Access Token | Regenerate token |
| 404.002.01 | Resource not found | Check bill reference |
| 500.001.1001 | Internal Server Error | Retry request |

## Testing
Use Daraja Simulator to create test bill references.

## Best Practices
- Use unique bill reference IDs
- Set realistic due dates
- Include detailed descriptions
- Maintain accurate customer information
- Regular bill status monitoring

## Support
- **Chatbot:** Daraja Chatbot
- **Email:** apisupport@safaricom.co.ke
