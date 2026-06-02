# Lipa na Bonga API

Loyalty program rewards redemption through M-Pesa.

**Endpoint:** `POST https://sandbox.safaricom.co.ke/mpesa/lipanabonga/v1/redeem`

## Overview
The Lipa na Bonga API enables customers to redeem loyalty points and rewards through M-Pesa. Customers can convert accumulated Bonga points into cash or discounts on purchases.

### Key Features
- Points redemption capability
- Instant point-to-cash conversion
- Loyalty program integration
- Secure transaction handling
- Real-time balance updates
- Automated reward tracking

## How It Works
1. Customer checks Bonga points balance
2. Customer initiates redemption request
3. Points are validated and deducted
4. M-Pesa credit issued to customer account
5. Transaction confirmed and logged
6. Points balance updated

## Use Cases
- Customer loyalty rewards
- Bonga points conversion
- Promotional incentives
- Customer retention
- Retail rewards programs
- Subscription discounts

## Getting Started

### Prerequisites
- Daraja Account on Safaricom Developer Portal
- Sandbox app with API credentials
- Safaricom customer account
- Active Bonga points balance
- M-Pesa account registered

### Good to Know
Bonga points earned through Safaricom services can be redeemed as cash via M-Pesa.

## Request Body
```json
{
  "PhoneNumber": "254722000000",
  "Amount": "100",
  "TransactionReference": "BONGA123456",
  "Remarks": "Redeem Bonga points"
}
```

## Request Parameter Definition

| Name | Description | Type | Sample |
|------|-------------|------|--------|
| PhoneNumber | Customer phone number (format: 2547XXXXXXXX) | String | 254722000000 |
| Amount | Amount of points to redeem | Numeric | 100 |
| TransactionReference | Unique transaction identifier | String | BONGA123456 |
| Remarks | Redemption notes | String | Redeem Bonga points |

## Response Body
```json
{
  "ResponseCode": "0",
  "ResponseDescription": "Bonga points redeemed successfully",
  "TransactionID": "OA90000000",
  "PhoneNumber": "254722000000",
  "PointsRedeemed": "100",
  "CreditAmount": "100.00",
  "NewBalance": "500"
}
```

## Response Parameter Definition

| Name | Description | Type | Sample |
|------|-------------|------|--------|
| ResponseCode | Transaction result | String | 0 |
| ResponseDescription | Status message | String | Bonga points redeemed successfully |
| TransactionID | M-Pesa transaction ID | String | OA90000000 |
| PhoneNumber | Customer phone number | String | 254722000000 |
| PointsRedeemed | Number of points redeemed | Numeric | 100 |
| CreditAmount | KES amount credited | Numeric | 100.00 |
| NewBalance | Remaining Bonga points | Numeric | 500 |

## Redemption Options

| Option | Conversion Rate | Description |
|--------|-----------------|-------------|
| Direct Cash | 1 point = 1 KES | Instant cash redemption |
| Airtime | 1 point = 0.8 KES | Airtime top-up value |
| Discount | Variable | Merchant discounts |

## Bonga Points Balance

| Status | Description |
|--------|-------------|
| Active | Points available for redemption |
| Pending | Points not yet processed |
| Expired | Points past expiry date |
| Blocked | Points temporarily unavailable |

## Error Codes

| errorCode | errorMessage | Mitigation |
|-----------|-------------|------------|
| 400.002.02 | Invalid amount | Check redemption amount |
| 401.002.01 | Invalid credentials | Regenerate token |
| 404.002.01 | Customer not found | Verify phone number |
| 500.001.1001 | Insufficient points | Check point balance |
| 500.003.02 | Service unavailable | Retry request |

## Minimum Redemption
- Minimum redemption: 100 Bonga points
- Maximum per transaction: 10,000 Bonga points
- Daily limit: 50,000 Bonga points

## Transaction Limits

| Limit | Value |
|-------|-------|
| Minimum | 100 points |
| Maximum per transaction | 10,000 points |
| Daily maximum | 50,000 points |
| Monthly maximum | 200,000 points |

## Bonga Points Earning
Points earned through:
- Mobile voice calls
- SMS usage
- Data bundle purchases
- Promotional campaigns
- Airtime purchases

## Testing
Use Daraja Simulator with test customer accounts.

## Support
- **Chatbot:** Daraja Chatbot
- **Email:** apisupport@safaricom.co.ke
