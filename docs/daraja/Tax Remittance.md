# Tax Remittance API

Facilitate tax payment remittances to Kenya Revenue Authority through M-Pesa.

**Endpoint:** `POST https://sandbox.safaricom.co.ke/mpesa/taxremittance/v1/remit`

## Overview
The Tax Remittance API enables businesses to pay taxes directly to Kenya Revenue Authority (KRA) through M-Pesa. This includes income tax, VAT, corporation tax, and other tax obligations through seamless integration.

### Key Features
- Direct tax payment to KRA
- Multiple tax type support
- Real-time payment processing
- Automatic receipt generation
- Tax compliance tracking
- Payment reconciliation

## How It Works
1. Business initiates tax payment through API
2. Tax amount and reference provided
3. Payment routed to KRA account
4. KRA processes payment immediately
5. Receipt issued and confirmed
6. Payment recorded in tax account

## Use Cases
- Income tax remittance
- VAT payment
- Corporation tax payments
- Customs duty payment
- Excise tax payments
- Professional tax
- Motor vehicle tax
- Penalty payments

## Getting Started

### Prerequisites
- Daraja Account on Safaricom Developer Portal
- Sandbox app with API credentials
- Business registered with KRA
- Valid KRA PIN
- Tax account setup
- Admin/Manager access

### Good to Know
Tax remittance directly credited to KRA account. Ensure correct tax type and reference.

## Request Body
```json
{
  "InitiatorName": "taxuser",
  "SecurityCredential": "SAFVNChNHfVtXEZMBuVo+a1Hwr+DtrUVN3zVg==",
  "CommandID": "TaxRemittance",
  "ShortCode": "600000",
  "TaxType": "IncomeTax",
  "KRAPINNumber": "A002345678Q",
  "Amount": "100000",
  "TransactionReference": "TAX-MONTHLY-JAN-2024",
  "Description": "January income tax"
}
```

## Request Parameter Definition

| Name | Description | Type | Sample |
|------|-------------|------|--------|
| InitiatorName | Initiator username | String | taxuser |
| SecurityCredential | Encrypted password | String | EToK4lNR... |
| CommandID | Command type (TaxRemittance) | String | TaxRemittance |
| ShortCode | Business short code | String | 600000 |
| TaxType | Tax category (IncomeTax, VAT, CorporationTax, etc) | String | IncomeTax |
| KRAPINNumber | Business KRA PIN | String | A002345678Q |
| Amount | Tax payment amount in KES | Numeric | 100000 |
| TransactionReference | Unique reference for tracking | String | TAX-MONTHLY-JAN-2024 |
| Description | Payment description | String | January income tax |

## Response Body
```json
{
  "ResponseCode": "0",
  "ResponseDescription": "Tax remittance processed successfully",
  "TransactionID": "OA90000000",
  "KRAPINNumber": "A002345678Q",
  "TaxType": "IncomeTax",
  "Amount": "100000",
  "ReceiptNumber": "KRA/2024/0000001",
  "PaymentDate": "2024-01-15 14:30:00",
  "Status": "Paid"
}
```

## Response Parameter Definition

| Name | Description | Type | Sample |
|------|-------------|------|--------|
| ResponseCode | Payment result | String | 0 |
| ResponseDescription | Status message | String | Tax remittance processed successfully |
| TransactionID | M-Pesa transaction ID | String | OA90000000 |
| KRAPINNumber | Business KRA PIN | String | A002345678Q |
| TaxType | Tax category | String | IncomeTax |
| Amount | Amount paid | Numeric | 100000 |
| ReceiptNumber | KRA receipt number | String | KRA/2024/0000001 |
| PaymentDate | Payment date/time | DateTime | 2024-01-15 14:30:00 |
| Status | Payment status | String | Paid |

## Tax Types Supported

| Tax Type | Description |
|----------|-------------|
| IncomeTax | Personal/business income tax |
| VAT | Value Added Tax |
| CorporationTax | Corporate tax |
| CustomsDuty | Import/export customs |
| ExciseTax | Excise duties |
| ProfessionalTax | Professional services tax |
| MotorVehicleTax | Vehicle registration tax |
| Penalties | Tax penalties and fines |
| Interest | Tax interest charges |

## KRA PIN Format
Kenya Revenue Authority PIN format: A###########Q
- First character: A (for individual) or P (for partnership) or C (for company)
- 9-10 digits
- Last character: Q (checksum)

## Payment Status

| Status | Description |
|--------|-------------|
| Submitted | Payment submitted to KRA |
| Processing | KRA processing payment |
| Paid | Payment successfully processed |
| Failed | Payment failed |
| Cancelled | Payment cancelled |

## Error Codes

| errorCode | errorMessage | Mitigation |
|-----------|-------------|------------|
| 400.002.02 | Invalid KRA PIN | Verify PIN format |
| 400.002.03 | Invalid tax type | Select valid tax type |
| 401.002.01 | Invalid credentials | Regenerate token |
| 404.002.01 | Tax account not found | Register with KRA |
| 500.001.1001 | Insufficient funds | Check balance |
| 500.003.02 | Service unavailable | Retry request |

## Tax Payment Limits

| Limit | Value |
|-------|-------|
| Minimum | KES 100 |
| Maximum per transaction | KES 5,000,000 |
| Daily limit | Unlimited |

## Receipt and Compliance
Tax receipts automatically:
- Generated and sent to business
- Recorded in KRA system
- Sent to KRA registered email
- Tracked for compliance reporting

## Payment Tracking
Monitor tax payments through:
- Transaction reference number
- KRA receipt number
- M-Pesa transaction ID
- M-Pesa portal history
- KRA online portal

## Testing
Use Daraja Simulator with test KRA PIN and tax reference.

## Best Practices
- Keep accurate KRA PIN records
- Submit payments before due dates
- Maintain payment receipts
- Track all remittances
- Regular compliance checks
- Monitor payment confirmations

## Support
- **Chatbot:** Daraja Chatbot
- **Email:** apisupport@safaricom.co.ke
- **KRA Support:** www.kra.go.ke
